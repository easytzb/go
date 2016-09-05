package main

import "fmt"

func xrange() chan int {
	ch := make(chan int)

	go func() {
		for i := 1; ; i++ {
			ch <- i
		}
	}()

	return ch
}

func main() {
	generator := xrange()

	for i := 0; i < 100; i++ {
		fmt.Println(<-generator)
	}
}
