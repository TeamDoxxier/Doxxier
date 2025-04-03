package cmix_e2e

import (
	"fmt"
	"os"

	"doxxier.tech/doxxier/assets"
	"doxxier.tech/doxxier/internal"
	"doxxier.tech/doxxier/pkg/models"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/client/v4/bindings"
	"gitlab.com/elixxir/client/v4/xxdk"
	"gitlab.com/elixxir/ekv/portableOS"
)

type CMixE2eConfig struct {
	NdfPath            string
	CertPath           string
	Secret             string
	Recipient          string
	StoragePath        string
	identityStorageKey string
	params             []byte
}

type CmixE2eTransport struct {
	ndfPath          string
	storagePath      string
	secret           string
	params           []byte
	cmix             *bindings.Cmix
	connection       *bindings.Connection
	e2e              *bindings.E2e
	identity         []byte
	callbackChan     chan string
	healthCallbackId int64
}

const (
	secretFlag             = "cE+a6mP>p`4b]F8;CY&*t}"
	certPathFlag           = "assets/main.cert"
	identityStorageKeyFlag = "identityStore"
)

func NewCmixE2eTransport(config CMixE2eConfig) (*CmixE2eTransport, error) {
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
	return &CmixE2eTransport{
		ndfPath:     config.NdfPath,
		storagePath: config.StoragePath,
		secret:      config.Secret,
		params:      []byte(""),
	}, nil
}

// Transport interface methods
func (c *CmixE2eTransport) Initialise() (string, error) {
	err := c.initialiseCMix()
	if err != nil {
		return "", errors.WithMessage(err, "Error initializing CMix")
	}
	identity, err := c.initialiseIdentity()
	if err != nil {
		return "", errors.WithMessage(err, "Error initializing identity")
	}
	c.identity = identity
	decodedIdentity, err := xxdk.UnmarshalReceptionIdentity(identity)
	if err != nil {
		return "", errors.WithMessage(err, "Error decoding identity")

	}
	contact := internal.BytesToBase64(decodedIdentity.GetContact().Marshal())

	return contact, nil
}

func (c *CmixE2eTransport) Connect(callbackChan chan string, recipient string) error {
	if callbackChan == nil {
		return errors.New("callbackChan is required")
	}
	if recipient == "" {
		return errors.New("recipient is required")
	}
	c.callbackChan = callbackChan
	err := c.connectCmixx(recipient)
	if err != nil {
		return errors.WithMessage(err, "Error connecting to CMix")
	}

	return nil
}

func (c *CmixE2eTransport) Disconnect() error {
	return nil
}

func (c *CmixE2eTransport) Send(contact string, doxxier models.Doxxier) error {
	return nil
}

func (c *CmixE2eTransport) SendMessage(message string) error {
	report, err := c.connection.SendE2E(2, []byte(message))
	println(fmt.Printf("Send report: %s\n", report))
	return err
}

// Internal helper methods
func (c *CmixE2eTransport) connectCmixx(recipient string) error {
	e2e, err := bindings.Login(c.cmix.GetID(), nil, c.identity, c.params)
	if err != nil {
		return errors.WithMessage(err, "Error logging in")
	}
	c.e2e = e2e
	recipientContact, err := internal.Base64ToBytes(recipient)
	if err != nil {
		return errors.WithMessage(err, "Error converting recipient to bytes")
	}
	c.cmix.StartNetworkFollower(10000)
	c.cmix.WaitForNetwork(30000)
	// waitUntilConnected := func(connected chan bool) {
	// 	waitTimeout := 30
	// 	timeoutTimer := time.NewTimer(time.Duration(waitTimeout) * time.Second)
	// 	isConnected := false
	// 	for !isConnected {
	// 		select {
	// 		case isConnected = <-connected:
	// 			sendEventToChannel(&c.callbackChan, []string{"Connected to cMixx network"})
	// 			break
	// 		case <-timeoutTimer.C:
	// 			sendEventToChannel(&c.callbackChan, []string{"Connecting to cMixx timed out"})
	// 			break
	// 		}
	// 	}
	// }
	// connected := make(chan bool, 10)

	// Use the struct to pass the callback
	c.healthCallbackId = c.cmix.AddHealthCallback(c)
	// waitUntilConnected(connected)

	connection, err := c.cmix.Connect(e2e.GetID(), recipientContact, c.params)
	if err != nil {
		return errors.WithMessage(err, "Error connecting to CMix")
	}
	c.e2e.RegisterListener([]byte{}, 0, c)
	c.connection = connection
	return nil
}

func (c *CmixE2eTransport) initialiseIdentity() ([]byte, error) {
	identity, err := c.cmix.MakeReceptionIdentity()
	if err != nil {
		return nil, errors.WithMessage(err, "Error making reception identity")
	}
	return identity, nil
}

func (c *CmixE2eTransport) initialiseCMix() error {
	println("Starting initialiseCMix")
	println(fmt.Printf("NDF Path: %s\n", c.ndfPath))
	println(fmt.Printf("Storage Path: %s\n", c.storagePath))

	ndf, err := c.getNdf()
	if err != nil {
		println(fmt.Printf("Error getting NDF: %v\n", err))
		return errors.WithMessage(err, "Error getting NDF")
	}
	println("NDF successfully retrieved")

	secretBytes := []byte(c.secret)
	stat, err := portableOS.Stat(c.storagePath + "/xx")
	if os.IsNotExist(err) || !stat.IsDir() {
		println("Creating new CMix")
		err = bindings.NewCmix(string(ndf), c.storagePath+"/xx", secretBytes, "")
		if err != nil {
			println(fmt.Printf("Error connecting to CMix: %v\n", err))
			return errors.WithMessage(err, "Error connecting to CMix")
		}
	}
	println("CMix initialized successfully")
	secretBytes = []byte(c.secret)
	cmixx, err := bindings.LoadCmix(c.storagePath+"/xx", secretBytes, c.params)
	if err != nil {
		return errors.WithMessage(err, "Error connecting to CMix")
	}
	c.cmix = cmixx
	return nil
}

func (c *CmixE2eTransport) getNdf() (string, error) {
	//Read NDF file
	if c.ndfPath == "" {
		ndf, err := assets.ParseNdf()
		if err != nil {
			return "", errors.WithMessage(err, "Error parsing NDF")
		}
		return ndf, nil
	}
	return "", errors.New("Unable to read NDF file")
}

//Helper functions

func sendEventToChannel(channel *chan string, events []string) {
	for _, event := range events {
		*channel <- event
	}
}

func (c *CmixE2eTransport) Callback(connected bool) {
	println(fmt.Printf("Connected: %v\n", connected))
	if connected {
		sendEventToChannel(&c.callbackChan, []string{"Connected to CMix"})
	} else {
		sendEventToChannel(&c.callbackChan, []string{"Disconnected from CMix"})
	}
}

func (c *CmixE2eTransport) Hear(item []byte) {
	sendEventToChannel(&c.callbackChan, []string{"Received message: " + string(item)})
}

func (l *CmixE2eTransport) Name() string {
	return "Transport listener"
}
