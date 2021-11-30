package node

import (
	"log"
	"testing"
	"vicoin/account"
	"vicoin/node"
)

func makeNode() *node.Node {
	ledger := account.NewLedger()
	n, err := node.NewNode(ledger)
	if err != nil {
		log.Fatalln(err)
	}
	return n
}
func TestNodesCanConnect(t *testing.T) {
	n1 := makeNode()
	n2 := makeNode()
	err := node.ConnectAndRequestPeers(n1, node.Self(n2))
	if err != nil {
		t.Error("Error when connecting nodes: ", err)
	}
	n1Peers := node.Peers(n1)
	if len(n1Peers) != 2 {
		t.Errorf("Peer not added!")
	}
}
func TestNodesAnnounceConnection(t *testing.T) {
	n1 := makeNode()
	n2 := makeNode()
	err := node.ConnectAndRequestPeers(n1, node.Self(n2))
	if err != nil {
		t.Error("Error when connecting nodes: ", err)
	}
	n2Peers := node.Peers(n2)
	if len(n2Peers) != 2 {
		t.Errorf("Peer not added!")
	}
}
func TestNodesPropagateConnectionAnnouncements(t *testing.T) {
	n1 := makeNode()
	n2 := makeNode()
	n3 := makeNode()
	node.ConnectAndRequestPeers(n2, node.Self(n1))
	node.ConnectAndRequestPeers(n3, node.Self(n1))
	n2Peers := node.Peers(n2)
	if len(n2Peers) != 3 {
		t.Errorf("Announcement not propagated from n1 to n2, got %d connections", len(n2Peers))
	}
}
func TestNodesAddTransactions(t *testing.T) {
	n1 := makeNode()
	self := node.Self(n1)
	self.Connect()
	transaction := &account.SignedTransaction{}
	self.SendTransaction(transaction)
	transactions := node.Transactions(n1)
	if len(transactions) != 1 {
		t.Errorf("Unexpected number of transactions %d, want 1", len(transactions))
	}
}
func TestNodesPropagateTransactions(t *testing.T) {
	n1 := makeNode()
	n2 := makeNode()
	node.ConnectAndRequestPeers(n1, node.Self(n2))
	self := node.Self(n1)
	self.Connect()
	transaction := &account.SignedTransaction{}
	self.SendTransaction(transaction)
	transactions := node.Transactions(n2)
	if len(transactions) != 1 {
		t.Errorf("Unexpected number of transactions %d, want 1", len(transactions))
	}
}
