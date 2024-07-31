package network

type PrivacyLevel int

const (
	LOW    PrivacyLevel = 0
	MEDIUM PrivacyLevel = 1
	HIGH   PrivacyLevel = 2
)

var transportMap = make(map[PrivacyLevel]Transport)

func RegisterTransport(privacyLevel PrivacyLevel, transport Transport) {
	transportMap[privacyLevel] = transport
}

func GetTransportation(privacyLevel PrivacyLevel) Transport {
	return transportMap[privacyLevel]
}
