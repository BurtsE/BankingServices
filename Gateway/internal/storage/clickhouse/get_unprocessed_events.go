package clickhouse

import (
	"context"
	"gateway/internal/domain"
)

func (c *ClickHouseStorage) GetUnprocessedEvents(ctx context.Context, events []domain.EventRequest) error {
	// query := `
	// 	SELECT FROM events
	// `
	return nil
}
