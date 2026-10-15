// Package api owns the HTTP/gRPC control-plane boundary.
package api

// Service is implemented by API transports; domain packages remain transport-neutral.
type Service interface {
	Start() error
}
