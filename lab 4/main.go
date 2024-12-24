

package main

import (
    "time"
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "errors"

    _ "github.com/lib/pq"
    "golang.org/x/crypto/bcrypt"
    "github.com/gorilla/mux"
    "github.com/golang-jwt/jwt/v4"

)
const secretKey = "your-256-bit-secret"

var db *sql.DB

// Инициализация базы данных
func initDB() {
    var err error
    connStr := "user=postgres password=laetitia dbname=auth_service sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Ошибка подключения к базе данных: %v", err)
    }
    if err = db.Ping(); err != nil {
        log.Fatalf("База данных недоступна: %v", err)
    }
    fmt.Println("Успешное подключение к базе данных")
}

// Обработчик для регистрации пользователя
func registerHandler(w http.ResponseWriter, r *http.Request) {
    type Customer struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    var customer Customer
    err := json.NewDecoder(r.Body).Decode(&customer)
    if err != nil {
        http.Error(w, "Некорректные данные", http.StatusBadRequest)
        return
    }

    // Проверяем, существует ли пользователь с таким email
    var exists bool
    err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM customers WHERE email=$1)`, customer.Email).Scan(&exists)
    if err != nil {
        http.Error(w, "Ошибка проверки существующего пользователя", http.StatusInternalServerError)
        return
    }
    if exists {
        http.Error(w, "Пользователь с таким email уже существует", http.StatusConflict)
        return
    }

    // Хешируем пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Ошибка при хешировании пароля", http.StatusInternalServerError)
        return
    }

    // Сохраняем пользователя в базу данных
    query := `INSERT INTO customers (name, email, password) VALUES ($1, $2, $3)`
    _, err = db.Exec(query, customer.Name, customer.Email, string(hashedPassword))
    if err != nil {
        http.Error(w, "Ошибка при сохранении данных", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Пользователь успешно зарегистрирован"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    type LoginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    var req LoginRequest
    // Декодируем тело запроса
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Некорректные данные", http.StatusBadRequest)
        return
    }

    // Проверяем, существует ли пользователь с таким email
    var id int
    var hashedPassword string
    err = db.QueryRow(`SELECT id, password FROM customers WHERE email=$1`, req.Email).Scan(&id, &hashedPassword)
    if err == sql.ErrNoRows {
        http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "Ошибка при запросе данных пользователя", http.StatusInternalServerError)
        return
    }

    // Проверяем пароль
    log.Println("Hashed password from DB:", hashedPassword)
    log.Println("Password provided:", req.Password)
    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
    if err != nil {
        http.Error(w, "Неверный пароль", http.StatusUnauthorized)
        return
    }

    // Генерируем JWT
    expirationTime := time.Now().Add(1 * time.Hour).Unix() // Токен действителен 1 час
    payload := jwt.MapClaims{
        "sub": fmt.Sprintf("%d", id), // ID пользователя
        "exp": expirationTime,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
    tokenStr, err := token.SignedString([]byte(secretKey))
    if err != nil {
        http.Error(w, "Ошибка при генерации токена", http.StatusInternalServerError)
        return
    }

    // Возвращаем токен клиенту
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}

// CustomClaims
type CustomClaims struct {
    Sub string `json:"sub"` // ID пользователя
    Exp int64  `json:"exp"` // Время истечения токена
}


func (c CustomClaims) Valid() error {
    if time.Now().Unix() > c.Exp {
        return errors.New("token expired")
    }
    return nil
}

func parseJWT(tokenStr string) (*CustomClaims, error) {
    claims := &CustomClaims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })
    if err != nil {
        return nil, fmt.Errorf("Ошибка проверки токена: %v", err)
    }

    if !token.Valid {
        return nil, fmt.Errorf("Токен недействителен")
    }

    return claims, nil
}

func getCustomerHandler(w http.ResponseWriter, r *http.Request) {
    // Получаем токен из заголовка Authorization
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Токен отсутствует", http.StatusForbidden)
        return
    }

    // Проверяем и парсим токен
    claims, err := parseJWT(authHeader)
    if err != nil {
        http.Error(w, err.Error(), http.StatusForbidden)
        return
    }

    // Получаем id пользователя из токена
    userID := claims.Sub // "sub" - это ID пользователя в токене

    // Получаем id пользователя из URL
    vars := mux.Vars(r)
    requestedID := vars["id"]

    // Проверяем соответствие id из токена и id из запроса
    if userID != requestedID {
        http.Error(w, "Доступ запрещен", http.StatusForbidden)
        return
    }

    // Запрашиваем данные пользователя из базы данных
    var name, email string
    err = db.QueryRow(`SELECT name, email FROM customers WHERE id=$1`, userID).Scan(&name, &email)
    if err == sql.ErrNoRows {
        http.Error(w, "Пользователь не найден", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Ошибка при запросе данных", http.StatusInternalServerError)
        return
    }

    // Отправляем данные пользователя клиенту
    response := map[string]string{
        "id":    userID,
        "name":  name,
        "email": email,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    // Инициализация базы данных
    initDB()
    defer db.Close()

    // Создаем маршрутизатор
    router := mux.NewRouter()

    // Добавляем маршруты
    router.HandleFunc("/register", registerHandler).Methods("POST") // Регистрация пользователя
    router.HandleFunc("/login", loginHandler).Methods("POST")       // Авторизация пользователя
    router.HandleFunc("/customer/{id}", getCustomerHandler).Methods("GET")   //Получение данных пользователя
    
    // Запуск сервера
    fmt.Println("Сервер запущен на порту 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}