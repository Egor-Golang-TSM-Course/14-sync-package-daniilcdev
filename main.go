package main

import (
	"fmt"
	mtbroker "lesson14/mtBroker"
)

func main() {
	done := make(chan struct{})

	// go mtBankAccount.StartConcurrentChanges(done)

	// <-done

	// go mtlogger.StartConcurrentLogging(done)

	// <-done

	// go mtcounter.StartConcurrentVisits(done)

	// <-done

	go mtbroker.StartRequestBroking(done)

	<-done

	close(done)

	fmt.Println("exit")
}
