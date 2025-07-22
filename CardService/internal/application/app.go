package application

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
	"syscall"
)

type Application struct {
	logger         *logrus.Logger
	config         *config.Config
	cardStorage    *postgres.PostgresRepository
	cardService    *card_service.CardService
	bankingService *banking_service.BankingService
	httpServer     *router.Router
	conn           *grpc.ClientConn
	errG           *errgroup.Group
	ctx            context.Context
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
	a.cardStorage, err = postgres.NewPostgresRepository(a.ctx, a.config)
	if err != nil {
		a.logger.Fatal(err)
	}

	// Init services
	a.cardService = card_service.NewCardService(a.cardStorage)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	a.conn, err = grpc.NewClient(config.GetBankingServiceGrpcURI(), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := protobuf.NewBankingServiceClient(a.conn)
	a.bankingService = banking_service.NewBankingService(client)

	a.httpServer = router.NewRouter(a.config, a.logger, a.cardService, a.bankingService)

	// Init routers
	a.httpServer = router.NewRouter(a.config, a.logger, a.cardService, a.bankingService)

	return a
}

func (a *Application) Start(group *errgroup.Group) {
	a.errG = group

	group.Go(func() error {
		a.logger.Printf("starting http server on port: %s", a.config.ServerPort)
		return a.httpServer.Start()
	})

	group.Go(func() error {
		<-a.ctx.Done()
		a.logger.Println("closing grpc connection...")
		return a.conn.Close()
	})

	group.Go(func() error {
		<-a.ctx.Done()
		a.logger.Printf("closing database...")
		a.cardStorage.Close()
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
