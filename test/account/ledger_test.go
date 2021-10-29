package account

import (
	"testing"
	"vicoin/account"
	"vicoin/crypto"
	"vicoin/registration"
)

func TestLedgersCanPerformTransactionSignedBySendingAccount(t *testing.T) {
	registration.RegisterStructsWithGob()
	ledger := account.NewLedger()
	public, private, _ := crypto.KeyGen(2048)
	senderAccount, _ := public.ToString()
	transaction, _ := account.NewSignedTransaction("id", senderAccount, "recipient", 0, private)
	err := ledger.SignedTransaction(transaction)
	if err != nil {
		t.Error("Error when performing legitimate transaction : ", err)
	}
}

func TestLedgersCantPerformTransactionSignedByForeignAccount(t *testing.T) {
	registration.RegisterStructsWithGob()
	ledger := account.NewLedger()
	sender, _, _ := crypto.KeyGen(2048)
	_, foreign, _ := crypto.KeyGen(2048)
	senderAccount, _ := sender.ToString()
	transaction, _ := account.NewSignedTransaction("id", senderAccount, "recipient", 0, foreign)
	err := ledger.SignedTransaction(transaction)
	if err == nil {
		t.Error("Error: allowed illegitemate transaction ")
	}
}

func TestLedgersCantPerformTransactionsThatViolateBalance(t *testing.T) {
	registration.RegisterStructsWithGob()
	ledger := account.NewLedger()
	public, private, _ := crypto.KeyGen(2048)
	senderAccount, _ := public.ToString()
	transaction, _ := account.NewSignedTransaction("id", senderAccount, "recipient", 1, private)
	err := ledger.SignedTransaction(transaction)
	if err == nil {
		t.Error("Error: allowed illegitemate transaction ")
	}
}
