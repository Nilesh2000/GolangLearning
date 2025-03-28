package main

import (
	"errors"
	"sync"
)

type BankAccount struct {
	balance int
	mu      sync.Mutex
}

func NewBankAccount() *BankAccount {
	return &BankAccount{}
}

func (b *BankAccount) Deposit(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.balance += amount
}

func (b *BankAccount) Withdraw(amount int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if amount > b.balance {
		return errors.New("Insufficient balance")
	}
	b.balance -= amount
	return nil
}

func (b *BankAccount) Balance() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.balance
}
