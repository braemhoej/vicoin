package network

import (
	"encoding/gob"
	"net"
	"testing"
	"time"
	"vicoin/network"
	mocks "vicoin/test/mocks/network"
)

func makeMockDependencies() (chan interface{}, *mocks.MockDialer, *mocks.MockListener) {
	dialer := mocks.NewMockDialer()
	listener := mocks.NewMockListener()
	return make(chan interface{}), dialer, listener
}

func TestPolysocketsAreInitialisedWithZeroConnections(t *testing.T) {
	channel, dialer, listener := makeMockDependencies()
	poly := network.NewPolysocket(channel, dialer, listener)
	if len(poly.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 0", len(poly.GetConnections()))
	}
}
func TestPolysocketsAddConnectionsToListUponDialing(t *testing.T) {
	local, _ := net.Pipe()
	channel, dialer, listener := makeMockDependencies()
	dialer.SetNextSocket(local)
	poly := network.NewPolysocket(channel, dialer, listener)
	poly.Connect(&net.TCPAddr{})
	if len(poly.GetConnections()) != 1 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly.GetConnections()))
	}
}

func TestPolysocketsAddConnectionsToListUponReceivingConnection(t *testing.T) {
	local, _ := net.Pipe()
	channel, dialer, listener := makeMockDependencies()
	poly := network.NewPolysocket(channel, dialer, listener)
	listener.SetNextSocket(local)
	time.Sleep(5 * time.Millisecond)
	if len(poly.GetConnections()) != 1 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly.GetConnections()))
	}
}

func TestPolysocketsRemoveConnectionsUponRemoteDisconnection(t *testing.T) {
	local, remote := net.Pipe()
	channel, dialer, listener := makeMockDependencies()
	poly := network.NewPolysocket(channel, dialer, listener)
	listener.SetNextSocket(local)
	remote.Close()
	time.Sleep(5 * time.Millisecond)
	if len(poly.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly.GetConnections()))
	}
}

func TestPolysocketsPurgeConnectionsUponLocalDisconnection(t *testing.T) {
	local, _ := net.Pipe()
	channel, dialer, listener := makeMockDependencies()
	poly := network.NewPolysocket(channel, dialer, listener)
	listener.SetNextSocket(local)
	time.Sleep(5 * time.Millisecond)
	poly.Close()
	if len(poly.GetConnections()) != 0 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly.GetConnections()))
	}
}

func TestPolysocketsDeliverReceivedMessagesOnChannel(t *testing.T) {
	local, remote := net.Pipe()
	channel, dialer, listener := makeMockDependencies()
	network.NewPolysocket(channel, dialer, listener)
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
