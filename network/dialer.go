package network

import "net"

type TCPDialer struct{}

func NewTCPDialer() (*TCPDialer, error) {
	return &TCPDialer{}, nil
}

func (dialer *TCPDialer) Dial(addr net.Addr) (net.Conn, error) {
	local, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	socket, err := net.DialTCP("tcp", local, addr.(*net.TCPAddr))
	return socket, err
}
