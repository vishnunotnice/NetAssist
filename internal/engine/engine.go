// Package engine reserves the boundary for a future independently designed networking engine.
// Existing connectivity technologies remain external integrations until a replacement is designed,
// audited, and explicitly approved.
package engine

// Service is the networking-engine application boundary.
type Service interface {
	Connect() error
}
