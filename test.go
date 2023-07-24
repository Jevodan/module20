package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// источник данных
	dataSource := func() <-chan int {
		c := make(chan int)
		done := make(chan bool)
		go func() {
			defer close(done)
			scanner := bufio.NewScanner(os.Stdin)
			var data string
			for {
				scanner.Scan()
				data = scanner.Text()
				if strings.EqualFold(data, "exit") {
					fmt.Println("Программа завершила работу!")
					close(c)
					return
				}
				i, err := strconv.Atoi(data)
				if err != nil {
					fmt.Println("Программа обрабатывает только целые числа!")
					continue
				}
				c <- i
			}
		}()
		return c
	}
	source := dataSource()
	for v := range source {
		fmt.Println(v)
	}

}
