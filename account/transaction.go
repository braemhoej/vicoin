package account

import (
	"encoding/base64"
	"vicoin/crypto"
)

type Transaction struct {
	ID     string
	From   string
	To     string
	Amount float64
}

type SignedTransaction struct {
	ID        string
	From      string
	To        string
	Amount    float64
	Signature string
}

func NewSignedTransaction(id string, from string, to string, amount float64, key *crypto.PrivateKey) (*SignedTransaction, error) {
	unsignedTransaction := Transaction{
		ID:     id,
		From:   from,
		To:     to,
		Amount: amount,
	}
	signature, err := crypto.Sign(unsignedTransaction, key)
	if err != nil {
		return nil, err
	}
	return &SignedTransaction{
		ID:        id,
		From:      from,
		To:        to,
		Amount:    amount,
		Signature: base64.StdEncoding.EncodeToString(signature),
	}, nil
}

func (signedTransaction *SignedTransaction) Validate(key *crypto.PublicKey) (isValid bool, err error) {
	unsignedTransaction := Transaction{
		ID:     signedTransaction.ID,
		From:   signedTransaction.From,
		To:     signedTransaction.To,
		Amount: signedTransaction.Amount,
	}
	bytes, err := base64.StdEncoding.DecodeString(signedTransaction.Signature)
	if err != nil {
		return false, err
	}
	isValid, err = crypto.Validate(unsignedTransaction, bytes, key)
	if err != nil {
		return false, err
	}
	return isValid, err
}
