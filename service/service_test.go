package service

import (
	"coinsWallet/db/postgres"
	"coinsWallet/transaction"
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	dStub "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/stretchr/testify/assert"
)

func setupTests() (*walletService, *migrate.Migrate) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	instance, err := dStub.WithInstance(db, &dStub.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fsrc, err := (&file.File{}).Open("file://../db/postgres")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance(
		"file",
		fsrc,
		"postgres",
		instance,
	)

	if err != nil {
		log.Fatal(err)
	}
	m.Up()
	conf, err := postgres.Parse("../configTest.json")
	if err != nil {
		panic(err)
	}
	cfg := &postgres.PostgresConfig{
		Host:     conf.Postgres.Host,
		Port:     conf.Postgres.Port,
		Database: conf.Postgres.Database,
		User:     conf.Postgres.User,
		Password: conf.Postgres.Password,
	}

	accounts, _ := postgres.NewAccountsRepository(cfg)
	transactions, _ := postgres.NewTransactionsRepository(cfg)
	return &walletService{
		accounts:     accounts,
		transactions: transactions,
	}, m
}

func teardownTests(m *migrate.Migrate) {
	err := m.Down()
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateTransactions(t *testing.T) {
	s, m := setupTests()
	defer teardownTests(m)
	ctx := context.Background()

	err := s.CreateTransaction(ctx, transaction.Transaction{
		Sender:   "mr house",
		Receiver: "courier",
		Currency: "CAPS",
		Amount:   1000,
	})
	assert.Equal(t, nil, err)
	sender, err := s.accounts.GetUser(ctx, "mr house")
	assert.Equal(t, nil, err)
	assert.Equal(t, 9000.0, sender.Balance)
	receiver, err := s.accounts.GetUser(ctx, "courier")
	assert.Equal(t, nil, err)
	assert.Equal(t, receiver.Balance, 1100.0)
	transactions, err := s.GetTransactions(ctx)
	assert.Equal(t, nil, err)
	expected := []*transaction.Transaction{
		{
			Sender:   "mr house",
			Receiver: "courier",
			Currency: "CAPS",
			Amount:   1000,
		},
	}
	assert.Equal(t, expected, transactions)
	err = s.CreateTransaction(ctx, transaction.Transaction{
		Sender:   "courier",
		Receiver: "mr house",
		Currency: "CAPS",
		Amount:   1000,
	})

	sender, err = s.accounts.GetUser(ctx, "mr house")
	assert.Equal(t, nil, err)
	assert.Equal(t, 10000.0, sender.Balance)
	receiver, err = s.accounts.GetUser(ctx, "courier")
	assert.Equal(t, nil, err)
	assert.Equal(t, receiver.Balance, 100.0)
	transactions, err = s.GetTransactions(ctx)
	assert.Equal(t, nil, err)
	expected = append(expected, &transaction.Transaction{
		Sender:   "courier",
		Receiver: "mr house",
		Currency: "CAPS",
		Amount:   1000,
	})
	assert.Equal(t, expected, transactions)
	assert.Equal(t, nil, err)

}

func TestAccountDoesNotExist(t *testing.T) {
	s, m := setupTests()
	defer teardownTests(m)
	ctx := context.Background()

	err := s.CreateTransaction(ctx, transaction.Transaction{
		Sender:   "asds",
		Receiver: "yes man",
		Currency: "CAPS",
		Amount:   1000,
	})
	assert.Equal(t, errors.As(err, &AccountNotFound), true)
	receiver, _ := s.accounts.GetUser(ctx, "yes man")
	assert.Equal(t, receiver.Balance, 10.0)
	transactions, err := s.GetTransactions(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(transactions))

	err = s.CreateTransaction(ctx, transaction.Transaction{
		Sender:   "yes man",
		Receiver: "asds",
		Currency: "CAPS",
		Amount:   1000,
	})
	assert.Equal(t, errors.As(err, &AccountNotFound), true)
	sender, _ := s.accounts.GetUser(ctx, "yes man")
	assert.Equal(t, sender.Balance, 10.0)
	transactions, err = s.GetTransactions(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(transactions))

}

func TestDifferentCurrencies(t *testing.T) {
	s, m := setupTests()
	defer teardownTests(m)

	ctx := context.Background()

	err := s.CreateTransaction(ctx, transaction.Transaction{
		Sender:   "mr house",
		Receiver: "yes man",
		Currency: "USD",
		Amount:   1000,
	})
	assert.Equal(t, errors.As(err, &CurrenciesDoNotMatch), true)
	sender, _ := s.accounts.GetUser(ctx, "mr house")
	receiver, _ := s.accounts.GetUser(ctx, "yes man")
	assert.Equal(t, sender.Balance, 10000.0)
	assert.Equal(t, receiver.Balance, 10.0)
	transactions, err := s.GetTransactions(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(transactions))
}

func TestLowBalance(t *testing.T) {
	s, m := setupTests()
	defer teardownTests(m)
	ctx := context.Background()

	err := s.CreateTransaction(ctx, transaction.Transaction{
		Sender:   "yes man",
		Receiver: "mr house",
		Currency: "CAPS",
		Amount:   10000,
	})
	assert.Equal(t, errors.As(err, &transaction.InsufficientFunds), true)
	sender, _ := s.accounts.GetUser(ctx, "yes man")
	receiver, _ := s.accounts.GetUser(ctx, "mr house")
	assert.Equal(t, sender.Balance, 10.0)
	assert.Equal(t, receiver.Balance, 10000.0)
	transactions, err := s.GetTransactions(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(transactions))
}
