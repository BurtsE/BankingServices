package main

import (
	"context"
	"gateway/generated/protobuf"
	"gateway/internal/cache/redis"
	"gateway/internal/config"
	"gateway/internal/metrics_server"
	"gateway/internal/router"
	"gateway/internal/service/user_service"
	"gateway/pkg/metrics"
	"gateway/pkg/tracing"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"syscall"
)

func main() {

	// context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	// initiating configuration for the app
	logger := logrus.New()

	cfg, err := config.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}

	if cfg.LogLevel == "DEBUG" {
		logger.SetLevel(logrus.DebugLevel)
	}

	// initializing jaeger tracer (global), http router should use middleware
	jaegerURL := config.GetJaegerUrl()
	tracerProvider, err := tracing.InitTracer(jaegerURL, "Gateway service")
	if err != nil {
		logger.Fatal(err)
	}

	// connecting to redis
	logger.Printf("connecting to redis with address %s:%s", cfg.Redis.Host, cfg.Redis.Port)
	cache, err := redis.NewRedisCache(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	// initializing grpc client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(config.GetUserServiceGrpcURI(), opts...)
	if err != nil {
		logger.Fatalf("fail to dial: %v", err)
	}

	client := protobuf.NewUserServiceClient(conn)
	userService := user_service.NewUserService(client)

	// Init prometheus metric registry
	registry := prometheus.NewRegistry()
	m := metrics.NewMetrics(registry)

	// initializing http router
	rtr := router.NewRouter(cfg, logger, cache, userService, m, tracerProvider.Tracer("gateway tracer name"))

	// initializing metric routes
	metricServer := metrics_server.NewMetricsServer(registry)

	errG, gCtx := errgroup.WithContext(ctx)

	errG.Go(func() error {
		logger.Printf("starting http server on port: %s", cfg.ServerPort)
		return rtr.Start()
	})

	errG.Go(func() error {
		logger.Printf("starting metric server on port: %s", config.GetPrometheusPort())
		return metricServer.Start()
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
		return rtr.Stop(gCtx)
	})

	errG.Go(func() error {
		<-gCtx.Done()
		logger.Println("closing metric server...")
		return metricServer.Stop(gCtx)
	})

	errG.Go(func() error {
		<-gCtx.Done()
		logger.Println("closing jaeger connection...")
		return tracerProvider.Shutdown(context.Background())
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
