package util

import (
	"net"
	"testing"
)

func TestExtractor(t *testing.T) {
	testData := []struct {
		addr   string
		expect string
		parse  bool
	}{
		{"127.0.0.1", "127.0.0.1", false},
		{"10.0.0.1", "10.0.0.1", false},
		{"", "", true},
		{"0.0.0.0", "", true},
		{"[::]", "", true},
	}

	for _, d := range testData {
		addr, err := Extract(d.addr)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if d.parse {
			ip := net.ParseIP(addr)
			if ip == nil {
				t.Error("Unexpected nil IP")
			}

		} else if addr != d.expect {
			t.Errorf("Expected %s got %s", d.expect, addr)
		}
	}

}

func TestExtract(t *testing.T) {
	addr, err := Extract("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addr)
}

func TestIPs(t *testing.T) {
	t.Log(IPs())
}

func TestAddress(t *testing.T) {
	addr, err := Address("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addr)
}
