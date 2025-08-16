package requestpusher

import (
	"context"
	"encoding/json"
	"gateway/internal/config"
	"gateway/internal/domain"
	"gateway/internal/storage"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// TODO moveto config
const requestDelay = time.Second

// KafkaWriter defines the interface for writing messages to Kafka.
// It's implemented by *kafka.Conn and can be mocked for tests.
type KafkaWriter interface {
	WriteMessages(msgs ...kafka.Message) (int, error)
	Close() error
}

type RequestPusher struct {
	storage storage.EventStorageReader
	logger  *logrus.Logger
	pusher  KafkaWriter
	ticker  *time.Ticker
	cfg     *config.Config
}

func NewRequestPusher(cfg *config.Config, storage storage.EventStorageReader, logger *logrus.Logger, pusher KafkaWriter) *RequestPusher {
	return &RequestPusher{
		storage: storage,
		logger:  logger,
		pusher:  pusher,
		ticker:  time.NewTicker(requestDelay),
		cfg:     cfg,
	}
}

func (r *RequestPusher) Start(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-r.ticker.C:
				events, err := r.storage.GetUnprocessedEvents(context.Background())
				if err != nil {
					r.logger.Errorf("error getting unprocessed events: %v", err)
					continue
				}
				r.processEvents(events)
			}
		}
	}(ctx)
}

func (r *RequestPusher) Stop() error {
	return r.pusher.Close()
}

// A struct representing the message payload for Kafka
type kafkaEventPayload struct {
	ID        uuid.UUID       `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	Data      json.RawMessage `json:"data"`
	Type      string          `json:"type"`
}

// send messages to kafka,
func (r *RequestPusher) processEvents(events []domain.EventRequest) {
	if len(events) == 0 {
		return
	}
	messages := make([]kafka.Message, 0, len(events))
	for _, event := range events {
		payload := kafkaEventPayload{
			ID:        event.ID(),
			CreatedAt: event.CreatedAt(),
			Data:      event.Data(),
			Type:      event.Type(),
		}
		messageBytes, err := json.Marshal(payload)
		if err != nil {
			r.logger.Errorf("error marshalling event %v: %v", event.ID(), err)
			continue
		}
		topic := r.eventTopic(event)
		if topic == "" {
			r.logger.Errorf("error getting topic for event %v: %v", event.ID(), domain.ErrInvalidEventType)
			continue
		}
		msg := kafka.Message{
			Topic: topic,
			Value: messageBytes,
		}
		messages = append(messages, msg)
	}
	if _, err := r.pusher.WriteMessages(messages...); err != nil {
		r.logger.Errorf("failed to write messages to kafka: %v", err)
	}
}

// configure topics based on event types
func (r *RequestPusher) eventTopic(event domain.EventRequest) string {
	switch event.Type() {
	case "create_user":
		return r.cfg.Kafka.UserTopic
	case "create_account":
		return r.cfg.Kafka.AccountTopic
	case "deposit_account":
		return r.cfg.Kafka.AccountTopic
	default:
		return ""
	}
}
