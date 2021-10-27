package account

import (
	"sync"
	"vicoin/crypto"
)

type Ledger struct {
	Accounts map[string]float64
	lock     sync.Mutex
}

func NewLedger() *Ledger {
	ledger := new(Ledger)
	ledger.Accounts = make(map[string]float64)
	return ledger
}
func (ledger *Ledger) SignedTransaction(transaction *SignedTransaction) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()
	sendersPublicKey, err := new(crypto.PublicKey).FromString(transaction.From)
	if err != nil {
		return err
	}
	validSignature, err := transaction.Validate(sendersPublicKey)
	if err != nil {
		return err
	}
	if validSignature {
		ledger.Accounts[transaction.From] -= transaction.Amount
		ledger.Accounts[transaction.To] += transaction.Amount
	}
	return nil
}
