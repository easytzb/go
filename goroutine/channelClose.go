package main

import (
	"fmt"
)

func main() {
	//no deadlock
	ch := make(chan int, 2)

	ch <- 1
	ch <- 2

	//if don't close  channel, deadlock will
	close(ch)

	//range will read until channel is close.
	for v := range ch {
		fmt.Println(v)
	}

}
