package main

import (
	"CardService/generated/protobuf"
	"CardService/internal/config"
	"CardService/internal/router"
	"CardService/internal/service/banking_service"
	"CardService/internal/service/card_service"
	"CardService/internal/storage/postgres"
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
)

func main() {

	logger := logrus.New()
	cfg, err := config.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}

	if cfg.LogLevel == "DEBUG" {
		logger.SetLevel(logrus.DebugLevel)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()

	db, err := postgres.NewPostgresRepository(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	// initializing grpc client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(config.GetBankingServiceGrpcURI(), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := protobuf.NewBankingServiceClient(conn)
	banking := banking_service.NewBankingService(client)

	s := card_service.NewCardService(db)
	httpRouter := router.NewRouter(cfg, logger, s, banking)

	errG, gCtx := errgroup.WithContext(ctx)

	errG.Go(func() error {
		logger.Printf("starting http server on port: %s", cfg.ServerPort)
		return httpRouter.Start()
	})

	// clearing resources
	errG.Go(func() error {
		<-gCtx.Done()
		logger.Println("closing grpc connection...")
		return conn.Close()
	})

	errG.Go(func() error {
		<-gCtx.Done()
		logger.Println("closing http client...")
		return httpRouter.Stop(gCtx)
	})
	
	errG.Go(func() error {
		<-gCtx.Done()
		logger.Println("closing database...")
		db.Close()
		return nil
	})

	if err = errG.Wait(); err != nil {
		logger.Printf("exit reason: %s \n", err)
	}
	logger.Println("app shutdown")

}
