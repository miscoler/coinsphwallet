package service

import (
	"coinsWallet/account"
	"coinsWallet/transaction"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type createTransactionRequest struct {
	Transaction transaction.Transaction
}

type createTransactionResponse struct {
	Err error `json:"err,omitempty"`
}

func makeCreateTransactionEndpoint(s WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*createTransactionRequest)
		err := s.CreateTransaction(ctx, req.Transaction)
		if err != nil {
			return nil, err
		}
		return &createTransactionResponse{Err: err}, nil
	}
}

type getTransactionsResponse struct {
	Transactions []*transaction.Transaction `json:"transactions"`
	Err          error                      `json:"err,omitempty"`
}

func makeGetTransactionsEndpoint(s WalletService) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		transactions, err := s.GetTransactions(ctx)
		if err != nil {
			return nil, err
		}
		if transactions == nil {
			transactions = []*transaction.Transaction{}
		}
		return &getTransactionsResponse{Transactions: transactions}, nil
	}
}

type getAccountsResponse struct {
	Accounts []*account.Account `json:"accounts"`
	Err      error              `json:"err,omitempty"`
}

func makeGetAccountsEndpoint(s WalletService) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		accounts, err := s.GetAccounts(ctx)
		if err != nil {
			return nil, err
		}
		if accounts == nil {
			accounts = []*account.Account{}
		}
		return &getAccountsResponse{Accounts: accounts}, nil
	}
}
