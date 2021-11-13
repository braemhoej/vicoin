package node_test

import (
	"net"
	"strconv"
	"testing"
	"time"
	"vicoin/internal/account"
	"vicoin/internal/node"
	"vicoin/network"
	mocks "vicoin/test/mocks/network"
)

func NewPolysocketMock(channel chan interface{}) *mocks.MockPolysocket {
	return &mocks.MockPolysocket{
		SentMessages:        make([]interface{}, 0),
		BroadcastedMessages: make([]interface{}, 0),
		Channel:             channel,
		Connections:         make([]net.Addr, 0),
	}
}

func TestNewNodeReturnsAPointerToANewNode(t *testing.T) {
	internal := make(chan interface{})
	external := make(chan account.SignedTransaction)
	node, err := node.NewNode(NewPolysocketMock(internal), internal, external)
	if err != nil {
		t.Error("Error when creating node: ", err)
	}
	if node == nil {
		t.Error("Nil returned")
	}
}

func TestNewNodesAddOwnAddressToPeerList(t *testing.T) {
	internal := make(chan interface{})
	external := make(chan account.SignedTransaction)
	node, _ := node.NewNode(NewPolysocketMock(internal), internal, external)
	peers := node.GetPeers()
	if len(peers) != 1 {
		t.Errorf("Unexpected number of peers %d, want 1", len(peers))
	}
}

func TestNewNodesSendPeerRequestUponConnection(t *testing.T) {
	internal := make(chan interface{})
	external := make(chan account.SignedTransaction)
	mock := NewPolysocketMock(internal)
	node, _ := node.NewNode(mock, internal, external)
	node.Connect(&net.TCPAddr{}) // Mock address
	if len(mock.SentMessages) != 1 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.SentMessages[0].(network.Packet)
	if msg.Instruction != network.PeerRequest {
		t.Errorf("Unexpected instruction %d, want 0", msg.Instruction)
	}
}

func TestNodesSendPeerReplyUponRequest(t *testing.T) {
	internal := make(chan interface{})
	external := make(chan account.SignedTransaction)
	mock := NewPolysocketMock(internal)
	node.NewNode(mock, internal, external)
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
	internal := make(chan interface{})
	external := make(chan account.SignedTransaction)
	mock := NewPolysocketMock(internal)
	n, _ := node.NewNode(mock, internal, external)
	for i := 1; i < 10; i++ {
		mock.InjectMessage(network.Packet{
			Instruction: network.ConnAnnouncment,
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
	internal := make(chan interface{})
	external := make(chan account.SignedTransaction)
	mock := NewPolysocketMock(internal)
	node, _ := node.NewNode(mock, internal, external)
	node.Connect(&net.TCPAddr{}) // Mock address
	time.Sleep(50 * time.Millisecond)
	if len(mock.BroadcastedMessages) != 1 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.BroadcastedMessages[0].(network.Packet)
	if msg.Instruction != network.ConnAnnouncment {
		t.Errorf("Unexpected instruction %d, want 3", msg.Instruction)
	}
}
func TestNewNodesAddAnnouncedConnectionsToPeerList(t *testing.T) {
	internal := make(chan interface{})
	external := make(chan account.SignedTransaction)
	mock := NewPolysocketMock(internal)
	n, _ := node.NewNode(mock, internal, external)
	for i := 1; i < 10; i++ {
		mock.InjectMessage(network.Packet{
			Instruction: network.ConnAnnouncment,
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
	internal := make(chan interface{})
	external := make(chan account.SignedTransaction)
	mock := NewPolysocketMock(internal)
	n, _ := node.NewNode(mock, internal, external)
	n.Connect(&net.TCPAddr{}) // Mock address
	mock.InjectMessage(network.Packet{Instruction: network.ConnAnnouncment, Data: node.Peer{}})
	time.Sleep(50 * time.Millisecond)
	if len(mock.BroadcastedMessages) != 2 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.BroadcastedMessages[1].(network.Packet)
	if msg.Instruction != network.ConnAnnouncment {
		t.Errorf("Unexpected instruction %d, want 3", msg.Instruction)
	}
}

func TestNodesSendSignedTransactionsOnChannel(t *testing.T) {
	internal := make(chan interface{})
	external := make(chan account.SignedTransaction)
	mock := NewPolysocketMock(internal)
	node, _ := node.NewNode(mock, internal, external)
	node.Connect(&net.TCPAddr{}) // Mock address
	mock.InjectMessage(network.Packet{Instruction: network.Transaction, Data: account.SignedTransaction{}})
	time.Sleep(50 * time.Millisecond)
	<-external
}

func TestNodesPropagateTransactions(t *testing.T) {
	internal := make(chan interface{})
	external := make(chan account.SignedTransaction)
	mock := NewPolysocketMock(internal)
	n, _ := node.NewNode(mock, internal, external)
	n.Connect(&net.TCPAddr{}) // Mock address
	mock.InjectMessage(network.Packet{Instruction: network.Transaction, Data: account.SignedTransaction{}})
	time.Sleep(50 * time.Millisecond)
	<-external
	if len(mock.BroadcastedMessages) != 2 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.BroadcastedMessages[1].(network.Packet)
	if msg.Instruction != network.Transaction {
		t.Errorf("Unexpected instruction %d, want 4", msg.Instruction)
	}
}
