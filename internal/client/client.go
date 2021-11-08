package client

import (
	"errors"
	"log"
	"net"
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
	client := Client{
		ledger:               ledger,
		node:                 node,
		internal:             internal,
		lock:                 sync.Mutex{},
		numberOfTransactions: 0,
		account:              "",
		public:               nil,
		private:              nil,
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
	if client.account == "" || client.private == nil || client.public == nil {
		return errors.New("invalid credentials")
	}
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

func (client *Client) ProvideCredentials(public *crypto.PublicKey, private *crypto.PrivateKey) error {
	client.lock.Lock()
	defer client.lock.Unlock()
	client.public = public
	client.private = private
	account, err := client.public.ToString()
	if err != nil {
		client.public = nil
		client.private = nil
		return err
	}
	client.account = account
	return nil
}

func (client *Client) GetBalance(account string) float64 {
	return client.ledger.GetBalance(account)
}

func (client *Client) GetAccount() string {
	return client.account
}

func (client *Client) GetPort() string {
	return strconv.Itoa(client.node.GetAddr().(*net.TCPAddr).Port)
}

func (client *Client) Connect(addr net.Addr) error {
	return client.node.Connect(addr)
}

func (client *Client) Close() []error {
	return client.node.Close()
}
