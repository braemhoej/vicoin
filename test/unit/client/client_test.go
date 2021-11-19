package client_test

import (
	"testing"
	"vicoin/account"
	"vicoin/client"
	"vicoin/crypto"
	"vicoin/registration"
	mocksAcc "vicoin/test/mocks/account"
	mocksNode "vicoin/test/mocks/node"
)

func makeDependencies() (*mocksAcc.MockLedger, *mocksNode.MockNode) {
	return mocksAcc.NewMockLedger(), mocksNode.NewMockNode()
}

func TestClientReturnsAPointerToANewClient(t *testing.T) {
	registration.RegisterStructsWithGob()
	ledger, node := makeDependencies()
	internal := make(chan account.SignedTransaction)
	client, err := client.NewClient(ledger, node, internal)
	if err != nil {
		t.Error(err)
	}
	if client == nil {
		t.Error("Nil returned")
	}
}

func TestClientAttemptsToPerformTransactionsReceivedOnInternalChannel(t *testing.T) {
	registration.RegisterStructsWithGob()
	ledger, node := makeDependencies()
	internal := make(chan account.SignedTransaction)
	client.NewClient(ledger, node, internal)
	internal <- account.SignedTransaction{}
	if len(ledger.Transactions) != 1 {
		t.Errorf("Unexpected number of transactions %d, want 1", len(ledger.Transactions))
	}
}

func TestTransferAttemptsToPerformTransactionAndIncrementsNumberOfTransactions(t *testing.T) {
	registration.RegisterStructsWithGob()
	public, private, _ := crypto.KeyGen(2048)
	ledger, node := makeDependencies()
	internal := make(chan account.SignedTransaction)
	c, err := client.NewClient(ledger, node, internal)
	c.ProvideCredentials(public, private)
	c.Transfer(10, "Santa")
	if err != nil {
		t.Error(err)
	}
	if len(ledger.Transactions) != 1 {
		t.Errorf("Unexpected number of transactions %d, want 1", len(ledger.Transactions))
	}
	if ledger.Transactions[0].ID != "1" {
		t.Error("Failed to increment transaction ID")
	}
}
