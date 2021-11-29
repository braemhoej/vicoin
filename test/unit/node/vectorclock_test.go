package node_test

import (
	"net"
	"strconv"
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

func TestVectorClocksUpdateByTakingMax(t *testing.T) {
	peers := make([]node.Peer, 0)
	vc1 := node.NewVectorClock()
	vc2 := node.NewVectorClock()
	for i := 1; i < 10; i++ {
		peer := node.Peer{
			Addr: &net.IPAddr{
				IP: []byte("mock" + strconv.Itoa(i)),
			},
		}
		peers = append(peers, peer)
	}
	for index, peer := range peers {
		if index%2 == 0 {
			vc1.Increment(peer)
			vc1.Increment(peer)
		} else {
			vc2.Increment(peer)
		}
	}
	updatedClock := node.Update(vc1, vc2)
	for index, peer := range peers {
		time := updatedClock.Get(peer)
		if index%2 == 0 {
			if time != 2 {
				t.Errorf("Unexpected time %d want 2", time)
			}
		} else {
			if time != 1 {
				t.Errorf("Unexpected time %d want 1", time)
			}
		}
	}
}
