package account_test

import (
	"testing"
	"vicoin/crypto"
	"vicoin/internal/account"
	"vicoin/internal/registration"
)

func TestLedgersCanSetAndGetAccountBalance(t *testing.T) {
	ledger := account.NewLedger()
	ledger.SetBalance("hej", 10)
	if ledger.GetBalance("hej") != 10 {
		t.Error("Error when setting or getting balance")
	}
}
func TestLedgersCanPerformTransactionSignedBySendingAccount(t *testing.T) {
	registration.RegisterStructsWithGob()
	ledger := account.NewLedger()
	public, private, _ := crypto.KeyGen(2048)
	senderAccount, _ := public.ToString()
	ledger.SetBalance(senderAccount, 42)
	transaction, _ := account.NewSignedTransaction("id", senderAccount, "recipient", 10, private)
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
