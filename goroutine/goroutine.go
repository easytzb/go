package main

import (
	"fmt"
)

var quit chan int = make(chan int)

func loop(cnt int) {
	for i := 0; i < cnt; i++ {
		fmt.Printf("%d\n", i)
	}
	quit <- 0
}

func main() {
	go loop(10)
	go loop(10)
	for i := 0; i < 2; i++ {
		fmt.Println(<-quit)
	}
}
