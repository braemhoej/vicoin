package network

import (
	"testing"
	"time"
	"vicoin/network"
)

func TestPolysocketsAreInitialisedWithZeroConnections(t *testing.T) {
	c1 := make(chan interface{})
	i1, err := network.NewPolysocket(c1)
	if err != nil {
		t.Error("Error creating polysocket: ", err)
	}
	if len(i1.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 1", len(i1.GetConnections()))
	}
	close(c1)
}
func TestPolysocketsAddConnectionsToListUponConnection(t *testing.T) {
	c1 := make(chan interface{})
	c2 := make(chan interface{})
	i1, err1 := network.NewPolysocket(c1)
	i2, err2 := network.NewPolysocket(c2)
	if err1 != nil || err2 != nil {
		t.Error("Error creating polysocket: ", err1, err2)
	}
	_, err := i2.Connect(network.TCP2Strings(i1.GetAddr()))
	if err != nil {
		t.Error("Error creating polysocket: ", err1, err2)
	}
	time.Sleep(5 * time.Millisecond) // Give the nodes a chance to update connection list.
	if len(i1.GetConnections()) != len(i2.GetConnections()) {
		t.Errorf("# of connections mismatched")
	}
	if len(i1.GetConnections()) != 1 {
		t.Errorf("Unexpected # of connections %d, want 1", len(i1.GetConnections()))
	}
}
func TestPolysocketsRemoveConnectionsFromListUponDisconnection(t *testing.T) {
	c1 := make(chan interface{})
	c2 := make(chan interface{})
	i1, _ := network.NewPolysocket(c1)
	i2, _ := network.NewPolysocket(c2)
	i2.Connect(network.TCP2Strings(i1.GetAddr()))
	time.Sleep(5 * time.Millisecond) // Give the nodes a chance to update connection list.
	i2.Close()
	time.Sleep(5 * time.Millisecond) // Give the nodes a chance to update connection list.
	if len(i1.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 0", len(i1.GetConnections()))
	}
}
func TestPolysocketsBroadcastToAllConnections(t *testing.T) {
	c1 := make(chan interface{})
	c2 := make(chan interface{})
	c3 := make(chan interface{})
	i1, _ := network.NewPolysocket(c1)
	i2, _ := network.NewPolysocket(c2)
	i3, _ := network.NewPolysocket(c3)
	i2.Connect(network.TCP2Strings(i1.GetAddr()))
	i3.Connect(network.TCP2Strings(i1.GetAddr()))
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
func TestPolysocketsSendsOnlyToTarget(t *testing.T) {
	c1 := make(chan interface{})
	c2 := make(chan interface{})
	c3 := make(chan interface{})
	i1, _ := network.NewPolysocket(c1)
	i2, _ := network.NewPolysocket(c2)
	i3, _ := network.NewPolysocket(c3)
	i2.Connect(network.TCP2Strings(i1.GetAddr()))
	i3.Connect(network.TCP2Strings(i1.GetAddr()))
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
