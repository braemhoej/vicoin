package network

import "net"

type Socket interface {
	Connect(addr net.Addr) (net.Conn, error)
	Close() []error
	Broadcast(data interface{}) []error
	Send(data interface{}, addr *net.TCPAddr) error
	GetAddr() *net.TCPAddr
}

type DialerStrategy interface {
	Dial(net.Addr) (net.Conn, error)
}

type ListenerStrategy interface {
	Accept() (net.Conn, error)
	Addr() net.Addr
}
