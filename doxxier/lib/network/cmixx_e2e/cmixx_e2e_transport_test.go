package network

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/elixxir/client/v4/xxdk"
	"gitlab.com/xx_network/primitives/id"
)

func TestNewCMixxE2eTransport(t *testing.T) {
	tests := []struct {
		name   string
		config CMixxConfig
		want   *CmixxE2eTransport
	}{
		{
			name: "All fields provided",
			config: CMixxConfig{
				NdfPath:            "../assets/ndf.json",
				CertPath:           "../assets/main.cert",
				Secret:             "cE+a6mP>p`4b]F8;CY&*t}",
				StoragePath:        "../assets/storage/",
				identityStorageKey: "customIdentityStorageKey",
			},
			want: &CmixxE2eTransport{
				secret:             "cE+a6mP>p`4b]F8;CY&*t}",
				ndfPath:            "../assets/ndf.json",
				certPath:           "../assets/main.cert",
				StoragePath:        "../assets/storage/",
				identityStorageKey: "customIdentityStorageKey",
			},
		},
		{
			name: "Missing optional fields",
			config: CMixxConfig{
				StoragePath: "../assets/storage/",
			},
			want: &CmixxE2eTransport{
				secret:             secretFlag,
				ndfPath:            ndfPathFlag,
				certPath:           certPathFlag,
				StoragePath:        "../assets/storage/",
				identityStorageKey: identityStorageKeyFlag,
			},
		},
		{
			name: "Missing StoragePath",
			config: CMixxConfig{
				NdfPath:            "custom/ndf/path",
				CertPath:           "custom/cert/path",
				Secret:             "customSecret",
				identityStorageKey: "customIdentityStorageKey",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.want != nil {
						t.Errorf("NewCMixxE2eTransport() panicked when it shouldn't have")
					}
				}
			}()
			got, _ := NewCMixxE2eTransport(tt.config)
			if tt.want == nil {
				if got != nil {
					t.Errorf("NewCMixxE2eTransport() = %v, want %v", got, tt.want)
				}
				return
			}
			if got.secret != tt.want.secret ||
				got.ndfPath != tt.want.ndfPath ||
				got.certPath != tt.want.certPath ||
				got.StoragePath != tt.want.StoragePath ||
				got.identityStorageKey != tt.want.identityStorageKey {
				t.Errorf("NewCMixxE2eTransport() = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockCmixxE2eTransport struct {
	mock.Mock
	CmixxE2eTransport
}

func (m *MockCmixxE2eTransport) initialiseCMixx() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCmixxE2eTransport) initialiseReceptionIdentity() (*xxdk.ReceptionIdentity, error) {
	args := m.Called()
	return args.Get(0).(*xxdk.ReceptionIdentity), args.Error(1)
}

func TestCmixxE2eTransport_Connect(t *testing.T) {
	tests := []struct {
		name                  string
		initialiseCMixxError  error
		initialiseIdentityErr error
		expectedError         string
		expectedOutput        string
		config                CMixxConfig
	}{
		{
			name: "Successful connection",
			config: CMixxConfig{
				StoragePath: "../assets/storage",
			},
			initialiseCMixxError:  nil,
			initialiseIdentityErr: nil,
			expectedError:         "",
			expectedOutput:        "PHh4YygyKUt3TExpaVpaV2xqaHdhMWxVK3I0aU5WOGJ0a3Q1NytJUkFRa2ErTzNuSU1Ea0FaaUI5a1pvK0RsM2ZYM2QzV3dDeTRLc2hZT1lWYzZKL1U3ODVURXpsNVFQYVhad04xeVZYRlpDZTNqS2cxeGR2RTY4L2hCR1Z6bWJwcWtGTXdIZStOVzRHaUh4aSs5aUcyZHVqQkdOeHEzTUpDRXZqOHlVZlNlb2dtTGIwbFc4eUozOE9BYzdNenR4U25ESDQwZmxITG9ZQ0pEUEFxNnl2azNYaGxWMXNGM1JhYXhBakdiejhML2xZc2NJNHVvMEZKSnZWSlZSQ2tYQVV2bkRFeEsvUjY0eTBRZWNVbVF1SjhqWGVFU2tUNU9SVU9qelpOQnZqajlEVTBlenZPWFAzT0Y2VC9UMTJad25BdlAvQmZjOE5HSW5UdDFSNkZVZzloSWdZdEJwUEUvN3R5dlNUbVM5TVdNYlVIM1J6S2dBZVQvdnNaUXRCQW1pRjIrSFFMeWVvQXRhY3ROZ0lmV3cyYUw3UGhGeVByOUxwSVRCdXczMzVIbm9pZ0lBeHF5UVVlTVpUTnA5Wk4vbWhhOFk2Ry8rTHZrbEM5bVZacHVuK1JmdENORlVuMThMNk5FYmFvS0ZGcFpOVU1sbkZrSGpBcjV0bUlaczg4S2l3U2JSR2Q4S0d4dUFXVzh3QXAzMTZmTkh1U2wzelV3R2tjZ3pMN1BCandjTk9nY1NoUVMzeG1KdW1VZ0RqL3dmQUFBQWdBN2VhM3d1RFpHMS9UZU43ZDNjTTMwbHc9PXh4Yz4=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransport := new(MockCmixxE2eTransport)
			transport, _ := NewCMixxE2eTransport(tt.config)
			mockTransport.CmixxE2eTransport = *transport
			mockTransport.On("initialiseCMixx").Return(tt.initialiseCMixxError)
			mockIdentity := &xxdk.ReceptionIdentity{}
			mockTransport.On("initialiseReceptionIdentity").Return(mockIdentity, tt.initialiseIdentityErr)

			if tt.initialiseIdentityErr == nil {
				mockTransport.On("initialiseReceptionIdentity").Return(mockIdentity, nil)
			}

			testChan := make(chan string)
			output, err := mockTransport.Connect(testChan)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, output)
			}

			mockTransport.AssertExpectations(t)
		})
	}
}

func TestCmixxE2eTransport_connectCmixx(t *testing.T) {
	tests := []struct {
		name               string
		loginError         error
		startFollowerError error
		expectedError      string
	}{
		{
			name:               "Successful connection",
			loginError:         nil,
			startFollowerError: nil,
			expectedError:      "",
		},
		{
			name:               "Login error",
			loginError:         errors.New("login error"),
			startFollowerError: nil,
			expectedError:      "Error logging in to cMixx: login error",
		},
		{
			name:               "Start network follower error",
			loginError:         nil,
			startFollowerError: errors.New("start follower error"),
			expectedError:      "Error starting network follower: start follower error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransport := new(MockCmixxE2eTransport)
			transport, _ := NewCMixxE2eTransport(CMixxConfig{StoragePath: "../assets/storage"})
			mockTransport.CmixxE2eTransport = *transport

			mockUser := new(MockUser)
			mockUser.On("StartNetworkFollower", mock.Anything).Return(tt.startFollowerError)
			mockUser.On("GetE2E").Return(new(MockE2E))
			mockUser.On("GetCmix").Return(new(MockCmix))

			mockTransport.On("initialiseCMixx").Return(nil)
			mockTransport.On("initialiseReceptionIdentity").Return(new(xxdk.ReceptionIdentity), nil)
			mockTransport.On("Login", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockUser, tt.loginError)

			err := mockTransport.connectCmixx()

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockTransport.AssertExpectations(t)
			mockUser.AssertExpectations(t)
		})
	}
}

type MockUser struct {
	mock.Mock
}

func (m *MockUser) StartNetworkFollower(timeout int) error {
	args := m.Called(timeout)
	return args.Error(0)
}

func (m *MockUser) GetE2E() *xxdk.E2e {
	args := m.Called()
	return args.Get(0).(*xxdk.E2e)
}

func (m *MockUser) GetCmix() *xxdk.Cmix {
	args := m.Called()
	return args.Get(0).(*xxdk.Cmix)
}

type MockE2E struct {
	mock.Mock
}

func (m *MockE2E) RegisterListener(id *id.ID, catalogType catalog.MessageType, listener xxdk.Listener) {
	m.Called(id, catalogType, listener)
}

type MockCmix struct {
	mock.Mock
}

func (m *MockCmix) AddHealthCallback(callback func(bool)) {
	m.Called(callback)
}
