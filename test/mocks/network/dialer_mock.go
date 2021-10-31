package mocks

import (
	"net"
)

type MockDialer struct {
	internal chan net.Conn
}

func NewMockDialer() *MockDialer {
	return &MockDialer{
		internal: make(chan net.Conn),
	}
}

func (dialer *MockDialer) Dial(net.Addr) (net.Conn, error) {
	return <-dialer.internal, nil
}

func (dialer *MockDialer) SetNextSocket(socket net.Conn) {
	go func() {
		dialer.internal <- socket
	}()
}
