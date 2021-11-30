package network

import (
	"errors"
	"net"
	"net/rpc"
)

func makeTCPListener() (*net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

type RPCServer interface {
	RegisterName(name string, object interface{}) error
	Close()
	Addr() net.Addr
}

type RPCServerTCP struct {
	listener *net.TCPListener
	close    chan bool
}

func NewRPCServerTCP() (server *RPCServerTCP, err error) {
	listener, err := makeTCPListener()
	if err != nil {
		return nil, err
	}
	server = &RPCServerTCP{
		listener: listener,
		close:    make(chan bool),
	}
	go server.listen()
	return server, err
}

func (server *RPCServerTCP) RegisterName(name string, object interface{}) error {
	return rpc.RegisterName(name, object)
}

func (server *RPCServerTCP) Close() {
	server.close <- true
}

func (server *RPCServerTCP) Addr() net.Addr {
	return server.listener.Addr()
}

func (server *RPCServerTCP) listen() error {
	if server.listener == nil {
		return errors.New("listener is nil")
	}
	for {
		select {
		case <-server.close:
			return nil
		default:
			go rpc.Accept(server.listener)
		}
	}
}
