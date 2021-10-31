package mocks

import (
	"net"
)

type MockListener struct {
	internal chan net.Conn
}

func NewMockListener() *MockListener {
	return &MockListener{
		internal: make(chan net.Conn),
	}
}

func (listener *MockListener) Accept() (net.Conn, error) {
	conn := <-listener.internal
	return conn, nil
}

func (listener *MockListener) Addr() net.Addr {
	return &net.IPAddr{}
}

func (listener *MockListener) SetNextSocket(socket net.Conn) {
	go func() {
		listener.internal <- socket
	}()
}
