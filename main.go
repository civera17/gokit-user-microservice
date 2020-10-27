package main

import (
	"account"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var httpAddr = flag.String("http",":8080","http listen address")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var db *sql.DB
	{
		var err error

		db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/gokittest")
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}
	flag.Parse()
	ctx := context.Background()
	var service account.Service
	{
		repository := account.NewRepo(db,logger)
		service = account.NewService(repository, logger)
	}

	errors := make(chan error)

	go func() {
		c := make(chan  os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errors <- fmt.Errorf("%s", <-c)
	}()

	endpoints := account.MakeEndpoints(service)

	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := account.NewHTTPServer(ctx, endpoints)
		errors <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <- errors)
}
