package router

import (
	"CardService/internal/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Router struct {
	logger    *logrus.Logger
	muxRouter *mux.Router
	service   service.BankingService
	srv       *http.Server
}
