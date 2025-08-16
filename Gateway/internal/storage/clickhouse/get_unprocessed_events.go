package clickhouse

import (
	"context"
	"fmt"
	"gateway/internal/domain"
	"time"

	"github.com/google/uuid"
)

func (c *ClickHouseStorage) GetUnprocessedEvents(ctx context.Context) ([]domain.EventRequest, error) {
	query := `
		SELECT uuid, event_type, data, created_at
		FROM(
			SELECT uuid, event_type, data, created_at, rowNUmber() OVER ()(
				PARTITION BY uuid
				ORDER BY created_at DESC
			) as rn
			FROM events	
		)
		WHERE rn = 1 AND status = $1
	`
	rows, err := c.conn.Query(ctx, query, domain.InProgress.Status())
	if err != nil {
		return nil, fmt.Errorf("failed to query unprocessed events: %w", err)
	}
	defer rows.Close()

	var events []domain.EventRequest
	for rows.Next() {
		var (
			id        uuid.UUID
			eventType string
			data      []byte
			createdAt time.Time
		)
		if err := rows.Scan(&id, &eventType, &data, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		event, err := domain.NewEventRequest(id, createdAt, data, eventType)
		if err != nil {
			return nil, fmt.Errorf("failed to create event: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}
