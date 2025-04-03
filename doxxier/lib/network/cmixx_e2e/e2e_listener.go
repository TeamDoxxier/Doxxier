package cmixx_e2e

import (
	"fmt"

	"gitlab.com/elixxir/client/v4/e2e/receive"
)

// listener implements the receive.Listener interface
type E2eListener struct {
	name string
}

// Hear will be called whenever a message matching
// the RegisterListener call is received
// User-defined message handling logic goes here
func (l E2eListener) Hear(item receive.Message) {
	fmt.Printf("Message received: %v", item)
}

// Name is used for debugging purposes
func (l E2eListener) Name() string {
	return l.name
}
