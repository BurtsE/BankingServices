package http_router

import (
	"UserService/internal/config"
	"UserService/internal/service"
	"UserService/pkg/middleware"
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/trace"
)

type Router struct {
	logger    *logrus.Logger
	muxRouter *mux.Router
	service   service.UserService
	srv       *http.Server
	tracer    trace.Tracer
}

func NewRouter(logger *logrus.Logger, cfg *config.Config, service service.UserService, tracer trace.Tracer) *Router {
	r := &Router{
		muxRouter: mux.NewRouter().PathPrefix("/api/v1/user").Subrouter(),
		logger:    logger.WithField("server", "http").Logger,
		service:   service,
		tracer:    tracer,
	}

	r.srv = &http.Server{
		Handler:      r.Handler(),
		Addr:         ":" + cfg.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	r.muxRouter.HandleFunc("/register", r.registerUserHandler).Methods("POST")
	r.muxRouter.HandleFunc("/login", r.loginHandler).Methods("POST")
	r.muxRouter.HandleFunc("/{id:[0-9]+}", r.getUserByIDHandler).Methods("GET")

	r.muxRouter.Use(otelmux.Middleware("User service"), )
	r.muxRouter.Use(middleware.NewLoggerMiddleware(logger))
	r.muxRouter.Use(middleware.NewPanicMiddleware(logger))

	return r
}

func (r *Router) Handler() http.Handler {
	return r.muxRouter
}

func (r *Router) Start() error {
	return r.srv.ListenAndServe()
}

func (r *Router) Stop(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}
