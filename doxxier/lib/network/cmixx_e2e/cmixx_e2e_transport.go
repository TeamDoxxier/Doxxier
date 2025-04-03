package cmixx_e2e

import (
	"crypto/ed25519"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"doxxier.tech/doxxier/assets"
	"doxxier.tech/doxxier/internal"
	"doxxier.tech/doxxier/pkg/models"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/client/v4/catalog"
	"gitlab.com/elixxir/client/v4/e2e/receive"
	"gitlab.com/elixxir/client/v4/xxdk"
	"gitlab.com/elixxir/crypto/codename"
	"gitlab.com/elixxir/crypto/contact"
	"gitlab.com/elixxir/crypto/nike"
	"gitlab.com/elixxir/ekv/portableOS"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/xx_network/primitives/id"
)

type CMixxE2eConfig struct {
	NdfPath            string
	CertPath           string
	Secret             string
	Recipient          string
	StoragePath        string
	identityStorageKey string
}

const (
	secretFlag             = "cE+a6mP>p`4b]F8;CY&*t}"
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
	recipientConnected bool
}

func NewCMixxE2eTransport(config CMixxE2eConfig) (*CmixxE2eTransport, error) {
	if config.Secret == "" {
		config.Secret = secretFlag
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
	println("Initialized CMixxE2eTransport")
	return transport, nil
}

// Transport interface methods
func (c *CmixxE2eTransport) Disconnect() error {
	c.cMixxNet.StopNetworkFollower()
	SendEventToChannel(&c.callbackChan, []string{"Disconnected from cMixx network"})
	return nil
}

func (c *CmixxE2eTransport) Connect(callbackChan chan string) error {
	c.callbackChan = callbackChan
	err := c.initialiseCMixx()
	identity, err := c.initialiseReceptionIdentity()
	if err != nil {
		return errors.WithMessage(err, "Error initialising reception identity")
	}
	c.identity = identity
	//contactOutput := internal.BytesToBase64(identity.GetContact().Marshal())
	err = c.connectCmixx()
	if err != nil {
		return errors.WithMessage(err, "Error connecting to cMixx")
	}

	return nil
}

func (c *CmixxE2eTransport) Send(contact string, doxxier models.Doxxier) error {
	return nil
}

func (c *CmixxE2eTransport) SendMessage(recipient string, message string) error {
	if !c.recipientConnected {
		err := c.connectRecipient(recipient)
		if err != nil {
			return errors.WithMessage(err, "Error connecting to recipient")
		}
	}
	client := c.user.GetE2E()
	params := xxdk.GetDefaultE2EParams()
	recipientDecoded, err := internal.Base64ToBytes(recipient)
	if err != nil {
		return errors.WithMessage(err, "Error decoding recipient")
	}
	e2eContact, err := contact.Unmarshal(recipientDecoded)
	if err != nil {
		return errors.WithMessage(err, "Error unmarshalling contact")
	}
	report, err := client.SendE2E(catalog.XxMessage, e2eContact.ID, []byte(message), params.Base)
	if err != nil {
		return errors.WithMessage(err, "Error sending message")
	}
	SendEventToChannel(&c.callbackChan, []string{"Message sent: " + string(report.MessageId.StringVerbose())})
	return nil
}

//Helper functions

func SendEventToChannel(channel *chan string, events []string) {
	for _, event := range events {
		*channel <- event
	}
}

func (c *CmixxE2eTransport) connectRecipient(recipient string) error {
	recipientDecoded, err := internal.Base64ToBytes(recipient)
	if err != nil {
		return errors.WithMessage(err, "Error decoding recipient")
	}
	recipientContact, err := contact.Unmarshal(recipientDecoded)
	if err != nil {
		return errors.WithMessage(err, "Error unmarshalling contact")
	}
	e2eClient := c.user.GetE2E()
	time.Sleep(30 * time.Second)

	_, err = e2eClient.GetPartner(recipientContact.ID)
	confirmChan := make(chan contact.Contact, 5)
	if err != nil {
		_, err = c.user.GetAuth().Request(recipientContact, fact.FactList{})
		if err != nil {
			return errors.WithMessage(err, "Error requesting contact")
		}
		timeout := time.NewTimer(30 * time.Second)

		select {
		case pc := <-confirmChan:
			if !pc.ID.Cmp(recipientContact.ID) {
				SendEventToChannel(&c.callbackChan, []string{"Contact confirmation failed"})
			}
			break
		case <-timeout.C:
			SendEventToChannel(&c.callbackChan, []string{"Contact confirmation timed out"})
			break
		}
	}
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
				break
			case <-timeoutTimer.C:
				SendEventToChannel(&c.callbackChan, []string{"Connecting to cMixx timed out"})
				break
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
	e2eClient.RegisterListener(&id.ZeroUser, catalog.NoType, c)
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
	println("Starting initialiseCMixx")
	println(fmt.Printf("NDF Path: %s\n", c.ndfPath))
	println(fmt.Printf("Storage Path: %s\n", c.StoragePath))

	ndf, err := c.getNdf()
	if err != nil {
		println(fmt.Printf("Error getting NDF: %v\n", err))
		return errors.WithMessage(err, "Error getting NDF")
	}
	println("NDF successfully retrieved")

	secretBytes := []byte(c.secret)
	stat, err := portableOS.Stat(c.StoragePath + "/xx")
	if os.IsNotExist(err) || !stat.IsDir() {
		println("Creating new CMixx")
		err = xxdk.NewCmix(string(ndf), c.StoragePath+"/xx", secretBytes, "")
		if err != nil {
			println(fmt.Printf("Error connecting to CMixx: %v\n", err))
			return errors.WithMessage(err, "Error connecting to CMixx")
		}
	}
	println("CMixx initialized successfully")

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
	if c.ndfPath == "" {
		ndf, err := assets.ParseNdf()
		if err != nil {
			return "", errors.WithMessage(err, "Error parsing NDF")
		}
		return ndf, nil
	}
	ndfPath := filepath.Join(internal.GetRootPath(), c.ndfPath)
	jsonNdf, err := os.ReadFile(ndfPath)
	if err != nil {
		return "", errors.WithMessage(err, "Error reading NDF file")
	}
	return string(jsonNdf), nil
}

func (c *CmixxE2eTransport) Hear(item receive.Message) {
	SendEventToChannel(&c.callbackChan, []string{"Received message: " + string(item.Payload)})
}

func (l *CmixxE2eTransport) Name() string {
	return "Transport listener"
}
