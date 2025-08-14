package requestservice

import (
	"context"
	"fmt"
	"gateway/internal/domain"
	"gateway/internal/storage"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	ErrServiceBisy      = fmt.Errorf("could not process request")
	ErrConnectionClosed = fmt.Errorf("service is closed")
)

// TODO move consts to config

const (
	requestFlushTimeout = time.Millisecond * 20
	batchSize           = 100
)

/*
RequestService batches incoming requests and adds them to storage with scheduled timeout
or when pool is filled. When provided context is cancelled, remaining events are sent to storage.
New events will not be accepted.
*/
type RequestService struct {
	logger      *logrus.Logger
	storage     storage.EventStorage
	requestChan chan domain.EventRequest
	pool        []domain.EventRequest
	ctx         context.Context
	cancelFunc  context.CancelFunc
}

func NewRequestService(logger *logrus.Logger, storage storage.EventStorage) *RequestService {
	return &RequestService{
		logger:      logger,
		storage:     storage,
		requestChan: make(chan domain.EventRequest, batchSize),
		pool:        make([]domain.EventRequest, 0, batchSize),
	}
}

/*
Starts requestService
*/
func (r *RequestService) Start(ctx context.Context) {
	r.ctx, r.cancelFunc = context.WithCancel(ctx)
	defer close(r.requestChan)
	defer r.cancelFunc() // cancel context before closing request channel in order to prevent panic from AddEvent
	ticker := time.NewTicker(requestFlushTimeout)

	for {
		select {
		case <-r.ctx.Done():
			return
		case event := <-r.requestChan:
			r.pool = append(r.pool, event)
			if len(r.pool) != batchSize { // prevent adding to storage if pool is not filled
				continue
			}
		case <-ticker.C:
			if len(r.pool) == 0 {
				continue
			}
		}
		err := r.storage.AddEvents(ctx, r.pool)
		if err != nil {
			r.logger.Errorf("error adding events to database: %v: %v", err, r.pool)
		}
		r.logger.Debugf("added %d events to db", len(r.pool))
		r.pool = r.pool[:0]
	}
}

func (r *RequestService) Stop() {
	r.cancelFunc()
}

func (r *RequestService) AddEvent(ctx context.Context, event domain.EventRequest) error {
	select {
	case <-ctx.Done():
		return ErrServiceBisy
	case <-r.ctx.Done():
		return ErrConnectionClosed
	case r.requestChan <- event:
		return nil
	}
}
