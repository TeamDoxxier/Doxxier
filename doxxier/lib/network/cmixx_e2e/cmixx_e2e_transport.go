package network

import (
	"crypto/ed25519"
	"os"
	"path/filepath"
	"time"

	"doxxier.tech/doxxier/internal"
	"doxxier.tech/doxxier/pkg/models"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/client/v4/catalog"
	"gitlab.com/elixxir/client/v4/xxdk"
	"gitlab.com/elixxir/crypto/codename"
	"gitlab.com/elixxir/crypto/contact"
	"gitlab.com/elixxir/crypto/nike"
	"gitlab.com/xx_network/primitives/id"
)

type CMixxConfig struct {
	NdfPath            string
	CertPath           string
	Secret             string
	Recipient          string
	StoragePath        string
	identityStorageKey string
}

const (
	secretFlag             = "cE+a6mP>p`4b]F8;CY&*t}"
	ndfPathFlag            = "assets/ndf.json"
	certPathFlag           = "assets/main.cert"
	identityStorageKeyFlag = "identityStore"
)

type CmixxE2eTransport struct {
	secret             string
	ndfPath            string
	certPath           string
	StoragePath        string
	identityStorageKey string
	cMixxNet           *xxdk.Cmix
	directId           codename.PrivateIdentity
	user               xxdk.E2e
	privateKey         nike.PrivateKey
	publicKey          nike.PublicKey
	receptionId        *id.ID
	selfReceptionId    *id.ID
	partnerPublicKey   ed25519.PublicKey
	partnerToken       uint32
	identity           *xxdk.ReceptionIdentity
	callbackChan       chan string
}

func NewCMixxE2eTransport(config CMixxConfig) (*CmixxE2eTransport, error) {
	if config.Secret == "" {
		config.Secret = secretFlag
	}
	if config.NdfPath == "" {
		config.NdfPath = ndfPathFlag
	}
	if config.CertPath == "" {
		config.CertPath = certPathFlag
	}
	if config.StoragePath == "" {
		return nil, errors.New("StoragePath is required")
	}
	if config.identityStorageKey == "" {
		config.identityStorageKey = identityStorageKeyFlag
	}
	transport := &CmixxE2eTransport{
		secret:             config.Secret,
		ndfPath:            config.NdfPath,
		certPath:           config.CertPath,
		StoragePath:        config.StoragePath,
		identityStorageKey: config.identityStorageKey,
	}
	return transport, nil
}

func (c *CmixxE2eTransport) Disconnect() error {
	c.cMixxNet.StopNetworkFollower()
	SendEventToChannel(&c.callbackChan, []string{"Disconnected from cMixx network"})
	return nil
}

func SendEventToChannel(channel *chan string, events []string) {
	for _, event := range events {
		*channel <- event
	}
}

func (c *CmixxE2eTransport) Connect(callbackChan chan string) (string, error) {
	c.callbackChan = callbackChan
	err := c.initialiseCMixx()
	if err != nil {
		return "", errors.WithMessage(err, "Error initialising CMixx")
	}

	identity, err := c.initialiseReceptionIdentity()
	if err != nil {
		return "", errors.WithMessage(err, "Error initialising reception identity")
	}
	c.identity = identity
	contactOutput := internal.BytesToBase64(identity.GetContact().Marshal())
	err = c.connectCmixx()
	if err != nil {
		return "", errors.WithMessage(err, "Error connecting to cMixx")
	}
	return contactOutput, nil
}

func (c *CmixxE2eTransport) Send(contact string, doxxier models.Doxxier) error {
	return nil
}

func (c *CmixxE2eTransport) connectCmixx() error {
	params := xxdk.GetDefaultE2EParams()
	confirmChan := make(chan contact.Contact, 5)
	user, err := xxdk.Login(c.cMixxNet, &auth{confirmChan: confirmChan}, *c.identity, params)
	if err != nil {
		return errors.WithMessage(err, "Error logging in to cMixx")
	}
	c.user = *user
	e2eClient := c.user.GetE2E()
	err = user.StartNetworkFollower(5)
	if err != nil {
		return errors.WithMessage(err, "Error starting network follower")
	}

	waitUntilConnected := func(connected chan bool) {
		waitTimeout := 30
		timeoutTimer := time.NewTimer(time.Duration(waitTimeout) * time.Second)
		isConnected := false
		for !isConnected {
			select {
			case isConnected = <-connected:
				SendEventToChannel(&c.callbackChan, []string{"Connected to cMixx network"})
			case <-timeoutTimer.C:
				SendEventToChannel(&c.callbackChan, []string{"Connecting to cMixx timed out"})
			}
		}
	}
	connected := make(chan bool, 10)
	user.GetCmix().AddHealthCallback(
		func(isConnected bool) {
			connected <- isConnected
		},
	)
	waitUntilConnected(connected)
	e2eClient.RegisterListener(&id.ZeroUser, catalog.NoType, listener{name: "listener"})
	return nil
}

func (c *CmixxE2eTransport) initialiseReceptionIdentity() (*xxdk.ReceptionIdentity, error) {
	identity, err := xxdk.LoadReceptionIdentity(c.identityStorageKey, c.cMixxNet)
	if err != nil {
		identity, err = xxdk.MakeReceptionIdentity(c.cMixxNet)
		if err != nil {
			return nil, errors.WithMessage(err, "Error making reception identity")
		}
		err = xxdk.StoreReceptionIdentity(c.identityStorageKey, identity, c.cMixxNet)
		if err != nil {
			return nil, errors.WithMessage(err, "Error storing reception identity")
		}
	}
	return &identity, nil
}

func (c *CmixxE2eTransport) initialiseCMixx() error {
	//Establish session
	ndf, err := c.getNdf()
	if err != nil {
		return errors.WithMessage(err, "Error getting NDF")
	}
	//Connect to CMixx
	secretBytes := []byte(c.secret)
	stat, err := os.Stat(c.StoragePath + "/xx")
	if os.IsNotExist(err) || !stat.IsDir() {
		err = xxdk.NewCmix(string(ndf), c.StoragePath+"/xx", secretBytes, "")
		if err != nil {
			return errors.WithMessage(err, "Error connecting to CMixx")
		}
	}

	params := xxdk.GetDefaultCMixParams()
	net, err := xxdk.LoadCmix(c.StoragePath+"/xx", secretBytes, params)
	if err != nil {
		return errors.WithMessage(err, "Error connecting to CMixx")
	}
	c.cMixxNet = net
	return nil
}

func (c *CmixxE2eTransport) getNdf() (string, error) {
	//Read NDF file
	ndfPath := filepath.Join(internal.GetRootPath(), c.ndfPath)
	jsonNdf, err := os.ReadFile(ndfPath)
	if err != nil {
		return "", errors.WithMessage(err, "Error reading NDF file")
	}
	return string(jsonNdf), nil
}
