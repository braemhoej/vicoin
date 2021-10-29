package account

import (
	"errors"
	"sync"
	"vicoin/crypto"
)

type Ledger struct {
	accounts map[string]float64
	lock     sync.Mutex
}

func NewLedger() *Ledger {
	ledger := new(Ledger)
	ledger.accounts = make(map[string]float64)
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
	if err != nil || !validSignature {
		return errors.New("unable to validate transaction")
	}
	if validSignature {
		err := ledger.transfer(transaction.From, transaction.To, transaction.Amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ledger *Ledger) GetBalance(account string) float64 {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()
	return ledger.accounts[account]
}

func (ledger *Ledger) transfer(from string, to string, amount float64) error {
	if ledger.accounts[from] < amount {
		return errors.New("insufficient funds")
	}
	ledger.accounts[from] -= amount
	ledger.accounts[to] += amount
	return nil
}

func (ledger *Ledger) SetBalance(account string, amount float64) {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()
	ledger.accounts[account] = amount
}
