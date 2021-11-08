package mocks

import (
	"net"
	"sync"
)

type MockPolysocket struct {
	SentMessages        []interface{}
	BroadcastedMessages []interface{}
	Channel             chan interface{}
	Connections         []net.Addr
	lock                sync.Mutex
}

func (pm *MockPolysocket) InjectMessage(data interface{}) {
	pm.Channel <- data
}

func (pm *MockPolysocket) Connect(addr net.Addr) (net.Conn, error) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.Connections = append(pm.Connections, addr)
	return &net.TCPConn{}, nil
}
func (pm *MockPolysocket) Close() []error {
	return nil
}
func (pm *MockPolysocket) Send(data interface{}, addr net.Addr) error {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.SentMessages = append(pm.SentMessages, data)
	return nil
}

func (pm *MockPolysocket) Broadcast(data interface{}) []error {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.BroadcastedMessages = append(pm.BroadcastedMessages, data)
	return nil
}

func (pm *MockPolysocket) GetAddr() net.Addr {
	return &net.TCPAddr{
		IP:   nil,
		Port: 0,
		Zone: "",
	}
}
