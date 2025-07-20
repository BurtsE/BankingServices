package router

import (
	"CardService/internal/config"
	"CardService/internal/service"
	"CardService/pkg/middleware"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const ROUTER_PREFIX = "/card"

type Router struct {
	logger *logrus.Logger
	srv    *http.Server

	service service.CardService
	banking service.IBankingService
}

func NewRouter(cfg *config.Config, logger *logrus.Logger, service service.CardService, banking service.IBankingService) *Router {

	muxRouter := mux.NewRouter().PathPrefix(ROUTER_PREFIX).Subrouter()
	muxRouter.Use(middleware.NewLoggerMiddleware(logger))
	muxRouter.Use(middleware.NewPanicMiddleware(logger))

	srv := &http.Server{
		Handler:      muxRouter,
		Addr:         ":" + cfg.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	r := &Router{
		logger:  logger,
		srv:     srv,
		service: service,
		banking: banking,
	}

	muxRouter.HandleFunc("/create", r.createCardHandler).Methods("POST")
	muxRouter.HandleFunc("/show", r.showCardsHandler).Methods("GET")

	return r
}

func (r *Router) Start() error {
	return r.srv.ListenAndServe()
}

func (r *Router) Stop(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}
