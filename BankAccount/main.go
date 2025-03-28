package main

import "errors"

type BankAccount struct {
	balance int
}

func NewBankAccount() *BankAccount {
	return &BankAccount{}
}

func (b *BankAccount) Deposit(amount int) {
	b.balance += amount
}

func (b *BankAccount) Withdraw(amount int) error {
	if amount > b.balance {
		return errors.New("Insufficient balance")
	}
	b.balance -= amount
	return nil
}

func (b *BankAccount) Balance() int {
	return b.balance
}
