package account

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"vicoin/crypto"
)

type account struct {
	public  string
	balance float64
}

type Ledger struct {
	accounts map[string]float64
	lock     sync.Mutex
}

func NewLedger() *Ledger {
	ledger := new(Ledger)
	ledger.accounts = make(map[string]float64)
	public := "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"
	ledger.accounts[public] = 1000
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

func readAccountsFromFile() []account {
	data, err := os.ReadFile("/workspaces/vicoin/account/accounts.txt")
	if err != nil {
		log.Panicln(err)
	}
	var accounts []account
	err = json.Unmarshal(data, &accounts)
	if err != nil {
		log.Panicln(err)
	}
	return accounts
}
