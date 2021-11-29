package mocks

import (
	"net"
	"sync"
	"vicoin/account"
)

type MockNode struct {
	sent []*account.SignedTransaction
	lock sync.Mutex
}

func NewMockNode() *MockNode {
	return &MockNode{
		sent: make([]*account.SignedTransaction, 0),
		lock: sync.Mutex{},
	}
}

func (mock *MockNode) Connect(addr net.Addr) error {
	return nil
}

func (mock *MockNode) Close() []error {
	return nil
}

func (mock *MockNode) SendTransaction(transaction account.SignedTransaction) {
	mock.lock.Lock()
	defer mock.lock.Unlock()
	mock.sent = append(mock.sent, &transaction)
}

func (mock *MockNode) GetAddr() net.Addr {
	return &net.TCPAddr{}
}
