package account

import (
	"vicoin/crypto"
)

type Transaction struct {
	ID     string  // Any string
	From   string  // A verification key coded as a string
	To     string  // A verification key coded as a string
	Amount float64 // Amount to transfer
}

func NewTransaction(id string, from string, to string, amount float64) *Transaction {
	return &Transaction{
		ID:     id,
		From:   from,
		To:     to,
		Amount: amount,
	}
}

type SignedTransaction struct {
	ID        string  // Any string
	From      string  // A verification key coded as a string
	To        string  // A verification key coded as a string
	Amount    float64 // Amount to transfer
	Signature string  // Potential signature coded as string
}

func NewSignedTransaction(id string, from string, to string, amount float64, key *crypto.PrivateKey) (*SignedTransaction, error) {
	transaction := NewTransaction(id, from, to, amount)
	signature, err := crypto.Sign(transaction, key)
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
func (transaction *SignedTransaction) Validate(key *crypto.PublicKey) (isValid bool, err error) {
	unsignedTransaction := NewTransaction(transaction.ID, transaction.From, transaction.To, transaction.Amount)
	isValid, err = crypto.Validate(unsignedTransaction, []byte(transaction.Signature), key)
	if err != nil {
		return false, err
	}
	return isValid, err
}
