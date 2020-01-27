package main

import (
	"coinsWallet/db/postgres"
	"coinsWallet/service"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
)

func main() {
	configPath := flag.String("config", "config.json", "path to config file")

	conf, err := postgres.Parse(*configPath)
	if err != nil {
		panic(err)
	}

	transactionsRepository, err := postgres.NewTransactionsRepository(conf.Postgres)
	if err != nil {
		panic(err)
	}
	accountsRepository, err := postgres.NewAccountsRepository(conf.Postgres)
	if err != nil {
		panic(err)
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	walletService := service.NewService(transactionsRepository, accountsRepository)
	walletService = service.NewLoggingService(log.With(logger), walletService)
	mux := http.NewServeMux()
	httpLogger := log.With(logger, "component", "http")
	mux.Handle("/wallet/v1/", service.MakeHandler(walletService, httpLogger))

	httpAddr := fmt.Sprintf(":%d", conf.ListenPort)
	server := &http.Server{Addr: httpAddr, Handler: mux}
	logger.Log("msg", "Start listening", "transport", "http", "address", httpAddr)
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err)
		}
	}()
	errs := make(chan error, 1)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
