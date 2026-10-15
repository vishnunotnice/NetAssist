package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func runCommand(t *testing.T, args ...string) string {
	t.Helper()
	cmd := NewRootCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute %v: %v", args, err)
	}
	return buf.String()
}

func TestVersion(t *testing.T) {
	if got := runCommand(t, "version"); !strings.Contains(got, "netassist") {
		t.Fatalf("version output = %q", got)
	}
}

func TestConfigInitAndShow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	runCommand(t, "config", "init", "--config", path)
	if got := runCommand(t, "config", "show", "--config", path); !strings.Contains(got, "enabled: false") {
		t.Fatalf("config output = %q", got)
	}
}

func TestDiscoverUsesDefaultHostLimit(t *testing.T) {
	got := runCommand(t, "discover", "192.168.10.0/30")
	if !strings.Contains(got, "Planned 2 host(s)") {
		t.Fatalf("discovery output = %q", got)
	}
}

func TestDeviceInventory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "devices.json")
	runCommand(t, "devices", "add", "router", "--address", "192.168.1.1", "--tags", "core", "--store", path)
	if got := runCommand(t, "devices", "list", "--store", path); !strings.Contains(got, "router\t192.168.1.1\tcore") {
		t.Fatalf("device output = %q", got)
	}
}

func TestSIEMEventStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "events.jsonl")
	runCommand(t, "siem", "ingest", "--source", "cli", "--kind", "audit", "--message", "test event", "--store", path)
	if got := runCommand(t, "siem", "events", "--store", path); !strings.Contains(got, "test event") {
		t.Fatalf("SIEM output = %q", got)
	}
}
