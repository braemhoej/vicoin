package network

import "net"

type TCPListener struct {
	ln net.Listener
}

func NewTCPListener() (*TCPListener, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	return &TCPListener{
		ln: listener,
	}, nil
}

func (listener *TCPListener) Addr() net.Addr {
	return listener.ln.Addr().(*net.TCPAddr)
}

func (listener *TCPListener) Accept() (net.Conn, error) {
	socket, err := listener.ln.Accept()
	return socket, err
}
