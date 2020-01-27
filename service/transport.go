package service

import (
	"coinsWallet/transaction"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(s WalletService, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	sendTransactionHandler := kithttp.NewServer(
		makeCreateTransactionEndpoint(s),
		decodeSendTransactionRequest,
		encodeResponse,
		opts...,
	)
	getTransactionsHandler := kithttp.NewServer(
		makeGetTransactionsEndpoint(s),
		decodeGetAllTransactionsRequest,
		encodeResponse,
		opts...,
	)
	getAccountsHandler := kithttp.NewServer(
		makeGetAccountsEndpoint(s),
		decodeGetAllAccountsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/wallet/v1/transaction", sendTransactionHandler).Methods("POST")
	r.Handle("/wallet/v1/transaction", getTransactionsHandler).Methods("GET")
	r.Handle("/wallet/v1/account", getAccountsHandler).Methods("GET")

	return r
}

type decodingError struct {
	Details string
}

func (de *decodingError) Error() string {
	return de.Details
}

func decodeSendTransactionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	sender := r.PostFormValue("sender")
	receiver := r.PostFormValue("receiver")
	if sender == "" || receiver == "" {
		return nil, &decodingError{"sender and receiver accounts needed"}
	}

	amount, err := strconv.ParseFloat(r.PostFormValue("amount"), 64)
	if err != nil {
		return nil, &decodingError{"invalid amount value"}
	}

	currency := r.PostFormValue("currency")
	if currency == "" {
		return nil, &decodingError{"currency is needed"}
	}

	return &createTransactionRequest{Transaction: transaction.Transaction{
		Sender:   sender,
		Receiver: receiver,
		Currency: currency,
		Amount:   amount,
	}}, nil
}

func decodeGetAllTransactionsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetAllAccountsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	var httpStatusCode int
	switch response.(type) {
	case *createTransactionResponse:
		httpStatusCode = http.StatusCreated
	default:
		httpStatusCode = http.StatusOK
	}
	w.WriteHeader(httpStatusCode)
	return json.NewEncoder(w).Encode(response)
}
