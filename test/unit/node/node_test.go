package node_test

import (
	"net"
	"strconv"
	"testing"
	"time"
	"vicoin/account"
	"vicoin/network"
	"vicoin/node"
	mocks "vicoin/test/mocks/network"
)

func makeDependecies() (*mocks.MockPolysocket, chan interface{}) {
	internal := make(chan interface{})
	return &mocks.MockPolysocket{
		SentMessages:        make([]interface{}, 0),
		BroadcastedMessages: make([]interface{}, 0),
		Channel:             internal,
		Connections:         make([]net.Addr, 0),
	}, internal
}

func TestNewNodeReturnsAPointerToANewNode(t *testing.T) {
	n, err := node.NewNode(makeDependecies())
	if err != nil {
		t.Error("Error when creating node: ", err)
	}
	if n == nil {
		t.Error("Nil returned")
	}
}

func TestNewNodesAddOwnAddressToPeerList(t *testing.T) {
	n, _ := node.NewNode(makeDependecies())
	peers := n.GetPeers()
	if len(peers) != 1 {
		t.Errorf("Unexpected number of peers %d, want 1", len(peers))
	}
}

func TestNewNodesSendPeerRequestUponConnection(t *testing.T) {
	mock, internal := makeDependecies()
	n, _ := node.NewNode(mock, internal)
	n.Connect(&node.Peer{}) // Mock address
	if len(mock.SentMessages) != 1 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.SentMessages[0].(network.Packet)
	if msg.Instruction != network.PeerRequest {
		t.Errorf("Unexpected instruction %d, want 0", msg.Instruction)
	}
}

func TestNodesSendPeerReplyUponRequest(t *testing.T) {
	mock, internal := makeDependecies()
	node.NewNode(mock, internal)
	mock.InjectMessage(network.Packet{
		Instruction: network.PeerRequest,
		Data:        &net.TCPAddr{},
	})
	time.Sleep(50 * time.Millisecond)
	if len(mock.SentMessages) != 1 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.SentMessages[0].(network.Packet)
	if msg.Instruction != network.PeerReply {
		t.Errorf("Unexpected instruction %d, want 1", msg.Instruction)
	}
}

func TestNewNodesMergePeerRequestIntoKnownPeers(t *testing.T) {
	mock, internal := makeDependecies()
	n, _ := node.NewNode(mock, internal)
	for i := 1; i < 10; i++ {
		mock.InjectMessage(network.Packet{
			Instruction: network.ConnAnn,
			Data: node.Peer{
				Addr: &net.IPAddr{
					IP: []byte("mock" + strconv.Itoa(i)),
				},
			},
		})
	}
	peers := n.GetPeers()
	if len(peers) != 10 {
		t.Errorf("Unexpected number of peers %d, want 10", len(peers))
	}
}

func TestNodesBroadcastConnectionAnnouncementUponConnection(t *testing.T) {
	mock, internal := makeDependecies()
	n, _ := node.NewNode(mock, internal)
	n.Connect(&node.Peer{}) // Mock address
	time.Sleep(50 * time.Millisecond)
	if len(mock.BroadcastedMessages) != 1 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.BroadcastedMessages[0].(network.Packet)
	if msg.Instruction != network.ConnAnn {
		t.Errorf("Unexpected instruction %d, want 3", msg.Instruction)
	}
}
func TestNewNodesAddAnnouncedConnectionsToPeerList(t *testing.T) {
	mock, internal := makeDependecies()
	n, _ := node.NewNode(mock, internal)
	for i := 1; i < 10; i++ {
		mock.InjectMessage(network.Packet{
			Instruction: network.ConnAnn,
			Data: node.Peer{
				Addr: &net.IPAddr{
					IP: []byte("mock" + strconv.Itoa(i)),
				},
			},
		})
	}
	peers := n.GetPeers()
	if len(peers) != 10 {
		t.Errorf("Unexpected number of peers %d, want 10", len(peers))
	}
}
func TestNodesPropagateConnectionAnnouncements(t *testing.T) {
	mock, internal := makeDependecies()
	n, _ := node.NewNode(mock, internal)
	n.Connect(&node.Peer{}) // Mock address
	mock.InjectMessage(network.Packet{Instruction: network.ConnAnn, Data: node.Peer{}})
	time.Sleep(50 * time.Millisecond)
	if len(mock.BroadcastedMessages) != 2 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.BroadcastedMessages[1].(network.Packet)
	if msg.Instruction != network.ConnAnn {
		t.Errorf("Unexpected instruction %d, want 3", msg.Instruction)
	}
}

func TestNodesPropagateTransactions(t *testing.T) {
	mock, internal := makeDependecies()
	n, _ := node.NewNode(mock, internal)
	n.Connect(&node.Peer{}) // Mock address
	mock.InjectMessage(network.Packet{Instruction: network.Transaction, Data: account.SignedTransaction{}})
	time.Sleep(50 * time.Millisecond)
	if len(mock.BroadcastedMessages) != 2 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.BroadcastedMessages[1].(network.Packet)
	if msg.Instruction != network.Transaction {
		t.Errorf("Unexpected instruction %d, want 4", msg.Instruction)
	}
}
