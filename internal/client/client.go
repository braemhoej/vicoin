package client

import (
	"sync"
	"vicoin/internal/account"
	"vicoin/internal/node"
)

type Client struct {
	ledger   Ledger
	node     Node
	internal chan account.SignedTransaction
	lock     sync.Mutex
}

func NewClient(ledger *account.Ledger, node *node.Node, internal chan account.SignedTransaction) *Client {
	client := Client{
		ledger:   ledger,
		node:     node,
		internal: internal,
		lock:     sync.Mutex{},
	}
	return &client
}
