// Package diagnostics owns safe, user-authorized network troubleshooting operations.
package diagnostics

// Service is the diagnostics application boundary.
type Service interface {
	Check() error
}
