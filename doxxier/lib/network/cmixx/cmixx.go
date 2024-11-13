package cmixx

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/elixxir/client/v4/collective/versioned"
	"gitlab.com/elixxir/client/v4/xxdk"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/codename"
	"gitlab.com/elixxir/crypto/nike"
	"gitlab.com/elixxir/crypto/nike/ecdh"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/crypto/tls"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/protobuf/proto"
)

type CMixxTransport struct {
	secret           string
	ndfPath          string
	certPath         string
	StoragePath      string
	cMixxNet         *xxdk.Cmix
	directId         codename.PrivateIdentity
	user             xxdk.E2e
	privateKey       nike.PrivateKey
	publicKey        nike.PublicKey
	receptionId      *id.ID
	selfReceptionId  *id.ID
	partnerPublicKey ed25519.PublicKey
	partnerToken     uint32
}

type CMixxConfig struct {
	NdfPath     string
	CertPath    string
	Secret      string
	Recipient   string
	StoragePath string
}

const (
	secretFlag   = "cE+a6mP>p`4b]F8;CY&*t}"
	ndfPathFlag  = "../assets/ndf.json"
	certPathFlag = "../assets/main.cert"
)

func NewCMixxTransport(config CMixxConfig) *CMixxTransport {
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
		panic("Storage path is required")
	}

	transport := &CMixxTransport{
		secret:      config.Secret,
		ndfPath:     config.NdfPath,
		certPath:    config.CertPath,
		StoragePath: config.StoragePath,
	}
	return transport
}

func (c *CMixxTransport) Connect(partnerPubKey string, partnerToken string) error {
	err := c.initialiseCMixx()
	if err != nil {
		return errors.WithMessage(err, "Error initialising CMixx")
	}
	err = c.initialiseReceptionIdentity()
	if err != nil {
		return errors.WithMessage(err, "Error initialising reception identity")
	}
	err = c.connectCmixx()
	if err != nil {
		return errors.WithMessage(err, "Error connecting to CMixx")
	}
	c.initialisePartner(partnerPubKey, partnerToken)
	return nil
}

func (c *CMixxTransport) initialiseCMixx() error {
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
	//Generate directId
	ekv := net.GetStorage().GetKV()
	directIdStorObj, err := ekv.Get("directId", 0)
	if err != nil && ekv.Exists(err) {
		return errors.WithMessage(err, "Error getting directId from storage")
	}

	var directId codename.PrivateIdentity
	if ekv.Exists(err) {
		directId, err = codename.UnmarshalPrivateIdentity(directIdStorObj.Data)
	} else {
		rng := net.GetRng().GetStream()
		defer rng.Close()
		directId, err = codename.GenerateIdentity(rng)
		if err != nil {
			return errors.WithMessage(err, "Error generating directId")
		}
		ekv.Set("directId", &versioned.Object{
			Data:      directId.Marshal(),
			Version:   0,
			Timestamp: time.Now(),
		})
	}

	privateEdwardsKey := directId.Privkey
	token := directId.GetDMToken()
	privateKey := ecdh.Edwards2EcdhNikePrivateKey(privateEdwardsKey)
	publicKey := ecdh.ECDHNIKE.DerivePublicKey(privateKey)

	c.privateKey = privateKey
	c.publicKey = publicKey
	c.receptionId = deriveReceptionID(publicKey.Bytes(), token)
	c.selfReceptionId = deriveReceptionID(publicKey.Bytes(), token)
	return nil
}

func (c *CMixxTransport) initialisePartner(partnerPubKey string, partnerToken string) error {
	decodedPubBytes, err := base64.StdEncoding.DecodeString(partnerPubKey)
	if err != nil {
		return errors.New("Error decoding partner public key")
	}
	publicKey, err := ecdh.ECDHNIKE.UnmarshalBinaryPublicKey(decodedPubBytes)
	if err != nil {
		return errors.WithMessage(err, "Error unmarshalling partner public key")
	}
	token := uint32(c.partnerToken)
	c.partnerPublicKey = ecdh.EcdhNike2EdwardsPublicKey(publicKey)
	c.partnerToken = token
	return nil
}

func (c *CMixxTransport) initialiseReceptionIdentity() error {
	// Get reception identity (automatically created if one does not exist)
	identityStorageKey := "identityStorageKey"
	identity, err := xxdk.LoadReceptionIdentity(identityStorageKey, c.cMixxNet)
	if err != nil {
		// If no extant xxdk.ReceptionIdentity, generate and store a new one
		identity, err = xxdk.MakeReceptionIdentity(c.cMixxNet)
		if err != nil {
			return errors.WithMessage(err, "Failed to generate reception identity")
		}
		err = xxdk.StoreReceptionIdentity(
			identityStorageKey, identity, c.cMixxNet)
		if err != nil {
			return errors.WithMessage(err, "Failed to store new reception identity")
		}
	}

	// Create an E2E client
	// The connect package handles AuthCallbacks,
	// xxdk.DefaultAuthCallbacks is fine here
	params := xxdk.GetDefaultE2EParams()
	user, err := xxdk.Login(c.cMixxNet, xxdk.DefaultAuthCallbacks{}, identity, params)
	if err != nil {
		return errors.WithMessage(err, "Failed to create E2E client")
	}
	c.user = *user
	return nil
}

func (c *CMixxTransport) connectCmixx() error {
	networkFollowerTimeout := 5 * time.Second
	err := c.user.StartNetworkFollower(networkFollowerTimeout)
	if err != nil {
		return errors.WithMessage(err, "Failed to start network follower")
	}

	// Set up a wait for the network to be connected
	waitUntilConnected := func(connected chan bool) error {
		waitTimeout := 30 * time.Second
		timeoutTimer := time.NewTimer(waitTimeout)
		isConnected := false
		// Wait until we connect or panic if we cannot before the timeout
		for !isConnected {
			select {
			case isConnected = <-connected:
				fmt.Printf("Network Status: %v", isConnected)
				break
			case <-timeoutTimer.C:
				return errors.New("Network connection timeout")
			}
		}
		return nil
	}
	// Create a tracker channel to be notified of network changes
	connected := make(chan bool, 10)
	// Provide a callback that will be signalled when network
	// health status changes
	c.user.GetCmix().AddHealthCallback(
		func(isConnected bool) {
			connected <- isConnected
		})
	// Wait until connected or crash on timeout
	err = waitUntilConnected(connected)
	if err != nil {
		return errors.WithMessage(err, "Failed to connect to network")
	}
	return nil
}

func (c *CMixxTransport) Send(content []byte) error {
	_, _, err := c.send(doxxier_init, content)
	return err
}

func (c *CMixxTransport) getNdf() (string, error) {
	//Read NDF file
	rawSignedMarshalledNdf, err := os.ReadFile(c.ndfPath)
	if err != nil {
		return "", errors.WithMessage(err, "Error reading NDF file")
	}

	signedMarshalledNdf, err := base64.StdEncoding.DecodeString(string(rawSignedMarshalledNdf))
	if err != nil {
		return "", errors.WithMessage(err, "Error reading NDF file")
	}
	signedNdfMsg := &pb.NDF{}
	err = proto.Unmarshal(signedMarshalledNdf, signedNdfMsg)
	if err != nil {
		return "", errors.WithMessage(err, "Error reading NDF file")
	}

	cert, err := os.ReadFile(c.certPath)
	if err != nil {
		return "", errors.WithMessage(err, "Error reading cert file")
	}
	certificate, err := tls.LoadCertificate(string(cert))

	if err != nil {
		return "", errors.WithMessage(err, "Error reading cert file")
	}

	pubKey, err := tls.ExtractPublicKey(certificate)
	if err != nil {
		return "", errors.WithMessage(err, "Error reading cert file")
	}

	err = signature.VerifyRsa(signedNdfMsg, pubKey)
	if err != nil {
		return "", errors.WithMessage(err, "Error verifying NDF")
	}
	return string(signedNdfMsg.Ndf), nil
}

func (c *CMixxTransport) Stop() error {
	c.user.StopNetworkFollower()
	return nil
}
