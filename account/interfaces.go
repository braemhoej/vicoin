package account

type LedgerInterface interface {
	SignedTransaction(transaction *SignedTransaction) error
	GetBalance(account string) float64
}
