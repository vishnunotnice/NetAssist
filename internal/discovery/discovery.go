// Package discovery owns authorized asset discovery and inventory enrichment.
package discovery

import (
	"fmt"
	"net/netip"
)

// Plan validates a local IPv4 subnet and lists its usable host addresses.
// It intentionally limits discovery to a caller-selected private, loopback, or
// link-local network; public address ranges are outside Phase 2 scope.
func Plan(cidr string, maxHosts int) ([]netip.Addr, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return nil, fmt.Errorf("parse CIDR: %w", err)
	}
	if !prefix.Addr().Is4() {
		return nil, fmt.Errorf("only IPv4 discovery is supported in Phase 2")
	}
	if !isLocal(prefix.Addr()) {
		return nil, fmt.Errorf("CIDR must be private, loopback, or link-local")
	}
	if maxHosts < 1 {
		return nil, fmt.Errorf("max hosts must be positive")
	}

	addresses := make([]netip.Addr, 0)
	for address := prefix.Masked().Addr(); prefix.Contains(address); address = address.Next() {
		if isNetworkOrBroadcast(prefix, address) {
			continue
		}
		addresses = append(addresses, address)
		if len(addresses) > maxHosts {
			return nil, fmt.Errorf("CIDR exceeds the %d-host safety limit", maxHosts)
		}
	}
	return addresses, nil
}

func isLocal(address netip.Addr) bool {
	return address.IsPrivate() || address.IsLoopback() || address.IsLinkLocalUnicast()
}

func isNetworkOrBroadcast(prefix netip.Prefix, address netip.Addr) bool {
	if prefix.Bits() >= 31 {
		return false
	}
	return address == prefix.Masked().Addr() || !prefix.Contains(address.Next())
}
