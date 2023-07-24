package main

import (
	"bufio"
	"fmt"
	"module20/konveyor"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	BUFFER_SIZE               = 50               // Буфер кольцевого массива
	INTERVAL    time.Duration = time.Second * 10 // Интервал времени, за который происходит очистка кольцевого буфера
)

func main() {
	// источник данных
	source, done := sourceData()
	// создание конввера
	konv := konveyor.NewKonveer(done, filterOne, filterTwo, filterThree)
	// запуск конвеера
	truba := konv.Start(source)

	for v := range truba {
		time.Sleep(time.Second * 2)
		fmt.Println("Конечная продукция: ", v)
	}
	time.Sleep(time.Second * 2)
}

// Источник данных
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
