package network

import (
	"testing"
	"time"
	"vicoin/internal/network"
)

func TestPolySocketsAreInitialisedWithZeroConnections(t *testing.T) {
	c1 := make(chan interface{})
	i1 := network.NewPolySocket(c1)
	if len(i1.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 1", len(i1.GetConnections()))
	}
	close(c1)
}
func TestPolySocketsAddConnectionsToListUponConnection(t *testing.T) {
	c1 := make(chan interface{})
	c2 := make(chan interface{})
	i1 := network.NewPolySocket(c1)
	i2 := network.NewPolySocket(c2)
	i2.Connect(network.TCP2Strings(i1.Addr))
	time.Sleep(5 * time.Millisecond) // Give the nodes a chance to update connection list.
	if len(i1.GetConnections()) != len(i2.GetConnections()) {
		t.Errorf("# of connections mismatched")
	}
	if len(i1.GetConnections()) != 1 {
		t.Errorf("Unexpected # of connections %d, want 1", len(i1.GetConnections()))
	}
}
func TestPolySocketsRemoveConnectionsFromListUponDisconnection(t *testing.T) {
	c1 := make(chan interface{})
	c2 := make(chan interface{})
	i1 := network.NewPolySocket(c1)
	i2 := network.NewPolySocket(c2)
	i2.Connect(network.TCP2Strings(i1.Addr))
	time.Sleep(5 * time.Millisecond) // Give the nodes a chance to update connection list.
	i2.Close()
	time.Sleep(5 * time.Millisecond) // Give the nodes a chance to update connection list.
	if len(i1.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 0", len(i1.GetConnections()))
	}
}
func TestPolySocketsBroadcastToAllConnections(t *testing.T) {
	c1 := make(chan interface{})
	c2 := make(chan interface{})
	c3 := make(chan interface{})
	i1 := network.NewPolySocket(c1)
	i2 := network.NewPolySocket(c2)
	i3 := network.NewPolySocket(c3)
	i2.Connect(network.TCP2Strings(i1.Addr))
	i3.Connect(network.TCP2Strings(i1.Addr))
	sent := "lorem ipsum"
	time.Sleep(50 * time.Millisecond)
	i1.Broadcast(sent)
	received := <-c2
	if received != sent {
		t.Errorf("Received (%d) message doesn't equal the sent (lorem ipsum) message", received)
	}
	received = <-c3
	if received != sent {
		t.Errorf("Received (%d) message doesn't equal the sent (lorem ipsum) message", received)
	}
}
func TestPolySocketsSendsOnlyToTarget(t *testing.T) {
	c1 := make(chan interface{})
	c2 := make(chan interface{})
	c3 := make(chan interface{})
	i1 := network.NewPolySocket(c1)
	i2 := network.NewPolySocket(c2)
	i3 := network.NewPolySocket(c3)
	i2.Connect(network.TCP2Strings(i1.Addr))
	i3.Connect(network.TCP2Strings(i1.Addr))
	sent := "lorem ipsum"
	time.Sleep(50 * time.Millisecond)
	i1.Broadcast(sent)
	received := <-c2
	if received != sent {
		t.Errorf("Received (%d) message doesn't equal the sent (lorem ipsum) message", received)
	}
	received = <-c3
	if received != sent {
		t.Errorf("Received (%d) message doesn't equal the sent (lorem ipsum) message", received)
	}
}
