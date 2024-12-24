package main
import (
	"errors"
	"fmt"
)

// Функция hello принимает строку и возвращает строку с приветствием.
func hello(name string) {
	fmt.Println("Hello, " + name + "!")
}

// Функция printEven выводит все чётные числа в диапазоне [a, b].
// Если a > b, возвращается ошибка.
func printEven(a, b int64) error {
	// Проверка, что левая граница не больше правой
	if a > b {
		return errors.New("левая граница больше правой")
	}

	// Цикл для вывода всех чётных чисел в диапазоне
    for i := a; i <= b; i++ {
        if i%2 == 0 {
            fmt.Println(i)
        }
    }

	// Если нет ошибок, возвращаем nil
	return nil	
}

// Функция apply выполняет арифметические операции с двумя числами
func apply(a, b float64, operator string)(float64, error)  {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		// Провераем деление на 0
		if b == 0 {
			return 0, errors.New("деление на ноль невозможно")		
		}
		return a / b, nil
	default:
		// Возвращаем ошибку, если оператор не поддерживается
		return 0, errors.New("неподдерживаемый оператор")	
	}
	
}


func main() {
	// Пример использования функции hello
	hello( "Laetitia" )
	hello("Astrid")


	// Пример использования функции printEven
	// Пример правильного вызова: диапазон [2, 10]
    fmt.Println("Вызов функции printEven(2, 10) :")
	err := printEven(2, 10)
	if err != nil {
	 fmt.Println("Ошибка:", err)
	} else {
	 fmt.Println("Четные числа выведены успешно.")
	}

	// Пример некорректного вызова: диапазон [10, 2]
    fmt.Println("Вызов функции printEven(10, 2) :")
	err = printEven(10, 2)
	if err != nil {
	 fmt.Println("Ошибка:", err)
	} else {
	 fmt.Println("Четные числа выведены успешно.")
	}


	// Пример использования функции apply
	result, err := apply(20, 4, "+")
	if err != nil {
		fmt.Println("Ошибка:", err)	
	} else {
		fmt.Println("Результат сложения:", result)
	}

	result, err = apply(5, 5, "-")
	if err != nil {
		fmt.Println("Ошибка:", err)	
	} else {
		fmt.Println("Результат вычитания:", result)
	}

	result, err = apply(6, 5, "*")
	if err != nil {
		fmt.Println("Ошибка:", err)	
	} else {
		fmt.Println("Результат умножения:", result)
	}

    result, err = apply(1, 5, "/")
	if err != nil {
		fmt.Println("Ошибка:", err)	
	} else {
		fmt.Println("Результат деления:", result)
	}

	// Пример деления на 0
	result, err = apply(5, 0, "/")
	if err != nil {
		fmt.Println("Ошибка:", err)	
	} else {
		fmt.Println("Результат:", result)
	}

	// Пример исползования неподдержтваемого оперетора
	result, err = apply(25, 5, "^")
	if err != nil {
		fmt.Println("Ошибка:", err)	
	} else {
		fmt.Println("Результат:", result)
	}


}



