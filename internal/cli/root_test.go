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
