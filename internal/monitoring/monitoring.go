// Package monitoring owns availability and performance measurements.
package monitoring

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/vishnunotnice/NetAssist/internal/diagnostics"
)

// Check is a timestamped explicit-target health observation.
type Check struct {
	Timestamp time.Time     `json:"timestamp"`
	Target    string        `json:"target"`
	Reachable bool          `json:"reachable"`
	Duration  time.Duration `json:"duration_ns"`
}

// FileHistory stores append-only local health observations as JSON Lines.
type FileHistory struct{ Path string }

func (h FileHistory) Append(check Check) error {
	if err := os.MkdirAll(filepath.Dir(h.Path), 0o755); err != nil {
		return fmt.Errorf("create monitoring directory: %w", err)
	}
	file, err := os.OpenFile(h.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open monitoring history: %w", err)
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(check)
}

func (h FileHistory) List(limit int) ([]Check, error) {
	file, err := os.Open(h.Path)
	if os.IsNotExist(err) {
		return []Check{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open monitoring history: %w", err)
	}
	defer file.Close()
	checks, err := decodeChecks(file)
	if err != nil {
		return nil, err
	}
	if limit > 0 && len(checks) > limit {
		checks = checks[len(checks)-limit:]
	}
	return checks, nil
}

func decodeChecks(reader io.Reader) ([]Check, error) {
	checks := []Check{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var check Check
		if err := json.Unmarshal(scanner.Bytes(), &check); err != nil {
			return nil, fmt.Errorf("parse monitoring history: %w", err)
		}
		checks = append(checks, check)
	}
	return checks, scanner.Err()
}

// CheckTCP performs and records one explicitly requested connectivity check.
func CheckTCP(ctx context.Context, address string, timeout time.Duration, history FileHistory) (Check, error) {
	result, err := diagnostics.CheckTCP(ctx, address, timeout)
	if err != nil {
		return Check{}, err
	}
	check := Check{Timestamp: time.Now().UTC(), Target: result.Address, Reachable: result.Reachable, Duration: result.Duration}
	if err := history.Append(check); err != nil {
		return Check{}, err
	}
	return check, nil
}
