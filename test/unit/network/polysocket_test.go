package network

import (
	"net"
	"testing"
	"time"
	"vicoin/network"
	"vicoin/test/mocks"
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
	client, _ := net.Pipe()
	channel, dialer, listener := makeMockDependencies()
	dialer.SetNextSocket(client)
	poly := network.NewPolysocket(channel, dialer, listener)
	poly.Connect(&net.TCPAddr{})
	if len(poly.GetConnections()) != 1 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly.GetConnections()))
	}
}

func TestPolysocketsAddConnectionsToListUponReceivingConnection(t *testing.T) {
	client, _ := net.Pipe()
	channel, dialer, listener := makeMockDependencies()
	poly := network.NewPolysocket(channel, dialer, listener)
	listener.SetNextSocket(client)
	time.Sleep(5 * time.Millisecond)
	if len(poly.GetConnections()) != 1 {
		t.Errorf("Unexpected # of connections %d, want 1", len(poly.GetConnections()))
	}
}
