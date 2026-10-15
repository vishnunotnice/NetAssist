// Package siem owns security-event intake, normalization, retention, query, and alert boundaries.
package siem

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Event is a normalized local security event. This phase has no external
// collectors, detection rules, automatic response, or background telemetry.
type Event struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
	Kind      string    `json:"kind"`
	Severity  string    `json:"severity"`
	Message   string    `json:"message"`
}

// FileStore is an append-only JSON Lines prototype for local SIEM events.
type FileStore struct{ Path string }

func (s FileStore) Append(event Event) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return fmt.Errorf("create SIEM directory: %w", err)
	}
	file, err := os.OpenFile(s.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open SIEM event store: %w", err)
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(event)
}

func (s FileStore) List(limit int) ([]Event, error) {
	file, err := os.Open(s.Path)
	if os.IsNotExist(err) {
		return []Event{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open SIEM event store: %w", err)
	}
	defer file.Close()
	events, err := decodeEvents(file)
	if err != nil {
		return nil, err
	}
	if limit > 0 && len(events) > limit {
		events = events[len(events)-limit:]
	}
	return events, nil
}

func decodeEvents(reader io.Reader) ([]Event, error) {
	events := []Event{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var event Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, fmt.Errorf("parse SIEM event store: %w", err)
		}
		events = append(events, event)
	}
	return events, scanner.Err()
}

// NewEvent validates and normalizes manual event input.
func NewEvent(source, kind, severity, message string) (Event, error) {
	source, kind, severity, message = strings.TrimSpace(source), strings.TrimSpace(kind), strings.ToLower(strings.TrimSpace(severity)), strings.TrimSpace(message)
	if source == "" || kind == "" || message == "" {
		return Event{}, fmt.Errorf("source, kind, and message are required")
	}
	switch severity {
	case "info", "low", "medium", "high", "critical":
	default:
		return Event{}, fmt.Errorf("severity must be info, low, medium, high, or critical")
	}
	identifier := make([]byte, 8)
	if _, err := rand.Read(identifier); err != nil {
		return Event{}, fmt.Errorf("generate event ID: %w", err)
	}
	return Event{ID: hex.EncodeToString(identifier), Timestamp: time.Now().UTC(), Source: source, Kind: kind, Severity: severity, Message: message}, nil
}
