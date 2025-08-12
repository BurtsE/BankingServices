package storage

import (
	"context"
	"fmt"
	"gateway/internal/domain"
)

var ErrInsertDatabase = fmt.Errorf("error inserting into database")

type EventStorage interface {
	AddEvents(context.Context, []domain.Event) error
}
