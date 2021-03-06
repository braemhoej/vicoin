package node

import (
	"net"
	"vicoin/internal/account"
)

type NodeInterface interface {
	Connect(addr net.Addr) error
	Close() []error
	SendTransaction(transaction account.SignedTransaction)
	GetAddr() net.Addr
}
