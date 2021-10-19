package network

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Link struct {
	Connections map[string]*net.Conn
	Addr        *net.TCPAddr
	channel     chan []byte
	lock        sync.Mutex
}

//NewLink (channel chan[]byte) (link *Link) :
//Creates a link which acts as both a server and a client.
//The link will listen at the TCP-address stored in link.Addr.
//The link can also send messages (of type []byte) directly to any connection in link.Connections,
//broadcast to all connections.
func NewLink(channel chan []byte) (link *Link) {
	link = &Link{
		Connections: make(map[string]*net.Conn),
		Addr:        nil,
		channel:     channel,
		lock:        sync.Mutex{},
	}
	link.Addr = link.listen()
	return link
}

//(link *Link) Broadcast(bytes []byte) :
//Broadcasts a byte array to all open connections in link.Connections.
//gob is used for encoding.
func (link *Link) Broadcast(bytes []byte) {
	link.lock.Lock()
	defer link.lock.Unlock()
	for _, conn := range link.Connections {
		enc := gob.NewEncoder(*conn)
		enc.Encode(bytes)
	}
}

//(link *Link) Connect(ip string, port string) (*net.Conn, error) :
//Establishes TCP-connection to server at ip:port.
//Returns a pointer to the connection, and error.
//Error is nil if connection was established successfully.
func (link *Link) Connect(ip string, port string) (*net.Conn, error) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Println("Connection failed", err)
	}
	if conn != nil {
		go link.handleConnection(conn)
		link.lock.Lock()
		link.Connections[conn.RemoteAddr().String()] = &conn
		link.lock.Unlock()
	}
	return &conn, err
}

//(link *Link) handleConnection(conn net.Conn)
//A method for continuously monitoring the connection conn
//for incomming messages. This link uses gob for marshalling,
//and expects gob-marshalled byte arrays. Received messages are sent on link.Channel.
//Prints to terminal if connection
//is closed.
func (link *Link) handleConnection(conn net.Conn) {
	defer conn.Close()
	var bytes []byte
	dec := gob.NewDecoder(conn)
	for {
		err := dec.Decode(&bytes)
		if err == io.EOF {
			fmt.Println("Connection closed by " + conn.RemoteAddr().String())
			return
		}
		if err != nil {
			log.Println(err.Error())
			return
		}
		link.channel <- bytes
	}
}

/*
	(link *Link) listen() *net.TCPAddr
		A method for starting a go-routine which continuously
		listens for incomming connections. Appends accepted
		connections to link.Connections. Returns a pointer to
		the address at which it is listening.
*/
func (link *Link) listen() *net.TCPAddr {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Panic(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Panic(err)
			}
			go link.handleConnection(conn)
			link.lock.Lock()
			link.Connections[conn.RemoteAddr().String()] = &conn
			link.lock.Unlock()
		}
	}()
	return listener.Addr().(*net.TCPAddr)
}
