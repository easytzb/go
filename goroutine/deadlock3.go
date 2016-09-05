package main

import (
	"fmt"
)

var ch1 chan int = make(chan int)
var ch2 chan int = make(chan int)

func say(s int) {
	fmt.Println(s)

	//no deadlock
	ch2 <- 1
	ch1 <- 0

	//deadlock
	//ch1 <- 1
	//ch2 <- 0
}

func main() {
	go say(1)
	<-ch2
}
