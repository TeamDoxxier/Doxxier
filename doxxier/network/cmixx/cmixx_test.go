package cmixx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	transport := &CMixxTransport{
		secret:      "test-secret",
		ndfPath:     "../../test_assets/ndf.json",
		certPath:    "../../test_assets/main.cert",
		StoragePath: "../../test_storage",
	}

	err := transport.Connect()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestInitaliseCMixx(t *testing.T) {
	transport := &CMixxTransport{
		secret:      "test-secret",
		ndfPath:     "../../test_assets/ndf.json",
		certPath:    "../../test_assets/main.cert",
		StoragePath: "../../test_storage",
	}

	err := transport.initialiseCMixx()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}
func TestGetNdf(t *testing.T) {
	transport := &CMixxTransport{
		ndfPath:  "../../test_assets/ndf.json",
		certPath: "../../test_assets/main.cert",
	}

	ndf, err := transport.getNdf()
	assert.NoError(t, err)
	assert.NotNil(t, ndf)
}
