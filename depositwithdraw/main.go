package main

import (
	"fmt"
	"os"
	"sync"
)

type bankOp struct {
	howMuch int
	confirm chan int
}

var accountBalance = 0
var bankRequests chan *bankOp

func updateBalance(amt int) int {
	update := &bankOp{howMuch: amt, confirm: make(chan int)}
	bankRequests <- update
	newBalance := <-update.confirm
	return newBalance
}

// For now a no-op, but could save balance to a file with a timestamp.
func logBalance(current int) {}

func reportAndExit(msg string) {
	fmt.Println(msg)
	os.Exit(-1)
}

func main() {
	iterations := 1000

	bankRequests = make(chan *bankOp, 8)

	var wg sync.WaitGroup

	go func() {
		for {

			select {
			case request := <-bankRequests:
				accountBalance += request.howMuch // update account
				request.confirm <- accountBalance // confirm with current balance
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			newBalance := updateBalance(1)
			logBalance(newBalance)
			//runtime.Gosched() // yield to another goroutine
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			newBalance := updateBalance(-1)
			logBalance(newBalance)
			//runtime.Gosched() // yield to another goroutine
		}
	}()

	wg.Wait()
	// confirm the balance is zero
	fmt.Println("Final balance: ", accountBalance)
}
