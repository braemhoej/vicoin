package node

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Node struct {
	Connections map[string]*net.Conn
	Addr        *net.TCPAddr
	channel     chan []byte
	lock        sync.Mutex
}

/*
	NewNode(channel chan[]byte) (node *Node) :
		Creates a node which acts as both a server and a client.
		The node will listen at the TCP-address stored in node.Addr.
		The node can also send messages (of type []byte) directly to any connection in node.Connections,
		or broadcast to all connections.
		channel :
*/
func NewNode(channel chan []byte) (node *Node) {
	node = &Node{
		Connections: make(map[string]*net.Conn),
		Addr:        nil,
		channel:     channel,
		lock:        sync.Mutex{},
	}
	node.Addr = node.listen()
	return node
}

/*
	(node *Node) Broadcast(bytes []byte) :
		Broadcasts a byte array to all open connections in node.Connections.
		gob is used for encoding.
*/
func (node *Node) Broadcast(bytes []byte) {
	node.lock.Lock()
	defer node.lock.Unlock()
	for _, conn := range node.Connections {
		enc := gob.NewEncoder(*conn)
		enc.Encode(bytes)
	}
}

/*
	(node *Node) Connect(ip string, port string) (*net.Conn, error) :
		Establishes TCP-connection to server at ip:port.
		Returns a pointer to the connection, and error.
		Error is nil if connection was established successfully.
*/
func (node *Node) Connect(ip string, port string) (*net.Conn, error) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Println("Connection failed", err)
	}
	if conn != nil {
		go node.handleConnection(conn)
		node.lock.Lock()
		node.Connections[conn.RemoteAddr().String()] = &conn
		node.lock.Unlock()
	}
	return &conn, err
}

/*
	(node *Node) handleConnection(conn net.Conn)
		A method for continuously monitoring the connection conn
		for incomming messages. This interface uses gob for marshalling,
		and expects gob-marshalled byte arrays. Received messages are sent on node.Channel.
		Prints to terminal if connection
		is closed.
*/
func (node *Node) handleConnection(conn net.Conn) {
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
		node.channel <- bytes
	}
}

/*
	(node *Node) listen() *net.TCPAddr
		A method for starting a go-routine which continuously
		listens for incomming connections. Appends accepted
		connections to node.Connections. Returns a pointer to
		the address at which it is listening.
*/
func (node *Node) listen() *net.TCPAddr {
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
			go node.handleConnection(conn)
			node.lock.Lock()
			node.Connections[conn.RemoteAddr().String()] = &conn
			node.lock.Unlock()
		}
	}()
	return listener.Addr().(*net.TCPAddr)
}
