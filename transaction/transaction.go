package transaction

import (
	"context"
	"fmt"
)

var InsufficientFunds = fmt.Errorf("insufficient funds")

type Transaction struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction Transaction) error
	GetAllTransactions(ctx context.Context) ([]*Transaction, error)
}
