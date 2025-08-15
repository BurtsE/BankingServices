package requestservice

import (
	"context"
	"gateway/internal/domain"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) AddEvents(ctx context.Context, events []domain.EventRequest) error {
	args := m.Called(ctx, events)
	return args.Error(0)
}

func TestRequestService(t *testing.T) {
	ctx := context.Background()
	mockStorage := &mockStorage{}
	mockStorage.On("AddEvents", mock.Anything, mock.Anything).
		Return(nil)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	service := NewRequestService(logger, mockStorage)

	clientsNumber := 10000
	requestTimeout := time.Millisecond * 100
	wg := new(sync.WaitGroup)
	wg.Add(clientsNumber)

	go service.Start(ctx)
	time.Sleep(time.Second)

	for range clientsNumber {
		go func() {
			defer wg.Done()
			clientCtx, cancel := context.WithTimeout(ctx, requestTimeout)
			defer cancel()
			err := service.AddEvent(clientCtx, domain.EventRequest{})
			assert.Nil(t, err, "should not return error")
		}()
	}
	wg.Wait()
	time.Sleep(time.Second)
	service.Stop()
	err := service.AddEvent(ctx, domain.EventRequest{})
	assert.Equal(t, err, ErrConnectionClosed, "should not accept new events")
}
