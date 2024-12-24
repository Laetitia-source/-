package main 
import (
	"errors"
	"fmt" 
	"math"
)

// Функция formatIP принимает массив из 4 байтов и возвращает строку в формате IP-адреса
func formatIP(ip [4]byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
   }



// Функция listEven принимает два целых числа, представляющих диапазон, и возвращает срез четных чисел и ошибку
func listEven(left, right int) ([]int, error) {
	// Проверка на корректность диапазона
	if left > right {
	 return nil, errors.New("левая граница больше правой")
	}
   
	var evens []int
   
	// Заполняем срез четными числами
	for i := left; i <= right; i++ {
	 if i%2 == 0 {
	  evens = append(evens, i)
	 }
	}
   
	return evens,nil
   }



// Функция, которая подсчитывает количество вхождений каждого символа в строке
func countCharacters(s string) map[rune]int {
	// Создаем карту для хранения количества вхождений символов
	charCount := make(map[rune]int)
   
	// Проходим по строке с помощью range, который дает символ и его индекс
	for _, char := range s {
	 // Увеличиваем счетчик для текущего символа
	 charCount[char]++
	}
   
	return charCount
   }



// Определяем структуру «Точка»
type Point struct {
 X, Y float64
}

// Определяем структуру «Отрезок»
type Segment struct {
 Start, End Point
}

// Метод для структуры «Отрезок», который возвращает длину отрезка
func (s Segment) Length() float64 {
 return math.Sqrt(math.Pow(s.End.X-s.Start.X, 2) + math.Pow(s.End.Y-s.Start.Y, 2))
}

// Определяем структуру «Треугольник»
type Triangle struct {
 A, B, C Point
}

// Метод Area для структуры «Треугольник», возвращающий площадь
func (t Triangle) Area() float64 {
 // Вычисляем длины сторон
 ab := Segment{t.A, t.B}.Length()
 bc := Segment{t.B, t.C}.Length()
 ca := Segment{t.C, t.A}.Length()

 // Используем формулу Герона для вычисления площади треугольника
 s := (ab + bc + ca) / 2
 return math.Sqrt(s * (s - ab) * (s - bc) * (s - ca))
}

// Определяем структуру «Круг»
type Circle struct {
 Center Point
 Radius float64
}

// Метод Area для структуры «Круг», возвращающий площадь
func (c Circle) Area() float64 {
 return math.Pi * math.Pow(c.Radius, 2)
}

// Определяем интерфейс «Фигура»
type Shape interface {
 Area() float64
}

// Функция для вывода площади фигуры
func printArea(s Shape) {
 result := s.Area()
 fmt.Printf("Площадь фигуры: %.2f\n", result)
}

   
   







   
   func main() {
 // Пример использования 1.1
	ip := [4]byte{127, 0, 0, 1}
	fmt.Println(formatIP(ip)) // Выведет: 127.0.0.1

	
 // Пример использования 1.2
	evens, err := listEven(10, 20)
	if err != nil {
	 fmt.Println("Ошибка:", err)
	} else {
	 fmt.Println("Чётные числа:", evens)
	}
   
	// Пример с ошибкой (левая граница больше правой)
	evens, err = listEven(20, 10)
	if err != nil {
	 fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Чётные числа:", evens)
	}



 // Пример использования 2
 str := "hello, world"
 // Получаем карту с количеством вхождений символов
 result := countCharacters(str)

 // Выводим результат
 for char, count := range result {
  fmt.Printf("'%c': %d\n", char, count)
 }


// Пример использования 3
 // Пример с отрезком
 point1 := Point{X: 0, Y: 0}
 point2 := Point{X: 3, Y: 4}
 segment := Segment{Start: point1, End: point2}
 fmt.Printf("Длина отрезка: %.2f\n", segment.Length())

 // Пример с треугольником
 triangle := Triangle{
  A: Point{X: 0, Y: 0},
  B: Point{X: 3, Y: 0},
  C: Point{X: 3, Y: 4},
 }
 printArea(triangle) // Вызов функции с треугольником

 // Пример с кругом
 circle := Circle{
  Center: Point{X: 0, Y: 0},
  Radius: 5,
 }
 printArea(circle) // Вызов функции с кругом
	 

   }