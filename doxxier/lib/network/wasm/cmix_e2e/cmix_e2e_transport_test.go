package cmix_e2e

import (
	"testing"

	"doxxier.tech/doxxier/pkg/models"
	"github.com/stretchr/testify/assert"
)

func getConfig() CMixE2eConfig {
	return CMixE2eConfig{
		NdfPath:     "ndf/path",
		StoragePath: "/tmp/test-cmix",
	}
}

func TestNewCmixE2eTransport_MissingStoragePath(t *testing.T) {
	config := getConfig()
	config.StoragePath = ""
	transport, err := NewCmixE2eTransport(config)
	assert.Nil(t, transport)
	assert.Error(t, err)
	assert.Equal(t, "StoragePath is required", err.Error())
}

func TestNewCmixE2eTransport_Defaults(t *testing.T) {
	config := CMixE2eConfig{
		NdfPath:     "ndf/path",
		StoragePath: "/tmp/test-cmix",
	}

	transport, err := NewCmixE2eTransport(config)
	assert.NoError(t, err)
	assert.Equal(t, config.NdfPath, transport.ndfPath)
	assert.Equal(t, config.StoragePath, transport.storagePath)
	assert.Equal(t, secretFlag, transport.secret)
	assert.NotNil(t, transport)
}

func TestInitialise_DefaultNdf(t *testing.T) {
	config := getConfig()
	config.NdfPath = ""
	transport, _ := NewCmixE2eTransport(config)
	contact, err := transport.Initialise()
	assert.NoError(t, err)
	assert.NotEmpty(t, contact)
}

func TestInitialise_MissingNdf(t *testing.T) {
	config := getConfig()
	config.NdfPath = "missing/path"
	transport, err := NewCmixE2eTransport(getConfig())
	_, err = transport.Initialise()
	assert.Error(t, err)
}

func TestGetNdf_DefaultPath(t *testing.T) {
	config := getConfig()
	config.NdfPath = ""
	transport, _ := NewCmixE2eTransport(config)
	ndf, err := transport.getNdf()
	assert.Nil(t, err)
	assert.NotNil(t, ndf)
}

func TestDisconnect(t *testing.T) {
	transport := &CmixE2eTransport{}
	err := transport.Disconnect()
	assert.NoError(t, err)
}

func TestSend(t *testing.T) {
	transport := &CmixE2eTransport{}
	err := transport.Send("someContact", models.Doxxier{})
	assert.NoError(t, err)
}

func TestSendMessage(t *testing.T) {
	transport := &CmixE2eTransport{}
	err := transport.SendMessage("recipient", "message")
	assert.NoError(t, err)
}
