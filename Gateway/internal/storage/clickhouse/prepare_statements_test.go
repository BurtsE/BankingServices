package clickhouse

import (
	"gateway/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPrepareQueryStatements(t *testing.T) {
	tests := []struct {
		name   string
		events []domain.Event
		result string
	}{
		{
			name:   "empty slice",
			events: []domain.Event{},
			result: "INSERT INTO events (uuid, value) VALUES",
		},
		{
			name: "single event",
			events: []domain.Event{
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: "INSERT INTO events (uuid, value) VALUES ($1, $2)",
		},
		{
			name: "multiple events",
			events: []domain.Event{
				{
					EventType: domain.CreateUserEvent,
				},
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: "INSERT INTO events (uuid, value) VALUES ($1, $2),($3, $4)",
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
		events []domain.Event
		result []any
	}{
		{
			name:   "empty slice",
			events: []domain.Event{},
			result: []any{},
		},
		{
			name: "single event",
			events: []domain.Event{
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: []any{
				any(uuid.UUID{}), domain.CreateUserEvent.Value(),
			},
		},
		{
			name: "multiple events",
			events: []domain.Event{
				{
					EventType: domain.CreateUserEvent,
				},
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: []any{
				any(uuid.UUID{}), domain.CreateUserEvent.Value(),
				any(uuid.UUID{}), domain.CreateUserEvent.Value(),
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
			counter += 2
		}
	}
}

func TestPrepareQueryArgsFail(t *testing.T) {
	tests := []struct {
		name   string
		events []domain.Event
		result []any
	}{
		{
			name: "single event",
			events: []domain.Event{
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: []any{
				uuid.New(), "delete_user",
			},
		},
		{
			name: "multiple events",
			events: []domain.Event{
				{
					EventType: domain.CreateUserEvent,
				},
				{
					EventType: domain.CreateUserEvent,
				},
			},
			result: []any{
				uuid.New(), "",
				uuid.New(), "",
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
			counter += 2
		}
	}
}
