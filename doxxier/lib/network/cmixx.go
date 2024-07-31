package network

import (
	"errors"
	"os"
	"time"

	"gitlab.com/elixxir/client/v4/collective/versioned"
	"gitlab.com/elixxir/client/v4/xxdk"
	"gitlab.com/elixxir/crypto/codename"
)

type CMixxTransport struct {
	secret   string
	ndfPath  string
	certPath string
	cMixxNet *xxdk.Cmix
	directId codename.PrivateIdentity
}

type CMixxConfig struct {
	NdfPath   string
	CertPath  string
	Secret    string
	Recipient string
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
	return &CMixxTransport{secret: config.Secret, ndfPath: config.NdfPath, certPath: config.CertPath}
}

func (c *CMixxTransport) Connect() error {
	//Read NDF file
	ndf, err := os.ReadFile(c.ndfPath)
	if err != nil {
		return errors.New("Error reading NDF file: " + err.Error())
	}
	//Connect to CMixx
	secretBytes := []byte(c.secret)
	err = xxdk.NewCmix(string(ndf), "xx", secretBytes, "")
	if err != nil {
		return errors.New("Error connecting to CMixx: " + err.Error())
	}
	params := xxdk.GetDefaultCMixParams()
	net, err := xxdk.LoadCmix("xx", secretBytes, params)
	if err != nil {
		return errors.New("Error connecting to CMixx: " + err.Error())
	}
	c.cMixxNet = net
	//Generate directId
	ekv := net.GetStorage().GetKV()
	directIdStorObj, err := ekv.Get("directId", 0)
	if err != nil {
		return errors.New("Error getting directId from storage: " + err.Error())
	}

	var directId codename.PrivateIdentity
	if ekv.Exists(err) {
		directId, err = codename.UnmarshalPrivateIdentity(directIdStorObj.Data)
	} else {
		rng := net.GetRng().GetStream()
		defer rng.Close()
		directId, err = codename.GenerateIdentity(rng)
		if err != nil {
			return errors.New("Error generating directId: " + err.Error())
		}
		ekv.Set("directId", &versioned.Object{
			Data:      directId.Marshal(),
			Version:   0,
			Timestamp: time.Now(),
		})
	}
	c.directId = directId
	return nil
}

func connectToRecipient(recipient string) error {

	return nil
}

func (c *CMixxTransport) Send() error {
	return nil
}
