package main

import (
	"fmt"
	"time"
)

func main() {
	channel := make(chan int)

	go func(msg int) {
		time.Sleep(time.Second)
		channel <- msg
		fmt.Println("middle")
	}(1)

	fmt.Println("begin")
	fmt.Println(<-channel)
	fmt.Println("end")
}
