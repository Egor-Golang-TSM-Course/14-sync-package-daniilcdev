package main

import (
	"fmt"
	"lesson14/mtBankAccount"
)

func main() {
	done := make(chan struct{})

	go mtBankAccount.StartConcurrentChanges(done)

	<-done

	fmt.Println("exit")
}
