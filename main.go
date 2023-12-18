package main

import (
	"fmt"
	"lesson14/mtBankAccount"
	mtlogger "lesson14/mtLogger"
)

func main() {
	done := make(chan struct{})

	go mtBankAccount.StartConcurrentChanges(done)

	<-done

	go mtlogger.StartConcurrentLogging(done)

	<-done

	fmt.Println("exit")
}
