package node

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"vicoin/account"
	"vicoin/network"
)

type Node struct {
	self         *Peer
	peers        []*Peer
	connections  []*Peer
	transactions map[account.SignedTransaction]bool
	ledger       account.LedgerInterface
	rpcServer    network.RPCServer
	lock         sync.Mutex
}

func NewNode(ledger account.LedgerInterface) (*Node, error) {
	server, err := network.NewRPCServerTCP()
	if err != nil {
		return nil, err
	}
	UUID := "node(" + server.Addr().(*net.TCPAddr).String() + ")" // TODO: Implement some UUID algorithm ... ?
	self := NewPeer(server.Addr().(*net.TCPAddr), UUID)
	peers := append(make([]*Peer, 0), self)
	node := &Node{
		self:         self,
		peers:        peers,
		connections:  make([]*Peer, 0),
		transactions: make(map[account.SignedTransaction]bool),
		ledger:       ledger,
		rpcServer:    server,
		lock:         sync.Mutex{},
	}
	server.RegisterName(UUID, node)
	return node, nil
}

// Establishes a new connection to a Peer RPC-server.
// Appends peer to node.connections if the given peer
// is not already contained in this list.
func Connect(node *Node, peer *Peer) error {
	node.lock.Lock()
	if !contains(node.connections, peer) {
		node.connections = append(node.connections, peer)
	}
	if !contains(node.peers, peer) {
		node.peers = append(node.peers, peer)
	}
	err := peer.Connect()
	if err != nil {
		return err
	}
	node.lock.Unlock()
	peer.AddConnection(node.self)
	return nil
}

// Establishes a new connection to a Peer RPC-server by calling
// node.Connect(peer). Furthemore a peer list is requested from
// the given peer through a RPC-Call. The received peer list is
// then merged into the list of peers maintained by the node.
func ConnectAndRequestPeers(node *Node, peer *Peer) error {
	err := Connect(node, peer)
	if err != nil {
		fmt.Println("Error: Connecting")
		return err
	}
	node.lock.Lock()
	peers, err := peer.RequestPeers()
	if err != nil {
		fmt.Println("Error: Requesting Peers")
		return err
	}
	node.peers = merge(node.peers, peers)
	node.lock.Unlock()
	peer.AnnounceConnection(node.self)
	return nil
}

func Self(node *Node) *Peer {
	return node.self
}

func Peers(node *Node) []*Peer {
	return node.peers
}

func Transactions(node *Node) map[account.SignedTransaction]bool {
	return node.transactions
}

/* RPC methods */

func (node *Node) AnnounceConnection(peer *Peer, _ *struct{}) error {
	if contains(node.peers, peer) {
		return errors.New("connection announcement replay")
	}
	node.lock.Lock()
	node.peers = append(node.peers, peer)
	node.lock.Unlock()
	connections := node.connections
	for _, connection := range connections {
		connection.AnnounceConnection(peer)
	}
	return nil
}

func (node *Node) AddConnection(peer *Peer, _ *struct{}) error {
	node.lock.Lock()
	defer node.lock.Unlock()
	if contains(node.connections, peer) {
		return errors.New("connection announcement replay")
	}
	peer.Connect()
	node.connections = append(node.connections, peer)
	return nil
}

func (node *Node) AddTransaction(transaction *account.SignedTransaction, _ *struct{}) error {
	node.lock.Lock()
	if _, contains := node.transactions[*transaction]; contains {
		return errors.New("transaction replay")
	}
	node.transactions[*transaction] = true
	connections := node.connections
	node.lock.Unlock()
	for _, peer := range connections {
		peer.SendTransaction(transaction)
	}
	return nil
}

func (node *Node) GetPeers(_ *struct{}, reply *[]*Peer) error {
	node.lock.Lock()
	defer node.lock.Unlock()
	*reply = node.peers
	return nil
}

// Returns true, used for checking liveness of RPC connections.
func (node *Node) Ping(_ *struct{}, reply *bool) error {
	*reply = true
	return nil
}

// Aux. functions.

// Merges listA []*Peer with listB []*Peer by appending
// every peer in listB which is not contained in listA
// to listA. Somewhat slow, could be improved.
func merge(listA []*Peer, listB []*Peer) (output []*Peer) {
	output = listA
	for _, peer := range listB {
		if !contains(output, peer) {
			output = append(output, peer)
		}
	}
	return output
}

// Checks if list []*Peer contains peer *Peer, by
// iterating through the list O(n).
func contains(list []*Peer, peer *Peer) bool {
	for _, known := range list {
		if known.UUID == peer.UUID {
			return true
		}
	}
	return false
}
