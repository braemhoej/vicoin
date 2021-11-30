package node_test

import (
	"net"
	"reflect"
	"strconv"
	"testing"
	"time"
	"vicoin/account"
	"vicoin/crypto"
	"vicoin/node"
	"vicoin/registration"
	accountMocks "vicoin/test/mocks/account"
	mocks "vicoin/test/mocks/network"
)

func makeDependecies() (*mocks.MockPolysocket, chan interface{}, account.LedgerInterface) {
	internal := make(chan interface{})
	return &mocks.MockPolysocket{
		SentMessages:        make([]interface{}, 0),
		BroadcastedMessages: make([]interface{}, 0),
		Channel:             internal,
		Connections:         make([]net.Addr, 0),
	}, internal, accountMocks.NewMockLedger()
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
	mock, internal, ledger := makeDependecies()
	n, _ := node.NewNode(mock, internal, ledger)
	n.Connect(&node.Peer{}) // Mock address
	if len(mock.SentMessages) != 1 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.SentMessages[0].(node.Packet)
	if msg.Instruction != node.PeerRequest {
		t.Errorf("Unexpected instruction %d, want 0", msg.Instruction)
	}
}

func TestNodesSendPeerReplyUponRequest(t *testing.T) {
	mock, internal, ledger := makeDependecies()
	node.NewNode(mock, internal, ledger)
	mock.InjectMessage(node.Packet{
		Instruction: node.PeerRequest,
		Data:        &net.TCPAddr{},
	})
	time.Sleep(50 * time.Millisecond)
	if len(mock.SentMessages) != 1 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.SentMessages[0].(node.Packet)
	if msg.Instruction != node.PeerReply {
		t.Errorf("Unexpected instruction %d, want 1", msg.Instruction)
	}
}

func TestNewNodesMergePeerRequestIntoKnownPeers(t *testing.T) {
	mock, internal, ledger := makeDependecies()
	n, _ := node.NewNode(mock, internal, ledger)
	for i := 1; i < 10; i++ {
		mock.InjectMessage(node.Packet{
			Instruction: node.ConnAnn,
			Data: node.Peer{
				Addr: &net.IPAddr{
					IP: []byte("mock" + strconv.Itoa(i)),
				},
			},
		})
	}
	time.Sleep(60 * time.Millisecond)
	peers := n.GetPeers()

	if len(peers) != 10 {
		t.Errorf("Unexpected number of peers %d, want 10", len(peers))
	}
}

func TestNodesBroadcastConnectionAnnouncementUponConnection(t *testing.T) {
	mock, internal, ledger := makeDependecies()
	n, _ := node.NewNode(mock, internal, ledger)
	n.Connect(&node.Peer{}) // Mock address
	time.Sleep(50 * time.Millisecond)
	if len(mock.BroadcastedMessages) != 1 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.BroadcastedMessages[0].(node.Packet)
	if msg.Instruction != node.ConnAnn {
		t.Errorf("Unexpected instruction %d, want 3", msg.Instruction)
	}
}
func TestNewNodesAddAnnouncedConnectionsToPeerList(t *testing.T) {
	mock, internal, ledger := makeDependecies()
	n, _ := node.NewNode(mock, internal, ledger)
	for i := 1; i < 10; i++ {
		mock.InjectMessage(node.Packet{
			Instruction: node.ConnAnn,
			Data: node.Peer{
				Addr: &net.IPAddr{
					IP: []byte("mock" + strconv.Itoa(i)),
				},
			},
		})
	}
	time.Sleep(50 * time.Millisecond)
	peers := n.GetPeers()
	if len(peers) != 10 {
		t.Errorf("Unexpected number of peers %d, want 10", len(peers))
	}
}
func TestNodesPropagateConnectionAnnouncements(t *testing.T) {
	mock, internal, ledger := makeDependecies()
	n, _ := node.NewNode(mock, internal, ledger)
	n.Connect(&node.Peer{}) // Mock address
	mock.InjectMessage(node.Packet{Instruction: node.ConnAnn, Data: node.Peer{}})
	time.Sleep(50 * time.Millisecond)
	if len(mock.BroadcastedMessages) != 2 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.BroadcastedMessages[1].(node.Packet)
	if msg.Instruction != node.ConnAnn {
		t.Errorf("Unexpected instruction %d, want 3", msg.Instruction)
	}
}

func TestNodesPropagateTransactions(t *testing.T) {
	mock, internal, ledger := makeDependecies()
	n, _ := node.NewNode(mock, internal, ledger)
	n.Connect(&node.Peer{}) // Mock address
	mock.InjectMessage(node.Packet{Instruction: node.Transaction, Data: account.SignedTransaction{}})
	time.Sleep(50 * time.Millisecond)
	if len(mock.BroadcastedMessages) != 2 {
		t.Errorf("Unexpected number of sent messages %d, want 1", len(mock.SentMessages))
	}
	msg := mock.BroadcastedMessages[1].(node.Packet)
	if msg.Instruction != node.Transaction {
		t.Errorf("Unexpected instruction %d, want 4", msg.Instruction)
	}
}

func TestNodesPropagateBlocks(t *testing.T) {
	registration.RegisterStructsWithGob()
	mock, internal, ledger := makeDependecies()
	n, _ := node.NewNode(mock, internal, ledger)
	n.Connect(&node.Peer{}) // Mock address
	public, private, _ := crypto.KeyGen(2048)
	block := node.Block{
		SequenceNumber: 1,
		Transactions:   make([]string, 0),
	}
	signedBlock, _ := block.Sign(private)
	peers := make([]node.Peer, 0)
	mock.InjectMessage(node.Packet{Instruction: node.PeerReply, Data: node.PeerData{
		Peers:        &peers,
		SequencerKey: public,
	}})
	mock.InjectMessage(node.Packet{Instruction: node.BlockAnn, Data: *signedBlock})
	time.Sleep(50 * time.Millisecond)
	if len(mock.BroadcastedMessages) != 2 {
		t.Errorf("Unexpected number of sent messages %d, want 2", len(mock.SentMessages))
	}
	msg := mock.BroadcastedMessages[1].(node.Packet)
	if msg.Instruction != node.BlockAnn {
		t.Errorf("Unexpected instruction %d, want 5", msg.Instruction)
	}
}

func TestPromotedNodesSendBlocksOnceEnoughTransactionsAreSeen(t *testing.T) {
	registration.RegisterStructsWithGob()
	mock, internal, ledger := makeDependecies()
	n, _ := node.NewNode(mock, internal, ledger)
	n.Connect(&node.Peer{}) // Mock address
	public, private, _ := crypto.KeyGen(2048)
	blockSize := 1
	n.Promote(public, private, blockSize)
	sender, _ := public.ToString()
	transaction, _ := account.NewSignedTransaction("randomID", sender, "randomAccount", 0, private)
	mock.InjectMessage(node.Packet{Instruction: node.Transaction, Data: *transaction})
	time.Sleep(50 * time.Millisecond)
	msg := mock.BroadcastedMessages[2].(node.Packet)
	if len(mock.BroadcastedMessages) != 3 {
		t.Errorf("Unexpected number of sent messages %d, want 2", len(mock.SentMessages))
	}
	if msg.Instruction != node.BlockAnn {
		t.Errorf("Unexpected instruction %d, want 5", msg.Instruction)
	}
}

func TestNodesUpdateSequencePublicKeyFromPeerReply(t *testing.T) {
	registration.RegisterStructsWithGob()
	mock, internal, ledger := makeDependecies()
	n, _ := node.NewNode(mock, internal, ledger)
	public, _, _ := crypto.KeyGen(2048)
	peers := make([]node.Peer, 0)
	data := node.Packet{
		Instruction: node.PeerReply,
		Data: node.PeerData{
			Peers:        &peers,
			SequencerKey: public,
		},
	}
	mock.InjectMessage(data)
	time.Sleep(100 * time.Millisecond)
	key := n.GetSequencerKey()
	if key == nil {
		t.Error("sequencer key unset")
	}
	if !reflect.DeepEqual(key, public) {
		t.Error("sequencer key doesn't match sent key")
	}
}

func TestNodesPerformTransactionsInValidBlocks(t *testing.T) {
	registration.RegisterStructsWithGob()
	mock, internal, ledger := makeDependecies()
	n, _ := node.NewNode(mock, internal, ledger)
	n.Connect(&node.Peer{}) // Mock address
	public, private, _ := crypto.KeyGen(2048)
	sender, _ := public.ToString()
	peers := make([]node.Peer, 0)
	peerReply := node.Packet{
		Instruction: node.PeerReply,
		Data: node.PeerData{
			Peers:        &peers,
			SequencerKey: public,
		},
	}
	mock.InjectMessage(peerReply)
	transaction, _ := account.NewSignedTransaction("randomID", sender, "randomAccount", 0, private)
	mock.InjectMessage(node.Packet{Instruction: node.Transaction, Data: *transaction})
	block := node.Block{
		SequenceNumber: 1,
		Transactions:   []string{"randomID"},
	}
	signedBlock, _ := block.Sign(private)
	mock.InjectMessage(node.Packet{Instruction: node.BlockAnn, Data: *signedBlock})
	time.Sleep(100 * time.Millisecond)
	mockLedger := ledger.(*accountMocks.MockLedger)
	transactions := mockLedger.Transactions
	if len(transactions) != 1 {
		t.Errorf("Unexpected number of transactions %d, want 1", len(transactions))
	}
	if transactions[0].ID != "randomID" {
		t.Error("Unexpected transaction id: ", transactions[0].ID)
	}
}

func TestNodesDoesNotPerformLooseTransactions(t *testing.T) {
	registration.RegisterStructsWithGob()
	mock, internal, ledger := makeDependecies()
	n, _ := node.NewNode(mock, internal, ledger)
	n.Connect(&node.Peer{}) // Mock address
	public, private, _ := crypto.KeyGen(2048)
	sender, _ := public.ToString()
	transaction, _ := account.NewSignedTransaction("randomID", sender, "randomAccount", 0, private)
	mock.InjectMessage(node.Packet{Instruction: node.Transaction, Data: *transaction})
	time.Sleep(100 * time.Millisecond)
	mockLedger := ledger.(*accountMocks.MockLedger)
	transactions := mockLedger.Transactions
	if len(transactions) != 0 {
		t.Errorf("Unexpected number of transactions %d, want 0", len(transactions))
	}
}
