package main

import (
	"context"
	"gateway/internal/cache/redis"
	"gateway/internal/config"
	"gateway/internal/router"
	"gateway/internal/service/user_service"
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

	cache := redis.NewRedisCache(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	userService := &user_service.UserService{}

	rtr := router.NewRouter(cfg, logger, cache, userService)

	errG, gCtx := errgroup.WithContext(ctx)

	errG.Go(func() error {
		logger.Printf("starting server on port: %s", cfg.ServerPort)
		return rtr.Start()
	})

	errG.Go(func() error {
		<-gCtx.Done()
		logger.Println("closing database...")
		return cache.Close(gCtx)
	})

	if err = errG.Wait(); err != nil {
		logger.Printf("exit reason: %s \n", err)
	}
	logger.Println("app shutdown")
}
