package discovery

import "testing"

func TestPlanListsUsableAddresses(t *testing.T) {
	addresses, err := Plan("192.168.10.0/30", 2)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(addresses), 2; got != want || addresses[0].String() != "192.168.10.1" || addresses[1].String() != "192.168.10.2" {
		t.Fatalf("addresses = %v", addresses)
	}
}

func TestPlanRejectsPublicNetwork(t *testing.T) {
	if _, err := Plan("198.51.100.0/30", 2); err == nil {
		t.Fatal("expected public CIDR rejection")
	}
}
