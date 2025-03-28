package main

import (
	"testing"
)

func TestBankAccount(t *testing.T) {
	t.Run("deposit increases balance", func(t *testing.T) {
		acc := NewBankAccount()
		acc.Deposit(100)
		if acc.Balance() != 100 {
			t.Errorf("got %d, want 100", acc.Balance())
		}
	})

	t.Run("withdraw reduces balance", func(t *testing.T) {
		acc := NewBankAccount()
		acc.Deposit(100)
		err := acc.Withdraw(50)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if acc.Balance() != 50 {
			t.Errorf("got %d, want 50", acc.Balance())
		}
	})

	t.Run("withdraw with insufficient balance returns error", func(t *testing.T) {
		acc := NewBankAccount()
		acc.Deposit(100)
		err := acc.Withdraw(120)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
