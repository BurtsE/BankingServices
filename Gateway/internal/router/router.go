package router

import (
	"context"
	"gateway/internal/cache"
	"gateway/internal/config"
	"gateway/internal/service"
	"gateway/pkg/metrics"
	"gateway/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const ROUTER_PREFIX = "/api/v1"

const caching_duration = time.Hour * 1

const (
	user_prefix    = "user"
	banking_prefix = "account"
	cards_prefix   = "card"
)

type Router struct {
	logger *logrus.Logger
	proxy  *httputil.ReverseProxy
	srv    *http.Server

	userService service.IUserService
	mapping     map[string]*url.URL
	tokenCache  cache.Cache
	metrics     *metrics.Metrics
}

func NewRouter(cfg *config.Config, logger *logrus.Logger, cache cache.Cache, userService service.IUserService,
	metrics *metrics.Metrics) *Router {

	rtr := &Router{
		logger:      logger,
		tokenCache:  cache,
		userService: userService,
	}

	muxRouter := mux.NewRouter().PathPrefix(ROUTER_PREFIX).Subrouter()

	rtr.srv = &http.Server{
		Handler:      muxRouter,
		Addr:         ":" + cfg.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	proxy := &httputil.ReverseProxy{
		Director:       rtr.director,
		ErrorHandler:   rtr.errorHandler,
		ModifyResponse: rtr.modifyResponse,
	}

	rtr.proxy = proxy
	rtr.InitServiceMapping()

	// Setup prometheus metrics
	rtr.metrics = metrics

	// register proxy handlers
	muxRouter.HandleFunc("/user/{*}", rtr.UserServiceHandler)
	muxRouter.HandleFunc("/account/{*}", rtr.BankingServiceHandler).Methods("GET", "POST")
	muxRouter.HandleFunc("/card/{*}", rtr.BankingServiceHandler).Methods("GET", "POST")
	muxRouter.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("pong")) }))

	// add middleware
	muxRouter.Use(middleware.NewLoggerMiddleware(logger))
	muxRouter.Use(middleware.NewPanicMiddleware(logger))

	return rtr
}

func (r *Router) InitServiceMapping() {
	r.mapping = make(map[string]*url.URL)

	uri, err := url.Parse(config.GetUserServiceHttpURI())
	if err != nil {
		r.logger.Fatal(err)
	}
	r.mapping[user_prefix] = uri

	uri, err = url.Parse(config.GetBankingServiceURI())
	if err != nil {
		r.logger.Fatal(err)
	}
	r.mapping[banking_prefix] = uri

	uri, err = url.Parse(config.GetCardServiceURI())
	if err != nil {
		r.logger.Fatal(err)
	}
	r.mapping[cards_prefix] = uri
}

func (r *Router) Start() error {
	return r.srv.ListenAndServe()
}

func (r *Router) Stop(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}
