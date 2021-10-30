package node

import (
	"net"
	"sync"
)

type PolysocketMock struct {
	SentMessages        []interface{}
	BroadcastedMessages []interface{}
	Channel             chan interface{}
	Connections         []*net.TCPAddr
	lock                sync.Mutex
}

func (pm *PolysocketMock) InjectMessage(data interface{}) {
	pm.Channel <- data
}

func (pm *PolysocketMock) Connect(addr *net.TCPAddr) (net.Conn, error) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.Connections = append(pm.Connections, addr)
	conn, _ := net.Pipe()
	return conn, nil
}
func (pm *PolysocketMock) Close() []error {
	return nil
}
func (pm *PolysocketMock) Send(data interface{}, addr *net.TCPAddr) error {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.SentMessages = append(pm.SentMessages, data)
	return nil
}

func (pm *PolysocketMock) Broadcast(data interface{}) error {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.BroadcastedMessages = append(pm.SentMessages, data)
	return nil
}

func (pm *PolysocketMock) GetAddr() *net.TCPAddr {
	return &net.TCPAddr{
		IP:   nil,
		Port: 0,
		Zone: "",
	}
}
