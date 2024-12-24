package main

import (
	"fmt"
	"net/http"
	"strconv"

   )
   
   // Обработчик для параметров name и age
   func nameAgeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")
   
	if name == "" || age == "" {
	 http.Error(w, "Параметры 'name' и 'age' обязательны", http.StatusBadRequest)
	 return
	}
   
	fmt.Fprintf(w, "Меня зовут %s, мне %s лет", name, age)
   }
   
   // Функция для вычисления арифметических операций
   func calculate(w http.ResponseWriter, r *http.Request, operation string) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")
   
	if aStr == "" || bStr == "" {
	 http.Error(w, "Параметры 'a' и 'b' обязательны", http.StatusBadRequest)
	 return
	}
   
	a, err1 := strconv.ParseFloat(aStr, 64)
	b, err2 := strconv.ParseFloat(bStr, 64)
   
	if err1 != nil || err2 != nil {
	 http.Error(w, "Параметры 'a' и 'b' должны быть числами", http.StatusBadRequest)
	 return
	}
   
	var result float64
	switch operation {
	case "add":
	 result = a + b
	case "sub":
	 result = a - b
	case "mul":
	 result = a * b
	case "div":
	 if b == 0 {
	  http.Error(w, "Деление на 0 невозможно", http.StatusBadRequest)
	  return
	 }
	 result = a / b
	default:
	 http.Error(w, "Неизвестная операция", http.StatusBadRequest)
	 return
	}
   
	fmt.Fprintf(w, "Результат: %f", result)
   }
   
   func addHandler(w http.ResponseWriter, r *http.Request) {
	calculate(w, r, "add")
   }
   
   func subHandler(w http.ResponseWriter, r *http.Request) {
	calculate(w, r, "sub")
   }
   
   func mulHandler(w http.ResponseWriter, r *http.Request) {
	calculate(w, r, "mul")
   }
   
   func divHandler(w http.ResponseWriter, r *http.Request) {
	calculate(w, r, "div")
   }



   

   func main() {
	// Маршрут для обработки name и age
	http.HandleFunc("/name_age", nameAgeHandler)
   
	// Маршруты для арифметических операций
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/sub", subHandler)
	http.HandleFunc("/mul", mulHandler)
	http.HandleFunc("/div", divHandler)
   
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)


	
   }
   