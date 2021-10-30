package network

import (
	"encoding/gob"
	"io"
	"log"
	"net"
	"sync"
)

type Polysocket struct {
	listener    ListenerStrategy
	dialer      DialerStrategy
	connections map[string]net.Conn
	addr        net.Addr
	channel     chan interface{}
	lock        sync.Mutex
}

func NewPolysocket(internal chan interface{}, dialerStrategy DialerStrategy, listenerStrategy ListenerStrategy) (polysocket *Polysocket) {
	polysocket = &Polysocket{
		listener:    listenerStrategy,
		dialer:      dialerStrategy,
		connections: make(map[string]net.Conn),
		addr:        nil,
		channel:     internal,
		lock:        sync.Mutex{},
	}
	polysocket.addr = polysocket.listener.Addr()
	go polysocket.listen()
	return polysocket
}

func (polysocket *Polysocket) Connect(addr net.Addr) (net.Conn, error) {
	socket, err := polysocket.dialer.Dial(addr)
	if err != nil {
		return nil, err
	}
	polysocket.lock.Lock()
	go polysocket.handle(socket)
	polysocket.connections[socket.RemoteAddr().String()] = socket
	polysocket.lock.Unlock()
	return socket, nil
}

func (polysocket *Polysocket) Close() []error {
	polysocket.lock.Lock()
	defer polysocket.lock.Unlock()
	var errors []error
	for _, socket := range polysocket.connections {
		err := socket.Close()
		if err != nil {
			errors = append(errors, err)
		}
	}
	polysocket.connections = make(map[string]net.Conn)
	return errors
}

func (polysocket *Polysocket) Broadcast(data interface{}) []error {
	polysocket.lock.Lock()
	defer polysocket.lock.Unlock()
	var errors []error
	for _, socket := range polysocket.connections {
		enc := gob.NewEncoder(socket)
		err := enc.Encode(&data)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (polysocket *Polysocket) Send(data interface{}, addr *net.TCPAddr) error {
	polysocket.lock.Lock()
	defer polysocket.lock.Unlock()
	socket := polysocket.connections[addr.String()]
	enc := gob.NewEncoder(socket)
	err := enc.Encode(&data)
	if err != nil {
		return err
	}
	return nil
}

func (polysocket *Polysocket) GetConnections() map[string]net.Conn {
	polysocket.lock.Lock()
	defer polysocket.lock.Unlock()
	return polysocket.connections
}

func (polysocket *Polysocket) GetAddr() net.Addr {
	return polysocket.addr
}

// Internal
func (polysocket *Polysocket) listen() {
	for {
		socket, err := polysocket.listener.Accept()
		if err != nil {
			log.Println("Incoming connection dropped: ", err)
		}
		log.Println("Incoming connection accepted: ", socket.RemoteAddr().String())
		polysocket.lock.Lock()
		go polysocket.handle(socket)
		polysocket.connections[socket.RemoteAddr().String()] = socket
		polysocket.lock.Unlock()
	}
}

func (polysocket *Polysocket) handle(socket net.Conn) {
	defer socket.Close()
	var buffer interface{}
	dec := gob.NewDecoder(socket)
	for {
		err := dec.Decode(&buffer)
		if err == io.EOF {
			log.Println("Connection closed by " + socket.RemoteAddr().String())
			polysocket.lock.Lock()
			delete(polysocket.connections, socket.RemoteAddr().String())
			polysocket.lock.Unlock()
			return
		}
		if err != nil {
			log.Println("Error when decoding: ", err.Error())
			break
		}
		polysocket.channel <- buffer
	}
}
