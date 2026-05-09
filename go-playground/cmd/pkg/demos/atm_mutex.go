package demos

import (
	"fmt"
	"sync"
	"time"
)

// ATM struct represents a simple bank account
type ATM struct {
	balance int        // the balance
	mu      sync.Mutex // mutex to protect the balance
}

// Withdraw method safely allows withdrawing from ATM
func (a *ATM) Withdraw(amount int, person string) {
	a.mu.Lock()
	defer a.mu.Unlock() // ensures unlock even if something goes wrong

	time.Sleep(500 * time.Millisecond) // simulate time delay

	if a.balance >= amount {
		a.balance -= amount
		fmt.Printf("%s successfully withdrew %d. Balance: %d\n", person, amount, a.balance)
	} else {
		fmt.Printf("%s tried to withdraw %d, but not enough funds. Balance: %d\n", person, amount, a.balance)
	}
}

// Deposit method allows depositing money into the ATM
func (a *ATM) Deposit(amount int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
	fmt.Printf("Deposited %d. Current balance: %d\n", amount, a.balance)
}

// ATM Mutex Demo
func ATMMutexDemo() {
	fmt.Println("---> ATM Mutex Demo")
	atm := &ATM{balance: 1000}

	// Simulate multiple people trying to withdraw money at the same time
	var wg sync.WaitGroup

	atm.Deposit(500)

	wg.Add(3) // 3 people trying to withdraw concurrently

	go func() {
		defer wg.Done()
		atm.Withdraw(400, "Person A")
	}()

	go func() {
		defer wg.Done()
		atm.Withdraw(300, "Person B")
	}()

	go func() {
		defer wg.Done()
		atm.Withdraw(200, "Person C")
	}()

	wg.Wait()
}
