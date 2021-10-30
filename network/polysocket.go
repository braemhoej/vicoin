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
	listener    net.Listener
	connections map[string]net.Conn
	addr        *net.TCPAddr
	channel     chan interface{}
	lock        sync.Mutex
}

func NewPolysocket(internal chan interface{}) (polysocket *Polysocket, err error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	polysocket = &Polysocket{
		connections: make(map[string]net.Conn),
		addr:        nil,
		channel:     internal,
		lock:        sync.Mutex{},
	}
	polysocket.listener = listener
	polysocket.addr = polysocket.listener.Addr().(*net.TCPAddr)
	go polysocket.listen()
	return polysocket, err
}

func (polysocket *Polysocket) Connect(addr *net.TCPAddr) (net.Conn, error) {
	local, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	socket, err := net.DialTCP("tcp", local, addr)
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
	fmt.Print("Address: " + addr.String())
	socket := polysocket.connections[addr.String()]
	fmt.Print(socket)
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

func (polysocket *Polysocket) GetAddr() *net.TCPAddr {
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
