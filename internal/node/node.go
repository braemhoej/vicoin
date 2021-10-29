package node

import (
	"log"
	"net"
	"sync"
	"vicoin/internal/account"
	"vicoin/network"
)

type Node struct {
	peers    []Peer
	history  map[account.SignedTransaction]bool
	socket   *network.Polysocket
	internal chan interface{}
	external chan account.SignedTransaction
	lock     sync.Mutex
}

func NewNode() (*Node, error) {
	internalChannel := make(chan interface{})
	externalChannel := make(chan account.SignedTransaction)
	polysocket, err := network.NewPolysocket(internalChannel)
	if err != nil {
		return nil, err
	}
	node := &Node{
		peers:    make([]Peer, 0),
		history:  make(map[account.SignedTransaction]bool),
		socket:   polysocket,
		internal: internalChannel,
		external: externalChannel,
		lock:     sync.Mutex{},
	}
	go node.handle()
	return node, nil
}

func (node *Node) handle() {
	for {
		msg := <-node.internal
		switch packet := msg.(type) {
		case network.Packet:
			switch packet.Instruction {
			case network.PeerRequest:
				requester := packet.Data.(net.TCPAddr)
				node.lock.Lock()
				node.socket.Send(node.peers, requester)
				node.lock.Unlock()
			case network.PeerReply:
				peers := packet.Data.([]Peer)
				node.lock.Lock()
				node.peers = merge(peers, node.peers)
				node.lock.Unlock()
			case network.ConnAnnouncment:
				peer := packet.Data.(Peer)
				node.lock.Lock()
				node.peers = append(node.peers, peer)
				node.lock.Unlock()
			case network.Transaction:
				signedTransaction := packet.Data.(account.SignedTransaction)
				node.external <- signedTransaction
			}
		default:
			log.Println("Unexpected message type, skipping")
			continue
		}
	}
}

func merge(received []Peer, known []Peer) []Peer {
	//TODO: Do merge
	return nil
}
