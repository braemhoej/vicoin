package network

import (
	"encoding/gob"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

func TCP2Strings(addr *net.TCPAddr) (string, string) {
	return "[" + addr.IP.String() + "]", strconv.Itoa(addr.Port)
}

type PolySocket struct {
	connections map[string]net.Conn
	Addr        *net.TCPAddr
	channel     chan interface{}
	lock        sync.Mutex
}

func NewPolySocket(internal chan interface{}) (polySocket *PolySocket) {
	polySocket = &PolySocket{
		connections: make(map[string]net.Conn),
		Addr:        nil,
		channel:     internal,
		lock:        sync.Mutex{},
	}
	polySocket.Addr = polySocket.listen()
	return polySocket
}

func (polySocket *PolySocket) Connect(ip string, port string) net.Conn {
	socket, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Println("Connection refused: ", err)
		return nil
	}
	polySocket.lock.Lock()
	go polySocket.handle(socket)
	polySocket.connections[socket.RemoteAddr().String()] = socket
	polySocket.lock.Unlock()
	return socket
}

func (polySocket *PolySocket) Close() {
	polySocket.lock.Lock()
	defer polySocket.lock.Unlock()
	for _, socket := range polySocket.connections {
		err := socket.Close()
		if err != nil {
			log.Println("Error when closing polySocket: ", err)
		}
	}
}

func (polySocket *PolySocket) Broadcast(data interface{}) {
	polySocket.lock.Lock()
	defer polySocket.lock.Unlock()
	for _, socket := range polySocket.connections {
		polySocket.Send(data, socket)
	}
}

func (polySocket *PolySocket) Send(data interface{}, socket net.Conn) {
	enc := gob.NewEncoder(socket)
	err := enc.Encode(&data)
	if err != nil {
		log.Panicln("Error sending package: ", err)
	}
}

func (polySocket *PolySocket) GetConnections() map[string]net.Conn {
	polySocket.lock.Lock()
	defer polySocket.lock.Unlock()
	return polySocket.connections
}

// Internal
func (polySocket *PolySocket) listen() *net.TCPAddr {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalln("Fatal error when starting listener: ", err)
	}
	go func() {
		for {
			socket, err := listener.Accept()
			if err != nil {
				log.Panicln("Incoming connection dropped: ", err)
			}
			log.Println("Incoming connection accepted: ", socket.RemoteAddr().String())
			polySocket.lock.Lock()
			go polySocket.handle(socket)
			polySocket.connections[socket.RemoteAddr().String()] = socket
			polySocket.lock.Unlock()
		}
	}()
	return listener.Addr().(*net.TCPAddr)
}

func (polySocket *PolySocket) handle(socket net.Conn) {
	defer socket.Close()
	var buffer interface{}
	dec := gob.NewDecoder(socket)
	for {
		err := dec.Decode(&buffer)
		if err == io.EOF {
			log.Println("Connection closed by " + socket.RemoteAddr().String())
			polySocket.lock.Lock()
			delete(polySocket.connections, socket.RemoteAddr().String())
			polySocket.lock.Unlock()
			return
		}
		if err != nil {
			log.Println("Error when decoding: ", err.Error())
			break
		}
		polySocket.channel <- buffer
	}
}
