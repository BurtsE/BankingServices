package application

import (
	"UserService/internal/config"
	"UserService/internal/servers/grpc_server"
	"UserService/internal/servers/http_router"
	"UserService/internal/service"
	"UserService/internal/storage/postgres"
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	logger      *logrus.Logger
	userStorage *postgres.PostgresRepository
	service     service.UserService
	httpServer  *http_router.Router
	grpcServer  *grpc_server.Server
	config      *config.Config
	errG        *errgroup.Group
	ctx         context.Context
}

func NewApp(ctx context.Context) *Application {
	var err error
	a := &Application{ctx: ctx}
	a.logger = logrus.New()

	// Init configuration from file
	a.config, err = config.InitConfig()
	if err != nil {
		a.logger.Fatal(err)
	}

	// Set logger level
	if a.config.LogLevel == "DEBUG" {
		a.logger.SetLevel(logrus.DebugLevel)
	}

	// Init repository
	a.userStorage, err = postgres.NewPostgresRepository(a.ctx, a.config)
	if err != nil {
		a.logger.Fatal(err)
	}

	// Init services
	s := service.NewUserService(a.userStorage, a.config)

	// Init routers
	a.httpServer = http_router.NewRouter(a.logger, a.config, s)
	a.grpcServer = grpc_server.NewGrpcServer(a.logger, a.config, s)

	return a
}

func (a *Application) Start(group *errgroup.Group) {
	a.errG = group

	group.Go(func() error {
		a.logger.Printf("starting grpc server on port: %d", a.config.GrpcPort)
		return a.grpcServer.Start()
	})

	group.Go(func() error {
		a.logger.Printf("starting http server on port: %s", a.config.ServerPort)
		return a.httpServer.Start()
	})

	// Cleaning resources
	group.Go(func() error {
		<-a.ctx.Done()
		a.logger.Println("stopping grpc server...")
		return a.grpcServer.Stop()
	})

	group.Go(func() error {
		<-a.ctx.Done()
		a.logger.Println("stopping http server...")
		return a.httpServer.Stop(context.Background())
	})

	group.Go(func() error {
		<-a.ctx.Done()
		a.logger.Println("closing database...")
		a.userStorage.Close()
		return nil
	})
}

// AwaitTermination makes the program wait for the signal termination
// Valid signal termination (SIGINT, SIGTERM)
func (a *Application) AwaitTermination(shutdownHook func()) {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM)
	<-interruptSignal
	shutdownHook()
	if err := a.errG.Wait(); err != nil {
		a.logger.Printf("exit reason: %s", err)
	}
	a.logger.Println("app shutdown")
}
