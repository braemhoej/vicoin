package network_test

import (
	"net"
	"testing"
	"time"
	"vicoin/network"
)

func makeDependencies() (chan interface{}, network.DialerStrategy, network.ListenerStrategy) {
	dialer, _ := network.NewTCPDialer()
	listener, _ := network.NewTCPListener()
	return make(chan interface{}), dialer, listener
}

func TestPolysocketsAreInitialisedWithZeroConnections(t *testing.T) {
	channel, dialer, listener := makeDependencies()
	poly := network.NewPolysocket(channel, dialer, listener)
	if len(poly.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly.GetConnections()))
	}
	close(channel)
}
func TestPolysocketsAddConnectionsToListUponConnection(t *testing.T) {
	channel1, dialer1, listener1 := makeDependencies()
	poly1 := network.NewPolysocket(channel1, dialer1, listener1)
	channel2, dialer2, listener2 := makeDependencies()
	poly2 := network.NewPolysocket(channel2, dialer2, listener2)
	_, err := poly2.Connect(poly1.GetAddr())
	if err != nil {
		t.Error("Error connection : ", err)
	}
	time.Sleep(5 * time.Millisecond) // Give the nodes a chance to update connection list.
	if len(poly1.GetConnections()) != len(poly2.GetConnections()) {
		t.Errorf("# of connections mismatched")
	}
	if len(poly1.GetConnections()) != 1 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly1.GetConnections()))
	}
}
func TestPolysocketsRemoveConnectionsFromListUponDisconnection(t *testing.T) {
	channel1, dialer1, listener1 := makeDependencies()
	poly1 := network.NewPolysocket(channel1, dialer1, listener1)
	channel2, dialer2, listener2 := makeDependencies()
	poly2 := network.NewPolysocket(channel2, dialer2, listener2)
	poly2.Connect(poly1.GetAddr())
	time.Sleep(5 * time.Millisecond) // Give the nodes a chance to update connection list.
	poly2.Close()
	time.Sleep(5 * time.Millisecond) // Give the nodes a chance to update connection list.
	if len(poly1.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 0", len(poly1.GetConnections()))
	}
}
func TestPolysocketsBroadcastToAllConnections(t *testing.T) {
	channel1, dialer1, listener1 := makeDependencies()
	poly1 := network.NewPolysocket(channel1, dialer1, listener1)
	channel2, dialer2, listener2 := makeDependencies()
	poly2 := network.NewPolysocket(channel2, dialer2, listener2)
	channel3, dialer3, listener3 := makeDependencies()
	poly3 := network.NewPolysocket(channel3, dialer3, listener3)
	poly2.Connect(poly1.GetAddr())
	poly3.Connect(poly1.GetAddr())
	sent := "lorem ipsum"
	time.Sleep(50 * time.Millisecond)
	poly1.Broadcast(sent)
	received := <-channel2
	if received != sent {
		t.Errorf("Received (%d) message doesn't equal the sent (lorem ipsum) message", received)
	}
	received = <-channel3
	if received != sent {
		t.Errorf("Received (%d) message doesn't equal the sent (lorem ipsum) message", received)
	}
}
func TestPolysocketsSendsOnlyToTarget(t *testing.T) {
	channel1, dialer1, listener1 := makeDependencies()
	poly1 := network.NewPolysocket(channel1, dialer1, listener1)
	channel2, dialer2, listener2 := makeDependencies()
	poly2 := network.NewPolysocket(channel2, dialer2, listener2)
	channel3, dialer3, listener3 := makeDependencies()
	poly3 := network.NewPolysocket(channel3, dialer3, listener3)
	conn2, _ := poly2.Connect(poly1.GetAddr())
	conn3, _ := poly3.Connect(poly1.GetAddr())
	sent := "lorem ipsum"
	secret := "ipsum lorem"
	time.Sleep(50 * time.Millisecond)
	poly1.Send(sent, conn2.LocalAddr().(*net.TCPAddr))
	poly1.Send(secret, conn3.LocalAddr().(*net.TCPAddr))
	received := <-channel2
	if received != sent {
		t.Errorf("Received (%d) message doesn't equal the sent (lorem ipsum) message", received)
	}
	received = <-channel3
	if received != secret {
		t.Errorf("Received (%d) message doesn't equal the sent (ipsum lorem) message", received)
	}
}
