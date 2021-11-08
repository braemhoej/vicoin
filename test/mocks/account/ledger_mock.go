package mocks

import (
	"sync"
	"vicoin/internal/account"
)

type MockLedger struct {
	Transactions []*account.SignedTransaction
	lock         sync.Mutex
}

func NewMockLedger() *MockLedger {
	return &MockLedger{
		Transactions: make([]*account.SignedTransaction, 0),
		lock:         sync.Mutex{},
	}
}

func (mock *MockLedger) SignedTransaction(transaction *account.SignedTransaction) error {
	mock.lock.Lock()
	defer mock.lock.Unlock()
	mock.Transactions = append(mock.Transactions, transaction)
	return nil
}

func (mock *MockLedger) GetBalance(account string) float64 {
	return 42
}
