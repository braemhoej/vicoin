package network

import (
	"encoding/gob"
	"net"
	"testing"
	"time"
	"vicoin/network"
	mocks "vicoin/test/mocks/network"
)

func makeMockDependencies() (*mocks.MockDialer, *mocks.MockListener) {
	dialer := mocks.NewMockDialer()
	listener := mocks.NewMockListener()
	return dialer, listener
}

func TestPolysocketsAreInitialisedWithZeroConnections(t *testing.T) {
	dialer, listener := makeMockDependencies()
	poly, _ := network.NewPolysocket(dialer, listener, 10)
	if len(poly.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 0", len(poly.GetConnections()))
	}
}
func TestPolysocketsAddConnectionsToListUponDialing(t *testing.T) {
	local, _ := net.Pipe()
	dialer, listener := makeMockDependencies()
	poly, _ := network.NewPolysocket(dialer, listener, 10)
	dialer.SetNextSocket(local)
	poly.Connect(&net.TCPAddr{})
	if len(poly.GetConnections()) != 1 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly.GetConnections()))
	}
}

func TestPolysocketsAddConnectionsToListUponReceivingConnection(t *testing.T) {
	local, _ := net.Pipe()
	dialer, listener := makeMockDependencies()
	poly, _ := network.NewPolysocket(dialer, listener, 10)
	listener.SetNextSocket(local)
	time.Sleep(5 * time.Millisecond)
	if len(poly.GetConnections()) != 1 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly.GetConnections()))
	}
}

func TestPolysocketsRemoveConnectionsUponRemoteDisconnection(t *testing.T) {
	local, remote := net.Pipe()
	dialer, listener := makeMockDependencies()
	poly, _ := network.NewPolysocket(dialer, listener, 10)
	listener.SetNextSocket(local)
	remote.Close()
	time.Sleep(5 * time.Millisecond)
	if len(poly.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly.GetConnections()))
	}
}

func TestPolysocketsPurgeConnectionsUponLocalDisconnection(t *testing.T) {
	local, _ := net.Pipe()
	dialer, listener := makeMockDependencies()
	poly, _ := network.NewPolysocket(dialer, listener, 10)
	listener.SetNextSocket(local)
	time.Sleep(5 * time.Millisecond)
	poly.Close()
	if len(poly.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly.GetConnections()))
	}
}

func TestPolysocketsDeliverReceivedMessagesOnChannel(t *testing.T) {
	local, remote := net.Pipe()
	dialer, listener := makeMockDependencies()
	_, channel := network.NewPolysocket(dialer, listener, 10)
	listener.SetNextSocket(local)
	time.Sleep(5 * time.Millisecond)
	enc := gob.NewEncoder(remote)
	var sent interface{} = "lorem ipsum" // Interface type is a consequence of implementation of decoding.
	enc.Encode(&sent)
	received := <-channel
	if received != sent {
		t.Error("Receive message doesn't equal sent message")
	}
}
