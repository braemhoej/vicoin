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
	accounts := readAccountsFromFile()
	for _, account := range accounts {
		ledger.accounts[account.public] = account.balance
	}
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
