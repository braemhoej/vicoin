package network

import (
	"net"
	"reflect"
	"strconv"
	"testing"
	"time"
	node "vicoin/internal/network"
)

func TCPToStrings(addr *net.TCPAddr) (ip string, port string) {
	return "[" + addr.IP.String() + "]", strconv.Itoa(addr.Port)
}

func TestNodesCanEstablishConnections(t *testing.T) {
	c1 := make(chan []byte)
	c2 := make(chan []byte)
	n1 := node.NewNode(c1)
	n2 := node.NewNode(c2)
	_, err := n2.Connect(TCPToStrings(n1.Addr))
	time.Sleep(5 * time.Millisecond)
	if err != nil {
		t.Error("Connection Error")
	}
	if len(n1.Connections) != len(n2.Connections) {
		t.Error("Number of connections mismatched")
	}
	if len(n1.Connections) != 1 {
		t.Errorf("Number of connections %d, want 1", len(n1.Connections))
	}
}

func TestBroadcastMessagesAreSentToAllConnections(t *testing.T) {
	c1 := make(chan []byte)
	c2 := make(chan []byte)
	c3 := make(chan []byte)
	n1 := node.NewNode(c1)
	n2 := node.NewNode(c2)
	n3 := node.NewNode(c3)
	_, err1 := n2.Connect(TCPToStrings(n1.Addr))
	_, err2 := n3.Connect(TCPToStrings(n1.Addr))
	time.Sleep(5 * time.Millisecond)
	if err1 != nil || err2 != nil {
		t.Error("Connection Error")
	}
	msg := []byte("lorem ipsum")
	n1.Broadcast(msg)
	rmsg := <-c2
	if !reflect.DeepEqual(rmsg, msg) {
		t.Error("Received message doesn't match sent message: ", rmsg, msg)
	}
	rmsg = <-c3
	if !reflect.DeepEqual(rmsg, msg) {
		t.Error("Received message doesn't match sent message: ", rmsg, msg)
	}
}
