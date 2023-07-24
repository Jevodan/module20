package main

import (
	"fmt"
	"module20/ring"
	"time"
)

// Стадия фильтрации отрицательных чисел (не пропускать отрицательные числа).
var filterOne = func(channel chan int, done chan int) (filterChannel chan int) {
	filterChannel = make(chan int)
	go func() {
		for {
			select {
			case v, isOpen := <-channel:
				if isOpen {
					if v >= 0 {
						filterChannel <- v
					}
				}

			case <-done:
				fmt.Println("Конвеер 1 закрыт")
				close(filterChannel)
				return
			}
		}
	}()
	return
}

// Стадия фильтрации чисел, не кратных 3 (не пропускать такие числа), исключая также и 0.
var filterTwo = func(channel chan int, done chan int) (filterChannel chan int) {
	filterChannel = make(chan int)
	go func() {
		for {
			select {
			case v := <-channel:
				if v != 0 && v%3 == 0 {
					filterChannel <- v
					time.Sleep(time.Second)
				}
			case <-done:
				fmt.Println("Конвеер 2 закрыт")
				close(filterChannel)
				return
			}
		}
	}()
	return
}

// Стадия буферизации данных
var filterThree = func(channel chan int, done chan int) (filterChannel chan int) {
	r := ring.NewRing(BUFFER_SIZE)
	filterChannel = make(chan int)
	go func() {
		for {
			select {
			case v := <-channel:
				r.SetValue(v)
				r = r.Next()
			case <-done:
				fmt.Println("Конвеер 3 закрыт")
				close(filterChannel)
				return
			}
		}

	}()

	go func() {
		for {
			select {
			case <-time.After(INTERVAL):
				arrayInt := r.Get()
				fmt.Print("Принятая продукция на третью стадию: ")
				for _, v := range arrayInt {
					if v != 0 {
						fmt.Print(v, " ")
					}
				}
				fmt.Println()
				for _, v := range arrayInt {
					if v != 0 {
						filterChannel <- v
					}
				}
			case <-done:
				return
			}
		}

	}()

	return filterChannel
}
