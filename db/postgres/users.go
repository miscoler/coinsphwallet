package postgres

import (
	"coinsWallet/account"
	"context"
	"fmt"

	"github.com/go-pg/pg"
)

type AccountsRepository struct {
	db *pg.DB
}

func NewAccountsRepository(dbSettings *PostgresConfig) (*AccountsRepository, error) {
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
	return &AccountsRepository{db: db}, nil
}

func (s *AccountsRepository) GetAllUsers(ctx context.Context) ([]*account.Account, error) {
	var records []*account.Account
	_, err := s.db.QueryContext(ctx,
		&records, "SELECT id, balance, currency FROM accounts",
	)
	return records, err
}

func (s *AccountsRepository) GetUser(ctx context.Context, accountID string) (*account.Account, error) {
	record := &account.Account{}
	_, err := s.db.QueryOneContext(ctx,
		record, "SELECT id, balance, currency FROM accounts WHERE id=?", accountID,
	)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return record, nil
}
