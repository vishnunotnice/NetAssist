// Package devices owns the locally managed device inventory.
package devices

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Device is a deliberately small inventory record. Credentials and connection
// secrets are never stored here.
type Device struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Store persists the local inventory.
type Store interface {
	Add(Device) error
	List() ([]Device, error)
}

// FileStore stores the inventory as a single JSON document suitable for a
// local prototype. A database-backed implementation can replace it later.
type FileStore struct{ Path string }

func (s FileStore) Add(device Device) error {
	devices, err := s.List()
	if err != nil {
		return err
	}
	for _, existing := range devices {
		if existing.Name == device.Name {
			return fmt.Errorf("device named %q already exists", device.Name)
		}
	}
	devices = append(devices, device)
	return s.write(devices)
}

func (s FileStore) List() ([]Device, error) {
	data, err := os.ReadFile(s.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return []Device{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read device inventory: %w", err)
	}
	var devices []Device
	if err := json.Unmarshal(data, &devices); err != nil {
		return nil, fmt.Errorf("parse device inventory: %w", err)
	}
	sort.Slice(devices, func(i, j int) bool { return devices[i].Name < devices[j].Name })
	return devices, nil
}

func (s FileStore) write(devices []Device) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return fmt.Errorf("create device inventory directory: %w", err)
	}
	data, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		return fmt.Errorf("encode device inventory: %w", err)
	}
	if err := os.WriteFile(s.Path, data, 0o600); err != nil {
		return fmt.Errorf("write device inventory: %w", err)
	}
	return nil
}

// New validates input and returns a new inventory record.
func New(name, address string, tags []string) (Device, error) {
	name, address = strings.TrimSpace(name), strings.TrimSpace(address)
	if name == "" || address == "" {
		return Device{}, fmt.Errorf("device name and address are required")
	}
	identifier := make([]byte, 8)
	if _, err := rand.Read(identifier); err != nil {
		return Device{}, fmt.Errorf("generate device ID: %w", err)
	}
	return Device{ID: hex.EncodeToString(identifier), Name: name, Address: address, Tags: cleanTags(tags), CreatedAt: time.Now().UTC()}, nil
}

func cleanTags(tags []string) []string {
	clean := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag = strings.TrimSpace(tag); tag != "" {
			clean = append(clean, tag)
		}
	}
	return clean
}
