// Package siem owns security-event intake, normalization, retention, query, and alert boundaries.
// It performs no collection, detection, or response actions in Phase 1.
package siem

import "time"

// Event is the transport-neutral normalized record for future SIEM ingestion.
type Event struct {
	Timestamp time.Time
	Source    string
	Kind      string
	Message   string
}

// Store defines the future persistence and query seam for security events.
type Store interface {
	Append(Event) error
}
