package node

import (
	"fmt"
	"log"
	"net"
	"sync"
	"vicoin/account"
	"vicoin/network"
)

type Node struct {
	peers    []Peer
	history  map[account.SignedTransaction]bool
	socket   network.Socket
	incoming chan interface{}
	internal chan Packet
	lock     sync.Mutex
}

func NewNode(polysocket network.Socket, incomingChannel chan interface{}) (*Node, error) {
	node := &Node{
		peers:    make([]Peer, 0),
		history:  make(map[account.SignedTransaction]bool),
		socket:   polysocket,
		incoming: incomingChannel,
		internal: make(chan Packet),
		lock:     sync.Mutex{},
	}
	self := Peer{
		Addr: polysocket.GetAddr(),
	}
	node.peers = append(node.peers, self)
	go node.listen()
	go node.handlePackets()
	return node, nil
}

func (node *Node) Connect(peer *Peer) error {
	conn, err := node.socket.Connect(peer.Addr)
	if err != nil {
		return err
	}
	peerRequest := Packet{
		Instruction: PeerRequest,
		Data:        conn.LocalAddr(),
	}
	node.socket.Send(peerRequest, conn.RemoteAddr())
	connAnnouncemet := Packet{
		Instruction: ConnAnn,
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

// Listen for messages received on the internal channel.
// Ensure received message is a packet through type assertion,
// and ensure that it is not malformed. If succesful, pass the
// packet along to the message handling method. Otherwise,
// skip message.
func (node *Node) listen() {
	for {
		fmt.Println("waiting for packaet")
		msg := <-node.incoming
		switch packet := msg.(type) {
		case Packet:
			if isMalformed(packet) {
				log.Println("Packet malformed, skipping")
				continue
			} else {
				node.internal <- packet
			}
		default:
			log.Println("Unexpected message type, skipping")
			continue
		}
	}
}

// Continuously listen for new packets on the internal channel.
// When a packet is received, check if it is deliverable. If it
// is, handle the packet, and handle all other undelivered packets.
// If it is not, append it to the list of undelivered packets.
func (node *Node) handlePackets() {
	undeliveredPackets := make([]Packet, 0)
	for {
		packet := <-node.internal
		fmt.Println("Message received")
		if deliverable(packet) {
			node.handle(packet)
			undeliveredPackets = node.handleUndelivered(undeliveredPackets)
		} else {
			undeliveredPackets = append(undeliveredPackets, packet)
		}
	}
}

func (node *Node) handleUndelivered(undeliveredPackets []Packet) []Packet {
	messageDelivered := false
	remainingPackets := undeliveredPackets
	for index, packet := range undeliveredPackets {
		if deliverable(packet) {
			node.handle(packet)
			remainingPackets = append(remainingPackets[0:index-1], remainingPackets[index+1:]...)
			messageDelivered = true
		}
	}
	if messageDelivered {
		return node.handleUndelivered(remainingPackets)
	} else {
		return remainingPackets
	}
}

// Checks is a packet is malformed. Returns true if the type of the data is
// not as expected given the instruction, or if the attatched vector clock is
// nil.
func isMalformed(packet Packet) bool {
	return false // TODO: Implement!
}

// Checks, if given the maintained vector clock, the packet is deliverable.
// Returns true if the vector clock attatched to the packet indicates that all
// previous messages have been delivered.
func deliverable(packet Packet) bool {
	return true // TODO: Implement!
}

// Decodes and handles the passed network.Packet according to the instruction.
// Assumes passed packet is deliverable.
func (node *Node) handle(packet Packet) {
	switch packet.Instruction {
	// NOTE: Currently vulnerable to malformed packages, i.e. data not of expected type !!!!
	case PeerRequest:
		requester := packet.Data.(net.Addr)
		node.handlePeerRequest(requester)
	case PeerReply:
		peers := packet.Data.([]Peer)
		node.handlePeerReply(peers)
	case ConnAnn:
		fmt.Println("ConnAnn")
		peer := packet.Data.(Peer)
		node.handleConnectionAnnouncement(peer, packet)
	case Transaction:
		fmt.Println("Transaction")
		signedTransaction := packet.Data.(account.SignedTransaction)
		node.handleTransaction(signedTransaction, packet)
	default:
		log.Printf("Unknown instruction %d \n", packet.Instruction)
	}
}

// Handles a peer request by sending a copy of this nodes peers
// to the address given as argument. Returns a error if no connection is
// established to the given address.
func (node *Node) handlePeerRequest(requester net.Addr) error {
	node.lock.Lock()
	defer node.lock.Unlock()
	reply := Packet{
		Instruction: PeerReply,
		Data:        node.peers,
	}
	return node.socket.Send(reply, requester)
}

// Handles a peer reply by merging the received list of peers with the
// the list of known peers maintained by the node. Strengthens the network
// by connecting to the last ten peers in the merged list.
func (node *Node) handlePeerReply(peers []Peer) {
	node.lock.Lock()
	defer node.lock.Unlock()
	for _, addr := range peers {
		if !contains(node.peers, addr) {
			node.peers = append(node.peers, addr)
		}
	}
	// Strengthen network TODO: Insert check if # connections is below some value, i.e. 10.
	index := len(node.peers)
	if len(node.peers) > 11 {
		index = 11
	}
	for _, peer := range node.peers[len(node.peers)-index : len(node.peers)-1] {
		node.socket.Connect(peer.Addr)
	}
}

// Handles a connection announcement by appending the received peer to
// the list of peers maintained by the node. Then propagate the announcenent
// by broadcasting the received packet to all open connections.
func (node *Node) handleConnectionAnnouncement(peer Peer, packet Packet) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.peers = append(node.peers, peer) //TODO: Ensure that known peers are not re-added.
	node.socket.Broadcast(packet)         //TODO: Strengthen network if # connections below some threshold.
}

// Handle a transaction by attempting to update the ledger. Then propagate the
// transaction by broadcasting the received packet to all open connections.
func (node *Node) handleTransaction(transaction account.SignedTransaction, packet Packet) {
	node.lock.Lock()
	defer node.lock.Unlock()
	// TODO: Handle transaction
	node.socket.Broadcast(packet)
}

func (node *Node) SendTransaction(transaction account.SignedTransaction) {
	node.lock.Lock()
	defer node.lock.Unlock()
	wrappedTransaction := Packet{
		Instruction: Transaction,
		Data:        transaction,
	}
	node.socket.Broadcast(wrappedTransaction)
}

func contains(list []Peer, peer Peer) bool {
	for _, known := range list {
		if known == peer {
			return true
		}
	}
	return false
}
