// Package diagnostics owns safe, user-authorized network troubleshooting operations.
package diagnostics

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"time"
)

// Resolver allows DNS resolution to be tested without external DNS traffic.
type Resolver interface {
	LookupNetIP(context.Context, string, string) ([]netip.Addr, error)
}

// DNSResult is the resolved address set for a hostname.
type DNSResult struct {
	Hostname  string
	Addresses []netip.Addr
}

// LookupDNS resolves a hostname using the supplied resolver.
func LookupDNS(ctx context.Context, resolver Resolver, hostname string) (DNSResult, error) {
	hostname = strings.TrimSpace(hostname)
	if hostname == "" {
		return DNSResult{}, fmt.Errorf("hostname is required")
	}
	addresses, err := resolver.LookupNetIP(ctx, "ip", hostname)
	if err != nil {
		return DNSResult{}, fmt.Errorf("resolve %q: %w", hostname, err)
	}
	return DNSResult{Hostname: hostname, Addresses: addresses}, nil
}

// TCPResult records an explicit, single-target TCP connectivity check.
type TCPResult struct {
	Address   string
	Reachable bool
	Duration  time.Duration
}

// CheckTCP connects to one explicitly supplied TCP address. It does not perform
// port discovery or attempt fallback ports.
func CheckTCP(ctx context.Context, address string, timeout time.Duration) (TCPResult, error) {
	if _, _, err := net.SplitHostPort(address); err != nil {
		return TCPResult{}, fmt.Errorf("address must be host:port: %w", err)
	}
	if timeout <= 0 {
		return TCPResult{}, fmt.Errorf("timeout must be positive")
	}
	started := time.Now()
	connection, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
	duration := time.Since(started)
	if err != nil {
		return TCPResult{Address: address, Duration: duration}, nil
	}
	if closeErr := connection.Close(); closeErr != nil {
		return TCPResult{}, fmt.Errorf("close connection: %w", closeErr)
	}
	return TCPResult{Address: address, Reachable: true, Duration: duration}, nil
}
