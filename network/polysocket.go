package network

import (
	"encoding/gob"
	"fmt"
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
	outgoing    chan interface{}
	close       chan bool
	lock        sync.Mutex
}

func NewPolysocket(dialerStrategy DialerStrategy, listenerStrategy ListenerStrategy, outgoingBuffer int) (polysocket *Polysocket, outgoing chan interface{}) {
	polysocket = &Polysocket{
		listener:    listenerStrategy,
		dialer:      dialerStrategy,
		connections: make(map[string]net.Conn),
		addr:        nil,
		outgoing:    make(chan interface{}, outgoingBuffer),
		close:       make(chan bool, 1),
		lock:        sync.Mutex{},
	}
	polysocket.addr = polysocket.listener.Addr()
	go polysocket.listen()
	return polysocket, polysocket.outgoing
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
	polysocket.close <- true
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
		fmt.Println("Writing")
		enc := gob.NewEncoder(socket)
		err := enc.Encode(&data)
		if err != nil {
			errors = append(errors, err)
		}
	}
	fmt.Println("Done")
	return errors
}

func (polysocket *Polysocket) Send(data interface{}, addr net.Addr) error {
	polysocket.lock.Lock()
	defer polysocket.lock.Unlock()
	socket := polysocket.connections[addr.String()]
	fmt.Println(socket.RemoteAddr())
	enc := gob.NewEncoder(socket)
	err := enc.Encode(&data)
	if err != nil {
		return err
	}
	return nil
}

func (polysocket *Polysocket) GetAddr() net.Addr {
	return polysocket.addr
}

func (polysocket *Polysocket) GetConnections() map[string]net.Conn {
	return polysocket.connections
}

// Internal
func (polysocket *Polysocket) listen() {
	for {
		select {
		case <-polysocket.close:
			return
		default:
			socket, err := polysocket.listener.Accept()
			if err != nil {
				log.Println("Incoming net.Conn dropped: ", err)
			}
			log.Println("Incoming net.Conn accepted: ", socket.RemoteAddr().String())
			polysocket.lock.Lock()
			go polysocket.handle(socket)
			polysocket.connections[socket.RemoteAddr().String()] = socket
			polysocket.lock.Unlock()
		}
	}
}

func (polysocket *Polysocket) handle(socket net.Conn) {
	// defer socket.Close()
	for {
		select {
		case <-polysocket.close:
			return
		default:
			packet, err := polysocket.read(socket)
			if err == io.EOF {
				log.Println("net.Conn closed by " + socket.RemoteAddr().String())
				polysocket.lock.Lock()
				delete(polysocket.connections, socket.RemoteAddr().String())
				polysocket.lock.Unlock()
				return
			} else if err != nil {
				log.Println("Error when decoding connection: ", err.Error())
			} else if packet != nil {
				polysocket.outgoing <- packet
			}
		}
	}
}

func (polysocket *Polysocket) read(socket net.Conn) (interface{}, error) {
	var buffer interface{}
	dec := gob.NewDecoder(socket)
	err := dec.Decode(&buffer)
	return buffer, err
}
