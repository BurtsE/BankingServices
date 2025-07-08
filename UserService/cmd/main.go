package main

import (
	"UserService/internal/config"
	"UserService/internal/server/grpc_server"
	"UserService/internal/server/http_router"
	"UserService/internal/service"
	"UserService/internal/storage/postgres"
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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

	if cfg.LogLevel == "DEBUG" {
		logger.SetLevel(logrus.DebugLevel)
	}

	db, err := postgres.NewPostgresRepository(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	s := service.NewUserService(db, cfg)

	httpRouter := http_router.NewRouter(logger, cfg, s)
	grpcServer := grpc_server.NewGrpcServer(logger, cfg, s)

	errG, gCtx := errgroup.WithContext(ctx)

	errG.Go(func() error {
		logger.Printf("starting grpc server on port: %d", cfg.GrpcPort)
		return grpcServer.Start()
	})

	errG.Go(func() error {
		logger.Printf("starting http server on port: %s", cfg.ServerPort)
		return httpRouter.Start()
	})

	errG.Go(func() error {
		<-gCtx.Done()
		logger.Println("closing database...")
		if db != nil {
			db.Close()
		}
		return nil
	})

	if err = errG.Wait(); err != nil {
		logger.Printf("exit reason: %s \n", err)
	}
	logger.Println("app shutdown")
}
