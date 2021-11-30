package node

import (
	"fmt"
	"log"
	"net"
	"sync"
	"vicoin/account"
	"vicoin/crypto"
	"vicoin/network"
)

type Node struct {
	peers              []Peer
	history            map[Packet]bool
	transactions       map[string]account.SignedTransaction
	socket             network.Socket
	incoming           chan interface{}
	internal           chan account.SignedTransaction
	ledger             account.LedgerInterface
	sequencer          bool
	public             *crypto.PublicKey
	private            *crypto.PrivateKey
	latestBlock        int
	blockSize          int
	transactionCounter int
	lock               sync.Mutex
}

func NewNode(polysocket network.Socket, incomingChannel chan interface{}, ledger account.LedgerInterface) (*Node, error) {
	node := &Node{
		peers:        make([]Peer, 0),
		history:      make(map[Packet]bool),
		transactions: make(map[string]account.SignedTransaction),
		socket:       polysocket,
		incoming:     incomingChannel,
		ledger:       ledger,
		sequencer:    false,
		latestBlock:  0,
		lock:         sync.Mutex{},
	}
	self := Peer{
		Addr: polysocket.GetAddr(),
	}
	node.peers = append(node.peers, self)
	go node.listen()
	return node, nil
}

func (node *Node) Promote(public *crypto.PublicKey, private *crypto.PrivateKey, blockSize int) {
	node.sequencer = true
	node.public = public
	node.private = private
	node.blockSize = blockSize
	node.transactionCounter = 0
	node.internal = make(chan account.SignedTransaction, blockSize)
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

func (node *Node) GetSequencerKey() *crypto.PublicKey {
	return node.public
}

func (node *Node) GetAddr() net.Addr {
	return node.socket.GetAddr()
}

// Listen for messages received on the internal channel.
// Ensure received message is a packet through type assertion,
// and ensure that it is not malformed. If succesful, pass the
// packet along to the packet handling method. Otherwise,
// skip message.
func (node *Node) listen() {
	for {
		msg := <-node.incoming
		switch packet := msg.(type) {
		case Packet:
			if isMalformed(packet) {
				log.Println("packet malformed, skipping")
				continue
			}
			isBlock := packet.Instruction == BlockAnn
			contains := false
			if !isBlock {
				_, found := node.history[packet]
				contains = found
				node.history[packet] = contains
			} else {
				contains = packet.Data.(SignedBlock).SequenceNumber <= node.latestBlock
			}
			if contains {
				log.Println("packet replay, skipping")
				continue
			}
			node.handle(packet)
		default:
			log.Println("unexpected message type, skipping")
			continue
		}
	}
}

// Checks is a packet is malformed. Returns true if the type of the data is
// not as expected given the instruction, or if the attatched vector clock is
// nil.
func isMalformed(packet Packet) bool {
	return false // TODO: Implement!
}

// Decodes and handles the passed network.Packet according to the instruction.
// Assumes passed packet is deliverable.
func (node *Node) handle(packet Packet) {
	node.lock.Lock()
	defer node.lock.Unlock()
	switch packet.Instruction {
	case PeerRequest:
		requester := packet.Data.(net.Addr)
		node.handlePeerRequest(requester)
	case PeerReply:
		peerData := packet.Data.(PeerData)
		node.handlePeerReply(peerData)
	case ConnAnn:
		peer := packet.Data.(Peer)
		node.handleConnectionAnnouncement(peer, packet)
	case Transaction:
		transaction := packet.Data.(account.SignedTransaction)
		node.handleTransaction(transaction, packet)
	case BlockAnn:
		block := packet.Data.(SignedBlock)
		node.handleBlock(block, packet)
	default:
		log.Printf("Unknown instruction %d \n", packet.Instruction)
	}
}

// Handles a peer request by sending a copy of this nodes peers
// to the address given as argument. Returns a error if no connection is
// established to the given address.
func (node *Node) handlePeerRequest(requester net.Addr) error {
	data := PeerData{
		Peers:        &node.peers,
		SequencerKey: node.public,
	}
	reply := Packet{
		Instruction: PeerReply,
		Data:        data,
	}
	return node.socket.Send(reply, requester)
}

// Handles a peer reply by merging the received list of peers with the
// the list of known peers maintained by the node. Strengthens the network
// by connecting to the last ten peers in the merged list.
func (node *Node) handlePeerReply(data PeerData) {
	node.public = data.SequencerKey
	for _, addr := range *data.Peers {
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
	node.peers = append(node.peers, peer) //TODO: Ensure that known peers are not re-added.
	node.socket.Broadcast(packet)         //TODO: Strengthen network if # connections below some threshold.
}

// Handle a transaction by propagating. If node is sequencer, increment transaction counter,
// and produce block once counter reaches blocksize. Wrap block in packet and send on "incoming"
// channel.
func (node *Node) handleTransaction(transaction account.SignedTransaction, packet Packet) {
	node.transactions[transaction.ID] = transaction
	node.socket.Broadcast(packet)
	if node.sequencer {
		node.internal <- transaction
		node.transactionCounter++
		if node.transactionCounter == node.blockSize {
			block, err := node.produceBlock()
			if err != nil {
				log.Println("error producing block: ", err)
			}
			blockPacket := Packet{
				Instruction: BlockAnn,
				Data:        *block,
			}
			node.handleBlock(*block, blockPacket)
			node.transactionCounter = 0
		}
	}
}

func (node *Node) handleBlock(block SignedBlock, packet Packet) {
	fmt.Println("Handling block")
	fmt.Println("begin validating")
	valid, _ := block.Validate(node.public)
	fmt.Println("finished validating")
	isNextBlock := block.SequenceNumber-node.latestBlock == 1
	if !valid {
		log.Println("unable to verify block: ", block.SequenceNumber)
		return
	}
	fmt.Println("block is valid")
	if !isNextBlock {
		log.Println("unexpected block sequence number: ", block.SequenceNumber, " expected: ", node.latestBlock+1)
		return
	}
	fmt.Println("block is next")
	for _, id := range block.Transactions {
		if transaction, contains := node.transactions[id]; contains {
			node.ledger.SignedTransaction(&transaction)
		} else {
			log.Panicln("transaction not found")
		}
	}
	node.latestBlock++
	fmt.Println("broadcasting blockpacket!")
	node.socket.Broadcast(packet)
}

func (node *Node) produceBlock() (*SignedBlock, error) {
	transactions := make([]string, 0)
	for i := 0; i < node.blockSize; i++ {
		transaction := <-node.internal
		transactions = append(transactions, transaction.ID)
	}
	block := Block{
		SequenceNumber: node.latestBlock + 1,
		Transactions:   transactions,
	}
	signedBlock, err := block.Sign(node.private)
	if err != nil {
		return nil, err
	}
	return signedBlock, nil
}

func (node *Node) SendTransaction(transaction account.SignedTransaction) {
	wrappedTransaction := Packet{
		Instruction: Transaction,
		Data:        transaction,
	}
	node.incoming <- wrappedTransaction
}

func contains(list []Peer, peer Peer) bool {
	for _, known := range list {
		if known == peer {
			return true
		}
	}
	return false
}
