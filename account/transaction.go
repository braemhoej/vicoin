package account

import (
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
		Signature: string(signature),
	}, nil
}
func (signedTransaction *SignedTransaction) Validate(key *crypto.PublicKey) (isValid bool, err error) {
	unsignedTransaction := Transaction{
		ID:     signedTransaction.ID,
		From:   signedTransaction.From,
		To:     signedTransaction.To,
		Amount: signedTransaction.Amount,
	}
	isValid, err = crypto.Validate(unsignedTransaction, []byte(signedTransaction.Signature), key)
	if err != nil {
		return false, err
	}
	return isValid, err
}
