package node

import (
	"testing"
	"vicoin/crypto"
	"vicoin/node"
	"vicoin/registration"
)

func TestBlocksCanBeSigned(t *testing.T) {
	registration.RegisterStructsWithGob()
	_, private, _ := crypto.KeyGen(2048)
	block := node.Block{
		SequenceNumber: 1,
		Transactions:   make([]string, 0),
	}
	signedBlock, err := block.Sign(private)
	if err != nil {
		t.Error("Error when signing block: ", err)
	}
	if signedBlock.Signature == "" {
		t.Error("No signature was generated")
	}
}
func TestCorrectlySignedTransactionsCanBeValidated(t *testing.T) {
	registration.RegisterStructsWithGob()
	public, private, _ := crypto.KeyGen(2048)
	block := node.Block{
		SequenceNumber: 1,
		Transactions:   make([]string, 0),
	}
	signedBlock, _ := block.Sign(private)
	isValid, err := signedBlock.Validate(public)
	if !isValid || err != nil {
		t.Error("Unable to validate correctly signed transaction, error : ", err)
	}
}
func TestCorrectlySignedTransactionsCantBeValidatedWithForeignKey(t *testing.T) {
	registration.RegisterStructsWithGob()
	_, private, _ := crypto.KeyGen(2048)
	foreignPublic, _, _ := crypto.KeyGen(2048)
	block := node.Block{
		SequenceNumber: 1,
		Transactions:   make([]string, 0),
	}
	signedBlock, _ := block.Sign(private)
	isValid, _ := signedBlock.Validate(foreignPublic)
	if isValid {
		t.Error("Unable to validate correctly signed transaction")
	}
}
func TestIncorrectlySignedTransactionsCantBeValidated(t *testing.T) {
	registration.RegisterStructsWithGob()
	public, private, _ := crypto.KeyGen(2048)
	block := node.Block{
		SequenceNumber: 1,
		Transactions:   make([]string, 0),
	}
	signedBlock, _ := block.Sign(private)
	signedBlock.SequenceNumber = 42
	isValid, _ := signedBlock.Validate(public)
	if isValid {
		t.Error("Unable to validate correctly signed transaction")
	}
}
