package node

import (
	"fmt"
	"io"
	"net"
	"net/rpc"
	"vicoin/account"
)

const (
	ping               = ".Ping"
	peerRequest        = ".GetPeers"
	sendTransaction    = ".AddTransaction"
	announceConnection = ".AnnounceConnection"
	addConnection      = ".AddConnection"
)

var noArg = new(struct{})

type Peer struct {
	Addr   *net.TCPAddr
	UUID   string
	client *rpc.Client
}

func NewPeer(addr *net.TCPAddr, name string) *Peer {
	return &Peer{
		Addr: addr,
		UUID: name,
	}
}

func (peer *Peer) AnnounceConnection(argPeer *Peer) error {
	fmt.Println("AnnounceConnection called on peer")
	err := peer.client.Call(peer.UUID+announceConnection, argPeer, &noArg)
	if err != nil {
		return err
	}
	return nil
}

func (peer *Peer) AddConnection(argPeer *Peer) error {
	err := peer.client.Call(peer.UUID+addConnection, argPeer, &noArg)
	if err != nil {
		return err
	}
	return nil
}

func (peer *Peer) SendTransaction(transaction *account.SignedTransaction) error {
	err := peer.client.Call(peer.UUID+sendTransaction, transaction, &noArg)
	if err != nil {
		return err
	}
	return nil
}

func (peer *Peer) Connect() (err error) {
	client, err := rpc.Dial(peer.Addr.Network(), peer.Addr.String())
	if err != nil {
		return err
	}
	peer.client = client
	return err
}

func (peer *Peer) IsConnected() (bool, error) {
	var pong = false
	err := peer.client.Call(peer.UUID+ping, noArg, &pong)
	if err == io.EOF {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return pong, nil
}

func (peer *Peer) RequestPeers() ([]*Peer, error) {
	var peers []*Peer
	err := peer.client.Call(peer.UUID+peerRequest, noArg, &peers)
	if err != nil {
		return nil, err
	}
	return peers, nil
}
