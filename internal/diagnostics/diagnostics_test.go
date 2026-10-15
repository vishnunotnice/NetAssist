package diagnostics

import (
	"context"
	"net/netip"
	"testing"
)

type fakeResolver struct{}

func (fakeResolver) LookupNetIP(context.Context, string, string) ([]netip.Addr, error) {
	return []netip.Addr{netip.MustParseAddr("192.0.2.10")}, nil
}

func TestLookupDNS(t *testing.T) {
	result, err := LookupDNS(context.Background(), fakeResolver{}, "example.test")
	if err != nil || len(result.Addresses) != 1 || result.Addresses[0].String() != "192.0.2.10" {
		t.Fatalf("result = %#v, err = %v", result, err)
	}
}

func TestCheckTCPRejectsMissingPort(t *testing.T) {
	if _, err := CheckTCP(context.Background(), "example.test", 1); err == nil {
		t.Fatal("expected invalid address error")
	}
}
