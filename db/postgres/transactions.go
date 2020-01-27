package postgres

import (
	"coinsWallet/transaction"
	"context"
	"fmt"

	"github.com/go-pg/pg"
)

type TransactionsRepository struct {
	db *pg.DB
}

func NewTransactionsRepository(dbSettings *PostgresConfig) (*TransactionsRepository, error) {
	db := pg.Connect(&pg.Options{
		User:     dbSettings.User,
		Password: dbSettings.Password,
		Database: dbSettings.Database,
		Addr:     fmt.Sprintf("%s:%d", dbSettings.Host, dbSettings.Port),
	})

	_, err := db.Exec("SELECT 1")
	if err != nil {
		return nil, fmt.Errorf("no connection to database: %w", err)
	}

	return &TransactionsRepository{db: db}, nil
}

func (s *TransactionsRepository) GetAllTransactions(ctx context.Context) ([]*transaction.Transaction, error) {
	var records []*transaction.Transaction
	_, err := s.db.QueryContext(ctx,
		&records,
		"SELECT sender, receiver, currency, amount FROM transactions ",
	)
	return records, err
}

func (s *TransactionsRepository) CreateTransaction(ctx context.Context, payment transaction.Transaction) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		var senderBalance float64

		_, err := tx.QueryOneContext(ctx,
			pg.Scan(&senderBalance),
			"SELECT balance FROM accounts where id=? FOR UPDATE",
			payment.Sender,
		)
		if err != nil {
			return err
		}

		if senderBalance-payment.Amount < 0 {
			return transaction.InsufficientFunds
		}

		_, err = tx.ExecOneContext(ctx,
			"INSERT INTO transactions (sender, receiver, amount, currency) VALUES (?, ?, ?, ?)",
			payment.Sender, payment.Receiver, payment.Amount, payment.Currency,
		)

		if err != nil {
			return err
		}

		_, err = tx.ExecOneContext(ctx,
			"UPDATE accounts SET balance = balance - ? WHERE id = ?",
			payment.Amount, payment.Sender,
		)
		if err != nil {
			return err
		}

		_, err = tx.ExecOneContext(ctx,
			"UPDATE accounts SET balance = balance + ? WHERE id = ?",
			payment.Amount, payment.Receiver,
		)
		return err
	})
	return err
}
