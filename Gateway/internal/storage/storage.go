package storage

import (
	"context"
	"fmt"
	"gateway/internal/domain"
)

var ErrInsertDatabase = fmt.Errorf("error inserting into database")

type EventStorageWriter interface {
	AddEvents(context.Context, []domain.EventRequest) error
}

type EventStorageReader interface {
	GetUnprocessedEvents(context.Context) ([]domain.EventRequest, error)
}
