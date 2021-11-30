package node

import (
	"net"
	"vicoin/crypto"
)

type Peer struct {
	Addr net.Addr
}

type PeerData struct {
	Peers        []Peer
	SequencerKey crypto.PublicKey
}
