package router

import (
	"UserService/internal/config"
	"UserService/internal/service"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Router struct {
	logger    *logrus.Logger
	muxRouter *mux.Router
	service   service.UserService
	srv       *http.Server
}

// NewRouter — конструктор роутера
func NewRouter(logger *logrus.Logger, cfg *config.Config, service service.UserService) *Router {
	r := &Router{
		muxRouter: mux.NewRouter().PathPrefix("/api/v1").Subrouter(),
		logger:    logger,
		service:   service,
	}
	r.srv = &http.Server{
		Handler:      r.Handler(),
		Addr:         ":" + cfg.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
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
