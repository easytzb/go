package main

import (
	"fmt"
	"time"
)

func loop(cnt int) {
	for i := 0; i < cnt; i++ {
		fmt.Printf("%d", i)
	}
}

func main() {
	go loop(10)
	loop(10)
	time.Sleep(time.Second)
}
