package network

import (
	"testing"
)

func TestConnect(t *testing.T) {
	transport := &CMixxTransport{
		secret:   "test-secret",
		ndfPath:  "../../test_assets/ndf.json",
		certPath: "../test_assets/main.cert",
	}

	err := transport.Connect()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}
