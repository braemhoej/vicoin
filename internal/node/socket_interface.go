package node

import "net"

type SocketInterface interface {
	Connect(addr *net.TCPAddr) (net.Conn, error)
	Close() []error
	Broadcast(data interface{}) []error
	Send(data interface{}, addr *net.TCPAddr) error
	GetAddr() *net.TCPAddr
}
