package metrics_server

import (
	"context"
	"gateway/internal/config"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type MetricsServer struct {
	srv *http.Server
}

func NewMetricsServer(reg *prometheus.Registry) *MetricsServer {
	metricsRouter := mux.NewRouter()

	port := config.GetPrometheusPort()

	metricsRouter.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      metricsRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	
	return &MetricsServer{srv: srv}
}

func (m *MetricsServer) Start() error {
	return m.srv.ListenAndServe()
}

func (m *MetricsServer) Stop(ctx context.Context) error {
	return m.srv.Shutdown(ctx)
}
