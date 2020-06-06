package main

import (
	"fmt"
	"runtime"
	"time"
)

func pinger(c chan string) {
	for i := 0; ; i++ {
		fmt.Printf("pinger loop %d\n", i)
		c <- fmt.Sprintf("ping%d", i)
		runtime.Gosched() // yield to another goroutine
	}
}

func ponger(c chan string) {
	for i := 0; ; i++ {
		fmt.Printf("ponger loop %d\n", i)
		c <- fmt.Sprintf("pong%d", i)
		runtime.Gosched() // yield to another goroutine
	}
}

func printer(c chan string) {
	for {
		msg := <-c
		fmt.Println(msg)
		time.Sleep(time.Second * 5)
	}
}

func main() {
	var c chan string = make(chan string, 5)

	go pinger(c)
	go ponger(c)
	go printer(c)
	select {}

}
