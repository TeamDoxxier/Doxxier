package cmixx_e2e

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/elixxir/client/v4/xxdk"
)

func TestNewCMixxE2eTransport(t *testing.T) {
	config := CMixxE2eConfig{
		StoragePath: "./test_storage",
	}

	transport, err := NewCMixxE2eTransport(config)
	assert.Nil(t, err)
	assert.NotNil(t, transport)
	assert.Equal(t, secretFlag, transport.secret)
	assert.Equal(t, certPathFlag, transport.certPath)
	assert.Equal(t, identityStorageKeyFlag, transport.identityStorageKey)
}

func TestCMixxE2eTransport_Disconnect(t *testing.T) {
	transport := &CmixxE2eTransport{
		cMixxNet:     &xxdk.Cmix{},
		callbackChan: make(chan string, 1),
	}

	err := transport.Disconnect()
	assert.Nil(t, err)
	select {
	case msg := <-transport.callbackChan:
		assert.Equal(t, "Disconnected from cMixx network", msg)
	case <-time.After(time.Second):
		t.Fatal("expected message not received")
	}
}

// func TestCMixxE2eTransport_Connect(t *testing.T) {
// 	transport := &CmixxE2eTransport{
// 		ndfPath:     "./test_ndf",
// 		StoragePath: "./test_storage",
// 	}
// 	callbackChan := make(chan string, 1)
// 	assert.NotNil(t, transport)
// 	assert.NotNil(t, callbackChan)
// 	// Mocking external dependencies
// 	// initialiseCMixx := func(transport *CmixxE2eTransport) error {
// 	// 	return nil
// 	// }
// 	// initialiseReceptionIdentity = func(transport *CmixxE2eTransport) (*xxdk.ReceptionIdentity, error) {
// 	// 	return &xxdk.ReceptionIdentity{}, nil
// 	// }
// 	// connectCmixx = func(transport *CmixxE2eTransport) error {
// 	// 	return nil
// 	//}

// 	// contactOutput, err := transport.Connect(callbackChan)
// 	// assert.Nil(t, err)
// 	// assert.NotEmpty(t, contactOutput)
// }

func TestCMixxE2eTransport_SendMessage_NotConnected(t *testing.T) {
	transport := &CmixxE2eTransport{
		user:               xxdk.E2e{},
		recipientConnected: false,
		callbackChan:       make(chan string, 1),
	}

	recipient := "testRecipient"
	message := "Hello, world!"

	// Mock connectRecipient to simulate successful connection
	// connectRecipient := func(transport *CmixxE2eTransport, recipient string) error {
	// 	return nil
	// }

	err := transport.SendMessage(recipient, message)
	assert.Nil(t, err)
	select {
	case msg := <-transport.callbackChan:
		assert.Contains(t, msg, "Message sent")
	case <-time.After(time.Second):
		t.Fatal("expected message not received")
	}
}

func TestCMixxE2eTransport_getNdf(t *testing.T) {
	config := CMixxE2eConfig{
		StoragePath: "./test_storage",
	}
	transport, err := NewCMixxE2eTransport(config)
	// 	ndfPath: "./test_ndf.json",
	// }

	// // Create mock NDF file
	// err := os.WriteFile(transport.ndfPath, []byte("mock NDF content"), 0644)
	// assert.Nil(t, err)
	// defer os.Remove(transport.ndfPath)

	ndf, err := transport.getNdf()
	assert.Nil(t, err)
	assert.NotNil(t, ndf)
}

func TestCMixxE2eTransport_initialiseReceptionIdentity(t *testing.T) {
	transport := &CmixxE2eTransport{
		identityStorageKey: "test_identity_key",
		cMixxNet:           &xxdk.Cmix{},
	}

	// Mock LoadReceptionIdentity
	// xxdk.LoadReceptionIdentity = func(key string, net *xxdk.Cmix) (xxdk.ReceptionIdentity, error) {
	// 	return xxdk.ReceptionIdentity{}, errors.New("not found")
	// }

	// xxdk.MakeReceptionIdentity = func(net *xxdk.Cmix) (xxdk.ReceptionIdentity, error) {
	// 	return xxdk.ReceptionIdentity{}, nil
	// }

	identity, err := transport.initialiseReceptionIdentity()
	assert.Nil(t, err)
	assert.NotNil(t, identity)
}
