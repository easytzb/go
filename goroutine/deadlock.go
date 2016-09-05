package main

import (
//"fmt"
)

func main() {
	channel := make(chan int)
	channel <- 0
	//or
	//fmt.Println(<-channel)
}
