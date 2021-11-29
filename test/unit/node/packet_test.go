package node_test

import (
	"net"
	"testing"
	"vicoin/node"
)

func TestVectorClockIncrementCreatesMapEntryForUnknownKeys(t *testing.T) {
	peer := node.Peer{
		Addr: &net.IPAddr{
			IP: net.ParseIP("192.168.0.1"),
		},
	}
	vectorClock := node.NewVectorClock()
	vectorClock.Increment(peer)
	clock := vectorClock.Get(peer)
	if clock != 1 {
		t.Errorf("Unexpected clock value %d expected 1", clock)
	}
}

func TestVectorClockGetReturnsZeroForUnknownKeys(t *testing.T) {
	peer := node.Peer{
		Addr: &net.IPAddr{
			IP: net.ParseIP("192.168.0.1"),
		},
	}
	vectorClock := node.NewVectorClock()
	clock := vectorClock.Get(peer)
	if clock != 0 {
		t.Errorf("Unexpected clock value %d expected 0", clock)
	}
}

func TestVectorClockIncrementOnlyTarget(t *testing.T) {
	peer1 := node.Peer{
		Addr: &net.IPAddr{
			IP: net.ParseIP("192.168.0.1"),
		},
	}
	peer2 := node.Peer{
		Addr: &net.IPAddr{
			IP: net.ParseIP("192.168.0.2"),
		},
	}
	vectorClock := node.NewVectorClock()
	vectorClock.Increment(peer1)
	clock1 := vectorClock.Get(peer1)
	if clock1 != 1 {
		t.Errorf("Unexpected clock1 value %d expected 1", clock1)
	}
	clock2 := vectorClock.Get(peer2)
	if clock2 != 0 {
		t.Errorf("Unexpected clock2 value %d expected 0", clock2)
	}
}
