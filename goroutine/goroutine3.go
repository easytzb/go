package main

import (
	"fmt"
	"time"
)

var ch chan int

func foo(id int) {
	fmt.Printf("%d", id)
	time.Sleep(time.Second)
	ch <- 0
}

func main() {
	count := 100
	ch = make(chan int, count)
	for i := 0; i < count; i++ {
		go foo(i + 1)
	}

	for i := 0; i < count; i++ {
		<-ch
	}
}
