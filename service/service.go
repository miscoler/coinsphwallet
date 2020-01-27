package service

import (
	"coinsWallet/transaction"
	"context"
	"fmt"

	"coinsWallet/account"
)

var (
	SenderNotFound         = fmt.Errorf("sender user not found")
	ReceiverNotFound       = fmt.Errorf("receiver user not found")
	TransactionNotPositive = fmt.Errorf("transaction amount should be positive")
	AccountNotFound        = fmt.Errorf("account not found")
	CurrenciesDoNotMatch   = fmt.Errorf("currencies do not match")
)

type WalletService interface {
	CreateTransaction(ctx context.Context, transaction transaction.Transaction) error
	GetTransactions(ctx context.Context) ([]*transaction.Transaction, error)
	GetAccounts(ctx context.Context) ([]*account.Account, error)
}

type walletService struct {
	accounts     account.AccountRepository
	transactions transaction.TransactionRepository
}

func NewService(transactions transaction.TransactionRepository, accounts account.AccountRepository) WalletService {
	return &walletService{transactions: transactions, accounts: accounts}
}

func (s walletService) GetAccounts(ctx context.Context) ([]*account.Account, error) {
	return s.accounts.GetAllUsers(ctx)
}

func (s walletService) GetTransactions(ctx context.Context) ([]*transaction.Transaction, error) {
	return s.transactions.GetAllTransactions(ctx)
}

func (s *walletService) CreateTransaction(ctx context.Context, transaction transaction.Transaction) error {

	//Correct float comparison omitted for the sake of simplicity
	if transaction.Amount <= 0 {
		return TransactionNotPositive
	}

	if transaction.Sender == transaction.Receiver {
		return fmt.Errorf("transaction sender can not be its receiver")
	}

	sender, err := s.accounts.GetUser(ctx, transaction.Sender)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if sender == nil {
		return fmt.Errorf("%w %s", AccountNotFound, transaction.Sender)
	}
	if sender.Currency != transaction.Currency {
		return fmt.Errorf("%w sender currency not match transaction currency", CurrenciesDoNotMatch)

	}

	receiver, err := s.accounts.GetUser(ctx, transaction.Receiver)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if receiver == nil {
		return fmt.Errorf("account not found: %s", transaction.Receiver)
	}
	if receiver.Currency != transaction.Currency {
		return fmt.Errorf("%w receiver currency not match transaction currency", CurrenciesDoNotMatch)
	}

	err = s.transactions.CreateTransaction(ctx, transaction)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}
	return nil
}
