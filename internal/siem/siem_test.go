package siem

import (
	"path/filepath"
	"testing"
)

func TestFileStore(t *testing.T) {
	store := FileStore{Path: filepath.Join(t.TempDir(), "events.jsonl")}
	event, err := NewEvent("cli", "audit", "high", "configuration changed")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Append(event); err != nil {
		t.Fatal(err)
	}
	events, err := store.List(20)
	if err != nil || len(events) != 1 || events[0].Severity != "high" {
		t.Fatalf("events = %#v, err = %v", events, err)
	}
}
