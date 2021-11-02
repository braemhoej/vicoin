package node

import (
	"log"
	"net"
	"sync"
	"vicoin/internal/account"
	"vicoin/network"
)

type Node struct {
	peers    map[Peer]bool
	history  map[network.Packet]bool
	socket   network.Socket
	internal chan interface{}
	external chan account.SignedTransaction
	lock     sync.Mutex
}

func NewNode(polysocket network.Socket, internalChannel chan interface{}, externalChannel chan account.SignedTransaction) (*Node, error) {
	node := &Node{
		peers:    make(map[Peer]bool),
		history:  make(map[network.Packet]bool),
		socket:   polysocket,
		internal: internalChannel,
		external: externalChannel,
		lock:     sync.Mutex{},
	}
	self := Peer{
		Addr: polysocket.GetAddr(),
	}
	node.peers[self] = true
	go node.handle()
	return node, nil
}

func (node *Node) Connect(addr net.Addr) error {
	go node.handle()
	conn, err := node.socket.Connect(addr)
	if err != nil {
		return err
	}
	peerRequest := network.Packet{
		Instruction: network.PeerRequest,
		Data:        conn.LocalAddr(),
	}
	node.socket.Send(peerRequest, conn.RemoteAddr())
	connAnnouncemet := network.Packet{
		Instruction: network.ConnAnnouncment,
		Data:        node.socket.GetAddr(),
	}
	node.socket.Broadcast(connAnnouncemet)
	return nil
}

func (node *Node) Close() []error {
	return node.socket.Close()
}

func (node *Node) GetPeers() map[Peer]bool {
	return node.peers
}

func (node *Node) handle() {
	for {
		msg := <-node.internal
		switch packet := msg.(type) {
		case network.Packet:
			if seen := node.history[msg.(network.Packet)]; seen {
				continue
			}
			switch packet.Instruction {
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
				peers := packet.Data.(map[Peer]bool)
				node.lock.Lock()
				node.peers = merge(peers, node.peers)
				node.lock.Unlock()
			case network.ConnAnnouncment:
				peer := packet.Data.(Peer)
				node.lock.Lock()
				node.peers[peer] = true
				node.socket.Broadcast(msg)
				node.lock.Unlock()
			case network.Transaction:
				signedTransaction := packet.Data.(account.SignedTransaction)
				node.socket.Broadcast(msg)
				node.external <- signedTransaction
			}
		default:
			log.Println("Unexpected message type, skipping")
			continue
		}
	}
}

func merge(received map[Peer]bool, known map[Peer]bool) map[Peer]bool {
	for peer := range received {
		known[peer] = true
	}
	return known
}
