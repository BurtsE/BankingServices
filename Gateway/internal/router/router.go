package router

import (
	"context"
	"gateway/internal/cache"
	"gateway/internal/config"
	"gateway/internal/service"
	"gateway/pkg/metrics"
	"gateway/pkg/middleware"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/trace"
)

const RouterPrefix = "/api/v1"

const cachingDuration = time.Hour * 1

type Router struct {
	logger *logrus.Logger
	srv    *http.Server
	cfg    *config.Config

	userService    service.IUserService
	requestService service.EventService
	tokenCache     cache.Cache

	metrics *metrics.Metrics
	tracer  trace.Tracer
}

func NewRouter(cfg *config.Config, logger *logrus.Logger, cache cache.Cache, userService service.IUserService,
	metrics *metrics.Metrics, tracer trace.Tracer) *Router {

	rtr := &Router{
		logger:      logger,
		cfg:         cfg,
		tokenCache:  cache,
		userService: userService,
		tracer:      tracer,
	}

	muxRouter := mux.NewRouter().PathPrefix(RouterPrefix).Subrouter()

	rtr.srv = &http.Server{
		Handler:      muxRouter,
		Addr:         ":" + cfg.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Setup prometheus metrics
	rtr.metrics = metrics

	// register user handlers
	muxRouter.HandleFunc("/user/register", rtr.RegisterUserHandler).Methods("POST")

	muxRouter.HandleFunc("/user/login", rtr.UserLoginHandler).Methods("POST")

	// register account handlers
	// muxRouter.HandleFunc("/account/create", rtr.AccountCreationHandler).Methods("POST")
	// muxRouter.HandleFunc("/account/deposit", rtr.AccountDepositHandler).Methods("POST")
	// muxRouter.HandleFunc("/account/withdraw", rtr.AccountWithdrawHandler).Methods("POST")
	// muxRouter.HandleFunc("/account/transfer", rtr.AccountTransferHandler).Methods("POST")

	// register card handlers
	// muxRouter.HandleFunc("/card/show/{accountId}", rtr.ShowCardHandler).Methods("GET")
	// muxRouter.HandleFunc("/card/issue", rtr.IssueCardHandler).Methods("POST")
	// muxRouter.HandleFunc("/card/block", rtr.BlockCardHandler).Methods("POST")

	// handler for checking request status
	// muxRouter.Handle("/request_status/{event_id}", rtr.EventStatusHandler).Methods("GET")

	// healthcheck handler
	muxRouter.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("pong")) }))

	// add middleware
	muxRouter.Use(otelmux.Middleware("Gateway"))
	muxRouter.Use(middleware.NewLoggerMiddleware(logger))
	muxRouter.Use(middleware.NewPanicMiddleware(logger))

	return rtr
}

func (r *Router) Start() error {
	return r.srv.ListenAndServe()
}

func (r *Router) Stop(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}
