package main

import (
	"sync"
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

func TestConcurrentDeposits(t *testing.T) {
	acc := NewBankAccount()

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			acc.Deposit(1)
		}()
	}
	wg.Wait()

	if acc.Balance() != 1000 {
		t.Errorf("expected balance 1000, got %d", acc.Balance())
	}
}
