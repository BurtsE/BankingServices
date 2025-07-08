package router

import (
	"BankingService/internal/config"
	"BankingService/internal/service"
	"BankingService/pkg/middleware"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type Router struct {
	logger    *logrus.Logger
	muxRouter *mux.Router
	service   service.BankingService
	srv       *http.Server
}

// NewRouter — конструктор роутера
func NewRouter(logger *logrus.Logger, cfg *config.Config, service service.BankingService) *Router {
	r := &Router{
		logger:    logger,
		muxRouter: mux.NewRouter().PathPrefix("/api/v1/account").Subrouter(),
		service:   service,
	}

	r.srv = &http.Server{
		Handler:      r.Handler(),
		Addr:         ":" + cfg.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	r.muxRouter.HandleFunc("/create", r.createAccountHandler).Methods("POST")
	r.muxRouter.HandleFunc("/deposit", r.depositHandler).Methods("POST")
	r.muxRouter.HandleFunc("/withdraw", r.withdrawHandler).Methods("POST")
	r.muxRouter.HandleFunc("/transfer", r.transferHandler).Methods("POST")

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

func GetAccountIDFromURI(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	accountIDStr := vars["id"]
	if accountIDStr == "" {
		return 0, fmt.Errorf(" missing account id")
	}
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(" could not parce account id")
	}
	return accountID, nil
}
