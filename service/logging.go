package service

import (
	"coinsWallet/account"
	"coinsWallet/transaction"
	"context"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	WalletService
}

func NewLoggingService(logger log.Logger, s WalletService) WalletService {
	return &loggingService{logger, s}
}

func (s *loggingService) CreateTransaction(ctx context.Context, transaction transaction.Transaction) error {
	s.logger.Log(
		"method", "create_transaction",
		"sender", transaction.Sender,
		"receiver", transaction.Receiver,
		"amount", transaction.Amount,
		"currency", transaction.Currency,
	)
	return s.WalletService.CreateTransaction(ctx, transaction)
}

func (s *loggingService) GetTransactions(ctx context.Context) ([]*transaction.Transaction, error) {
	s.logger.Log(
		"method", "get_payments",
	)
	return s.WalletService.GetTransactions(ctx)
}

func (s *loggingService) GetAllAccounts(ctx context.Context) ([]*account.Account, error) {
	s.logger.Log(
		"method", "get_accounts",
	)
	return s.WalletService.GetAccounts(ctx)
}
