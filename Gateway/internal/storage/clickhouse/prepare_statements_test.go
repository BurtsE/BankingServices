package clickhouse

import (
	"gateway/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPrepareQueryStatements(t *testing.T) {
	tests := []struct {
		name   string
		events []domain.EventRequest
		result string
	}{
		{
			name:   "empty slice",
			events: []domain.EventRequest{},
			result: "INSERT INTO events (uuid, event_type, data, created_at) VALUES",
		},
		{
			name: "single event",
			events: []domain.EventRequest{
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: "INSERT INTO events (uuid, event_type, data, created_at) VALUES ($1, $2, $3, $4)",
		},
		{
			name: "multiple events",
			events: []domain.EventRequest{
				{
					EventType: domain.CreateUserEvent,
				},
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: "INSERT INTO events (uuid, event_type, data, created_at) VALUES ($1, $2, $3, $4),($5, $6, $7, $8)",
		},
	}

	for _, test := range tests {
		query := prepareQueryStatement(test.events)
		assert.Equal(t, test.result, query, "should be equal")
	}
}

func TestPrepareQueryArgs(t *testing.T) {
	tests := []struct {
		name   string
		events []domain.EventRequest
		result []any
	}{
		{
			name:   "empty slice",
			events: []domain.EventRequest{},
			result: []any{},
		},
		{
			name: "single event",
			events: []domain.EventRequest{
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: []any{
				any(uuid.UUID{}), domain.CreateUserEvent.Type(), []byte(nil), time.Time{},
			},
		},
		{
			name: "multiple events",
			events: []domain.EventRequest{
				{
					EventType: domain.CreateUserEvent,
				},
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: []any{
				any(uuid.UUID{}), domain.CreateUserEvent.Type(), []byte(nil), time.Time{},
				any(uuid.UUID{}), domain.CreateUserEvent.Type(), []byte(nil), time.Time{},
			},
		},
	}

	for _, test := range tests {
		args := prepareQueryArgs(test.events)
		if len(test.result) != len(args) {
			t.Fatalf("%s failed: length does not match:", test.name)
		}
		counter := 0
		for counter < len(test.result) {
			assert.Equal(t, test.result[counter], args[counter], "ids should be equal")
			assert.Equal(t, test.result[counter+1], args[counter+1], "types should be equal")
			assert.Equal(t, test.result[counter+2], args[counter+2], "data should be equal")
			assert.Equal(t, test.result[counter+3], args[counter+3], "creation time should be equal")

			counter += 4
		}
	}
}

// TODO add checks for time and data
func TestPrepareQueryArgsFail(t *testing.T) {
	tests := []struct {
		name   string
		events []domain.EventRequest
		result []any
	}{
		{
			name: "single event",
			events: []domain.EventRequest{
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: []any{
				uuid.New(), "delete_user", []byte(nil), time.Time{},
			},
		},
		{
			name: "multiple events",
			events: []domain.EventRequest{
				{
					EventType: domain.CreateUserEvent,
				},
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: []any{
				uuid.New(), "", []byte(nil), time.Time{},
				uuid.New(), "", []byte(nil), time.Time{},
			},
		},
	}

	for _, test := range tests {
		args := prepareQueryArgs(test.events)
		if len(test.result) != len(args) {
			t.Fatalf("%s failed: length does not match:", test.name)
		}

		counter := 0
		for counter < len(test.result) {
			assert.NotEqual(t, test.result[counter], args[counter], "%s: ids should not be equal", test.name)
			assert.NotEqual(t, test.result[counter+1], args[counter+1], "%s: types should not be equal", test.name)
			counter += 4
		}
	}
}
