package client

import (
	"net"
	"vicoin/internal/account"
	"vicoin/internal/node"
)

type Node interface {
	Connect(addr net.Addr) error
	Close() []error
	GetPeers() []node.Peer
	SendTransaction(transaction account.SignedTransaction)
}

type Ledger interface {
	SignedTransaction(transaction *account.SignedTransaction) error
	GetBalance(account string) float64
}
