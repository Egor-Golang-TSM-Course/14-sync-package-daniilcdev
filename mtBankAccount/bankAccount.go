package mtBankAccount

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type BankAccount struct {
	balance int
	mutex   sync.Mutex
}

func StartConcurrentChanges(done chan<- struct{}) {
	ba := newBankAccount()
	wg := sync.WaitGroup{}
	for i := int64(0); i < 100; i++ {
		r := rand.New(rand.NewSource(i))
		wg.Add(1)

		id := i
		go func(ba *BankAccount) {
			time.Sleep(time.Duration(100+(r.Int63()%900)) * time.Millisecond)

			s := r.Int63()
			if s%100 < 50 {
				fmt.Printf("worker-%d - deposit 100\n", id)
				ba.Deposit(100)
			} else {
				fmt.Printf("worker-%d - withdraw 30\n", id)
				ba.Withdraw(30)
			}

			wg.Done()
		}(ba)
	}

	wg.Wait()

	fmt.Printf("final balance: %d\n", ba.balance)

	done <- struct{}{}
}

func newBankAccount() *BankAccount {
	return &BankAccount{
		mutex: sync.Mutex{},
	}
}

func (ba *BankAccount) Deposit(amount int) {
	defer ba.mutex.Unlock()

	ba.mutex.Lock()
	ba.balance += amount
}

func (ba *BankAccount) Withdraw(amount int) {
	defer ba.mutex.Unlock()

	ba.mutex.Lock()
	if ba.balance >= amount {
		ba.balance -= amount
	} else {
		fmt.Printf("can't withdraw, insufficient balance: %d\n", ba.balance)
	}
}
