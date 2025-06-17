package events

import (
	"github.com/google/uuid"
)

// Envelope wraps any event payload with metadata.
type Envelope[T any] struct {
	EventID   uuid.UUID `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp string    `json:"timestamp"`
	// TODO: Добавить позже
	// Source        string    `json:"source,omitempty"`
	// SchemaVersion string    `json:"schema_version,omitempty"`
	// TraceID       string    `json:"trace_id,omitempty"`
	// RequestID     string    `json:"request_id,omitempty"`
	Payload T `json:"payload"`
}
