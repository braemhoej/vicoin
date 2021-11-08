package client

import (
	"log"
	"strconv"
	"sync"
	"vicoin/crypto"
	"vicoin/internal/account"
	"vicoin/internal/node"
)

type Client struct {
	ledger               account.LedgerInterface
	node                 node.NodeInterface
	internal             chan account.SignedTransaction
	lock                 sync.Mutex
	numberOfTransactions int
	account              string
	public               *crypto.PublicKey
	private              *crypto.PrivateKey
}

func NewClient(ledger account.LedgerInterface, node node.NodeInterface, internal chan account.SignedTransaction) (*Client, error) {
	public, private, err := crypto.KeyGen(2048)
	if err != nil {
		return nil, err
	}
	account, err := public.ToString()
	if err != nil {
		return nil, err
	}
	client := Client{
		ledger:               ledger,
		node:                 node,
		internal:             internal,
		lock:                 sync.Mutex{},
		numberOfTransactions: 0,
		account:              account,
		public:               public,
		private:              private,
	}
	go client.handle()
	return &client, nil
}

func (client *Client) handle() {
	for {
		transaction := <-client.internal
		err := client.ledger.SignedTransaction(&transaction)
		log.Println(err)
		//TODO: Handle errors?
	}
}

func (client *Client) Transfer(amount float64, to string) error {
	client.lock.Lock()
	defer client.lock.Unlock()
	client.numberOfTransactions += 1
	transaction, err := account.NewSignedTransaction(strconv.Itoa(client.numberOfTransactions), client.account, to, amount, client.private)
	if err != nil {
		client.numberOfTransactions -= 1
		return err
	}
	err = client.ledger.SignedTransaction(transaction)
	if err != nil {
		client.numberOfTransactions -= 1
		return err
	}
	client.node.SendTransaction(*transaction)
	return nil
}

func (Client *Client) GetBalance(account string) float64 {
	return Client.ledger.GetBalance(account)
}
