package requestpusher

import (
	"context"
	"encoding/json"
	"errors"
	"gateway/internal/config"
	"gateway/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockStorage is a mock for the storage.EventStorageReader interface.
type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) GetUnprocessedEvents(ctx context.Context) ([]domain.EventRequest, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.EventRequest), args.Error(1)
}

// mockKafkaWriter is a mock for the KafkaWriter interface.
type mockKafkaWriter struct {
	mock.Mock
}

func (m *mockKafkaWriter) WriteMessages(msgs ...kafka.Message) (int, error) {
	// The mock library doesn't handle variadic arguments well when they are the only argument.
	// We'll pass it as a slice.
	args := m.Called(msgs)
	return args.Int(0), args.Error(1)
}

func (m *mockKafkaWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestRequestPusher_processEvents(t *testing.T) {
	logger := logrus.New()
	logger.SetReportCaller(true)

	logger.SetLevel(logrus.ErrorLevel) // Keep test output clean

	cfg := &config.Config{
		Kafka: config.Kafka{
			UserTopic:    "user-events",
			AccountTopic: "account-events",
		},
	}

	t.Run("success - process and send multiple events", func(t *testing.T) {
		// Arrange
		mockWriter := new(mockKafkaWriter)
		pusher := NewRequestPusher(cfg, nil, logger, mockWriter) // storage is not used by processEvents

		userEvent, _ := domain.NewEventRequest(uuid.New(), time.Now(), []byte(`{"username":"test"}`), "create_user")
		accountEvent, _ := domain.NewEventRequest(uuid.New(), time.Now(), []byte(`{"currency":"USD"}`), "create_account")
		events := []domain.EventRequest{userEvent, accountEvent}

		// We expect WriteMessages to be called once with a slice of 2 messages.
		mockWriter.On("WriteMessages", mock.AnythingOfType("[]kafka.Message")).Return(len(events), nil).Run(func(args mock.Arguments) {
			msgs := args.Get(0).([]kafka.Message)
			require.Len(t, msgs, 2)

			// Check first message (user event)
			assert.Equal(t, cfg.Kafka.UserTopic, msgs[0].Topic)
			var payload kafkaEventPayload
			err := json.Unmarshal(msgs[0].Value, &payload)
			require.NoError(t, err)
			assert.Equal(t, userEvent.ID(), payload.ID)
			assert.Equal(t, "create_user", payload.Type)
			assert.JSONEq(t, `{"username":"test"}`, string(payload.Data))

			// Check second message (account event)
			assert.Equal(t, cfg.Kafka.AccountTopic, msgs[1].Topic)
			err = json.Unmarshal(msgs[1].Value, &payload)
			require.NoError(t, err)
			assert.Equal(t, accountEvent.ID(), payload.ID)
			assert.Equal(t, "create_account", payload.Type)
			assert.JSONEq(t, `{"currency":"USD"}`, string(payload.Data))
		})

		// Act
		pusher.processEvents(events)

		// Assert
		mockWriter.AssertExpectations(t)
	})

	t.Run("kafka write fails", func(t *testing.T) {
		mockWriter := new(mockKafkaWriter)
		pusher := NewRequestPusher(cfg, nil, logger, mockWriter)
		event, _ := domain.NewEventRequest(uuid.New(), time.Now(), []byte(`{}`), "create_user")
		events := []domain.EventRequest{event}

		kafkaErr := errors.New("kafka is down")
		mockWriter.On("WriteMessages", mock.AnythingOfType("[]kafka.Message")).Return(0, kafkaErr)

		pusher.processEvents(events)

		mockWriter.AssertExpectations(t)
	})
}

