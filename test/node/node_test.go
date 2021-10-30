package node_test

import (
	"net"
	"testing"
	"vicoin/internal/account"
	"vicoin/internal/node"
	"vicoin/network"
)

func NewPolysocketMock(channel chan interface{}) *network.PolysocketMock {
	return &network.PolysocketMock{
		SentMessages:        make([]interface{}, 0),
		BroadcastedMessages: make([]interface{}, 0),
		Channel:             channel,
		Connections:         make([]*net.TCPAddr, 0),
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
