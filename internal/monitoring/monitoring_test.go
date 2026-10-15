package monitoring

import (
	"path/filepath"
	"testing"
	"time"
)

func TestHistory(t *testing.T) {
	history := FileHistory{Path: filepath.Join(t.TempDir(), "history.jsonl")}
	if err := history.Append(Check{Timestamp: time.Now(), Target: "localhost:80", Reachable: true}); err != nil {
		t.Fatal(err)
	}
	checks, err := history.List(20)
	if err != nil || len(checks) != 1 || !checks[0].Reachable {
		t.Fatalf("checks = %#v, err = %v", checks, err)
	}
}
