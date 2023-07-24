package main

import (
	"bufio"
	"fmt"
	"module20/ring"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	BUFFER_SIZE               = 50               // Буфер кольцевого массива
	INTERVAL    time.Duration = time.Second * 10 // Интервал времени, за который происходит очистка кольцевого буфера
)

type Stadies func(channel chan int, done chan int) (filterChannel chan int) // Тип данных Стадия

type konveer struct {
	stage []Stadies
	done  chan int
}

func NewKonveer(done chan int, stadies ...Stadies) *konveer {
	k := konveer{stage: stadies, done: done}
	return &k
}

/*
Запуск конвеера. выходной  канал в цикле предыдущей стадии
является входным каналом слкдующей стадии. Возвращает выходной канал последней стадии
source - источник данных и входной канал первой стадии
с - обработанные данные после конвеера
*/
func (k *konveer) start(source chan int) chan int {
	var c chan int = source
	for _, v := range k.stage {
		c = v(c, k.done) // Вызов внутренних функций filterOne filterSecond filterThree
	}
	return c
}

func main() {
	//Стадия фильтрации отрицательных чисел (не пропускать отрицательные числа).
	filterOne := func(channel chan int, done chan int) (filterChannel chan int) {
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
	//Стадия фильтрации чисел, не кратных 3 (не пропускать такие числа), исключая также и 0.
	filterTwo := func(channel chan int, done chan int) (filterChannel chan int) {
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

	//Стадия буферизации данных
	filterThree := func(channel chan int, done chan int) (filterChannel chan int) {
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

	// источник данных
	source, done := sourceData()
	// создание конввера
	konv := NewKonveer(done, filterOne, filterTwo, filterThree)
	// запуск конвеера
	truba := konv.start(source)

	for v := range truba {
		time.Sleep(time.Second * 2)
		fmt.Println("Конечная продукция: ", v)
	}
	time.Sleep(time.Second * 2)
}

func sourceData() (chan int, chan int) {
	var (
		str string
	)
	channel := make(chan int, 10)
	done := make(chan int)

	go func() {
		defer close(done)
		fmt.Print("Введите числа через пробел и формируйте источники данных, затем нажмите enter(exit - завершение): ")
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if err := scanner.Err(); err != nil {
				fmt.Println("Непредвиденная ошибка:", err)
				os.Exit(1)
			}
			if scanner.Scan() {
				str = scanner.Text()
			}
			if strings.EqualFold(str, "exit") {
				fmt.Println("Продукции больше нет.")
				close(channel)
				return
			}
			data := strings.Split(str, " ")

			for _, v := range data {
				if n, err := strconv.Atoi(v); err == nil {
					channel <- n
					fmt.Println("нулевой цикл. Сырьё: ", n, " ушло в работу")
					time.Sleep(time.Second * 2)
				} else {
					fmt.Println("Бракованная продукция, только целые числа: ", v)
				}
			}
		}
	}()
	return channel, done
}
