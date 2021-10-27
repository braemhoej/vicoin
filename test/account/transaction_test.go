package account

import (
	"encoding/gob"
	"testing"
	"vicoin/account"
	"vicoin/crypto"
)

func TestTransactionsCanBeSigned(t *testing.T) {
	gob.Register(account.Transaction{})
	_, private, _ := crypto.KeyGen(2048)
	signedTransaction, _ := account.NewSignedTransaction("id", "claus", "santa", 24.12, private)
	if signedTransaction.Signature == "" {
		t.Error("No signature was generated")
	}
}
func TestCorrectlySignedTransactionsCanBeValidated(t *testing.T) {
	gob.Register(account.Transaction{})
	public, private, _ := crypto.KeyGen(2048)
	signedTransaction, _ := account.NewSignedTransaction("id", "claus", "santa", 24.12, private)
	isValid, err := signedTransaction.Validate(public)
	if !isValid || err != nil {
		t.Error("Unable to validate correctly signed transaction, error : ", err)
	}
}
func TestCorrectlySignedTransactionsCantBeValidatedWithForeignKey(t *testing.T) {
	gob.Register(account.Transaction{})
	_, private, _ := crypto.KeyGen(2048)
	foreignPublic, _, _ := crypto.KeyGen(2048)
	signedTransaction, _ := account.NewSignedTransaction("id", "claus", "santa", 24.12, private)
	isValid, _ := signedTransaction.Validate(foreignPublic)
	if isValid {
		t.Error("Unable to validate correctly signed transaction")
	}
}
func TestIncorrectlySignedTransactionsCantBeValidated(t *testing.T) {
	gob.Register(account.Transaction{})
	public, private, _ := crypto.KeyGen(2048)
	signedTransaction, _ := account.NewSignedTransaction("id", "claus", "santa", 24.12, private)
	signedTransaction.From = "darthvader"
	isValid, _ := signedTransaction.Validate(public)
	if isValid {
		t.Error("Unable to validate correctly signed transaction")
	}
}
