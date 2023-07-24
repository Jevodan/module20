package main

import (
	"fmt"
)

func main() {
	// Создаем новое кольцо размером 5
	r := NewRing(5)

	// Получаем длину кольца
	n := r.Len()
	fmt.Println(n)
	// Инициализируем кольцо
	// некоторыми целочисленными значениями
	for i := 0; i < 100; i++ {
		if i > n {
			fmt.Println("Конец кольца, пора его очистить")
			break
		}
		r.Value = i
		r = r.Next()

	}

	n = r.Len()
	fmt.Println(r.Get(), r.Len())
	n = r.Len()
	fmt.Println(r.Get(), r.Len())

	for i := 0; i < 100; i++ {
		if i > n {
			fmt.Println("Конец кольца, пора его очистить")
			break
		}
		r.Value = i
		r = r.Next()

	}
	n = r.Len()
	fmt.Println(r.Get(), r.Len())

}
