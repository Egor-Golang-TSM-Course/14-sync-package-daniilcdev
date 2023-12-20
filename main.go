package main

import (
	"fmt"
	"lesson14/mtBankAccount"
	mtcounter "lesson14/mtCounter"
	mtlogger "lesson14/mtLogger"
)

func main() {
	done := make(chan struct{})

	go mtBankAccount.StartConcurrentChanges(done)

	<-done

	go mtlogger.StartConcurrentLogging(done)

	<-done

	go mtcounter.StartConcurrentVisits(done)

	<-done

	close(done)

	fmt.Println("exit")
}
