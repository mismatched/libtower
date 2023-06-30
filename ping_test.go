package libtower

import (
	"testing"
)

func TestPing(t *testing.T) {
	_, _, err := Ping("google.com", 1)
	if err != nil {
		t.Errorf("test failed %v", err)
	}
}

func TestPingIPv6(t *testing.T) {
	_, _, err := Ping("ipv6.google.com", 1)
	if err != nil {
		t.Errorf("test failed %v", err)
	}
}
