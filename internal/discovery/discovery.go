// Package discovery owns authorized asset discovery and inventory enrichment.
package discovery

// Service is the discovery application boundary.
type Service interface {
	Discover() error
}
