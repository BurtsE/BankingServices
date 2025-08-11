package router

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"gateway/internal/config"
	"gateway/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func SetupEnvironmentVariables() {
	os.Setenv("USER_SERVICE_HTTP_URI", "test_uri")
	os.Setenv("BANKING_SERVICE_URI", "test_uri")
	os.Setenv("CARD_SERVICE_URI", "test_uri")
}

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Validate(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

type mockCache struct {
	mock.Mock
}

func (m *mockCache) Get(ctx context.Context, key string) (string, bool, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *mockCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *mockCache) Save(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func TestNewRouter(t *testing.T) {
	SetupEnvironmentVariables()

	mockUserSvc := new(mockUserService)
	mockCache := new(mockCache)
	logger := logrus.New()
	registry := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(registry)
	tracer := noop.NewTracerProvider().Tracer("test")

	cfg := &config.Config{
		ServerPort: "8080",
	}

	router := NewRouter(cfg, logger, mockCache, mockUserSvc, metrics, tracer)

	assert.NotNil(t, router)
	assert.NotNil(t, router.proxy)
	assert.NotNil(t, router.srv)
}

func TestRouter_StartStop(t *testing.T) {
	SetupEnvironmentVariables()
	mockUserSvc := new(mockUserService)
	mockCache := new(mockCache)
	logger := logrus.New()
	registry := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(registry)
	tracer := noop.NewTracerProvider().Tracer("test")

	cfg := &config.Config{
		ServerPort: "8080",
	}

	router := NewRouter(cfg, logger, mockCache, mockUserSvc, metrics, tracer)

	// Test starting and stopping the router
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := router.Start()
		assert.Error(t, err) // Should error when we shutdown
	}()

	time.Sleep(100 * time.Millisecond) // Give server time to start
	err := router.Stop(ctx)
	assert.NoError(t, err)
}

func TestRouter_ProxyHandlers(t *testing.T) {
	SetupEnvironmentVariables()

	mockUserSvc := new(mockUserService)
	mockCache := new(mockCache)
	logger := logrus.New()
	registry := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(registry)
	tracer := noop.NewTracerProvider().Tracer("test")

	cfg := &config.Config{
		ServerPort: "8080",
	}

	router := NewRouter(cfg, logger, mockCache, mockUserSvc, metrics, tracer)

	testCases := []struct {
		name       string
		path       string
		method     string
		statusCode int
	}{
		{"UserService", "/api/v1/user/test", "GET", http.StatusServiceUnavailable},
		{"BankingService", "/api/v1/account/test", "GET", http.StatusUnauthorized},
		{"CardService", "/api/v1/card/test", "GET", http.StatusUnauthorized},
		{"Ping", "/api/v1/ping", "GET", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			router.srv.Handler.ServeHTTP(w, req)

			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}

func TestRouter_ErrorHandler(t *testing.T) {
	mockUserSvc := new(mockUserService)
	mockCache := new(mockCache)
	logger := logrus.New()
	registry := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(registry)
	tracer := trace.NewNoopTracerProvider().Tracer("test")

	cfg := &config.Config{
		ServerPort: "8080",
	}

	router := NewRouter(cfg, logger, mockCache, mockUserSvc, metrics, tracer)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	err := errors.New("test error")

	router.errorHandler(w, req, err)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestRouter_ModifyResponse(t *testing.T) {
	mockUserSvc := new(mockUserService)
	mockCache := new(mockCache)
	logger := logrus.New()
	registry := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(registry)
	tracer := trace.NewNoopTracerProvider().Tracer("test")

	cfg := &config.Config{
		ServerPort: "8080",
	}

	router := NewRouter(cfg, logger, mockCache, mockUserSvc, metrics, tracer)

	res := &http.Response{
		Header: make(http.Header),
	}

	err := router.modifyResponse(res)

	assert.NoError(t, err)
	assert.Equal(t, "GoG", res.Header.Get("gateway"))
}
