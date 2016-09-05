package main

import (
	"fmt"
)

var ch1 chan int = make(chan int)
var ch2 chan int = make(chan int)

func say(s int) {
	fmt.Println(s)
	ch1 <- <-ch2
}

func main() {
	go say(1)
	<-ch1
}
