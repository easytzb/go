package main

import (
	"fmt"
)

func main() {
	//no deadlock
	ch := make(chan int, 2)

	ch <- 1
	ch <- 2

	fmt.Println(123)
}
