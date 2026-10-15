// Package monitoring owns availability and performance measurement.
package monitoring

// Service is the monitoring application boundary.
type Service interface {
	Collect() error
}
