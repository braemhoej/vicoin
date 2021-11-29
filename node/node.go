package node

import (
	"log"
	"net"
	"sync"
	"vicoin/account"
	"vicoin/network"
)

type Node struct {
	peers    []Peer
	history  map[network.Packet]bool
	socket   network.Socket
	internal chan interface{}
	external chan account.SignedTransaction
	lock     sync.Mutex
}

func NewNode(polysocket network.Socket, internalChannel chan interface{}, externalChannel chan account.SignedTransaction) (*Node, error) {
	node := &Node{
		peers:    make([]Peer, 0),
		history:  make(map[network.Packet]bool),
		socket:   polysocket,
		internal: internalChannel,
		external: externalChannel,
		lock:     sync.Mutex{},
	}
	self := Peer{
		Addr: polysocket.GetAddr(),
	}
	node.peers = append(node.peers, self)
	go node.listen()
	return node, nil
}

func (node *Node) Connect(peer *Peer) error {
	conn, err := node.socket.Connect(peer.Addr)
	if err != nil {
		return err
	}
	peerRequest := network.Packet{
		Instruction: network.PeerRequest,
		Data:        conn.LocalAddr(),
	}
	node.socket.Send(peerRequest, conn.RemoteAddr())
	connAnnouncemet := network.Packet{
		Instruction: network.ConnAnn,
		Data:        node.socket.GetAddr(),
	}
	node.socket.Broadcast(connAnnouncemet)
	return nil
}

func (node *Node) Close() []error {
	return node.socket.Close()
}

func (node *Node) GetPeers() []Peer {
	return node.peers
}

func (node *Node) GetAddr() net.Addr {
	return node.socket.GetAddr()
}

func (node *Node) listen() {
	for {
		msg := <-node.internal
		switch packet := msg.(type) {
		case network.Packet:
			node.handle(packet)
		default:
			log.Println("Unexpected message type, skipping")
			continue
		}
	}
}

func (node *Node) handle(packet network.Packet) {
	switch packet.Instruction {
	// NOTE: Currently vulnerable to malformed packages, i.e. data not of expected type !!!!
	case network.PeerRequest:
		requester := packet.Data.(net.Addr)
		node.lock.Lock()
		reply := network.Packet{
			Instruction: network.PeerReply,
			Data:        node.peers,
		}
		node.socket.Send(reply, requester)
		node.lock.Unlock()
	case network.PeerReply:
		peers := packet.Data.([]Peer)
		node.lock.Lock()
		node.peers = merge(peers, node.peers)
		node.lock.Unlock()
		node.strengthenNetwork()
	case network.ConnAnn:
		peer := packet.Data.(Peer)
		node.lock.Lock()
		node.peers = append(node.peers, peer)
		node.socket.Broadcast(packet)
		node.lock.Unlock()
	case network.Transaction:
		signedTransaction := packet.Data.(account.SignedTransaction)
		node.socket.Broadcast(packet)
		node.external <- signedTransaction
	}
}

func (node *Node) SendTransaction(transaction account.SignedTransaction) {
	node.lock.Lock()
	defer node.lock.Unlock()
	wrappedTransaction := network.Packet{
		Instruction: network.Transaction,
		Data:        transaction,
	}
	node.socket.Broadcast(wrappedTransaction)
}

func (node *Node) strengthenNetwork() {
	node.lock.Lock()
	defer node.lock.Unlock()
	index := len(node.peers)
	if len(node.peers) > 11 {
		index = 11
	}
	for _, peer := range node.peers[len(node.peers)-index : len(node.peers)-1] {
		node.socket.Connect(peer.Addr)
	}
}

func contains(list []Peer, peer Peer) bool {
	for _, known := range list {
		if known == peer {
			return true
		}
	}
	return false
}

func merge(received []Peer, known []Peer) []Peer {
	output := known
	for _, addr := range received {
		if !contains(output, addr) {
			output = append(output, addr)
		}
	}
	return output
}
