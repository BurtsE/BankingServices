package clickhouse

import (
	"bytes"
	"context"
	"fmt"
	"gateway/internal/domain"
	"gateway/internal/storage"
)

func (c *ClickHouseStorage) AddEvent(ctx context.Context, events []domain.EventRequest) error {
	query := prepareQueryStatement(events)
	args := prepareQueryArgs(events)
	err := c.conn.Exec(ctx, query, args...)
	if err != nil {
		return storage.ErrInsertDatabase
	}
	return nil
}

func prepareQueryStatement(events []domain.EventRequest) string {
	buf := bytes.NewBuffer([]byte("INSERT INTO events (uuid, event_type, data, created_at) VALUES "))
	counter := 1
	for range events {
		fmt.Fprintf(buf, "($%d, $%d, $%d, $%d),", counter, counter+1, counter+2, counter+3)
		counter += 4
	}
	buf.Truncate(buf.Len() - 1)
	return buf.String()
}

func prepareQueryArgs(events []domain.EventRequest) []any {
	args := make([]any, 0, len(events)*4)
	for i := range events {
		args = append(args, events[i].ID(), events[i].EventType.Type(), events[i].Data(), events[i].CreatedAt())
	}
	return args
}
