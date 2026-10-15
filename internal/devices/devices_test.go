package devices

import (
	"path/filepath"
	"testing"
)

func TestFileStore(t *testing.T) {
	store := FileStore{Path: filepath.Join(t.TempDir(), "devices.json")}
	device, err := New("router", "192.168.1.1", []string{"network", " core "})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Add(device); err != nil {
		t.Fatal(err)
	}
	items, err := store.List()
	if err != nil || len(items) != 1 || items[0].Tags[1] != "core" {
		t.Fatalf("items = %#v, err = %v", items, err)
	}
}
