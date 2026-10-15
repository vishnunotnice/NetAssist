// Package agent owns cross-platform device enrollment and telemetry collection.
package agent

// Service is the future agent control-plane contract.
type Service interface {
	Register() error
}
