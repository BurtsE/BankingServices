package main

import (
	"BankingService/internal/config"
	"BankingService/internal/router"
	"BankingService/internal/service"
	"BankingService/internal/storage/postgres"
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"

	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	logger := logrus.New()

	cfg, err := config.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := postgres.NewPostgresRepository(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	s := service.NewBankingService(db)

	rtr := router.NewRouter(logger, cfg, s)

	errG, gCtx := errgroup.WithContext(ctx)

	errG.Go(func() error {
		logger.Printf("starting server on port: %s", cfg.ServerPort)
		return rtr.Start()
	})

	errG.Go(func() error {
		<-gCtx.Done()
		log.Println("closing database...")
		if db != nil {
			db.Close()
		}
		return nil
	})

	if err = errG.Wait(); err != nil {
		log.Printf("exit reason: %s \n", err)
	}
	log.Println("app shutdown")
}
