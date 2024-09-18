package cmixx

import (
	"time"

	"gitlab.com/elixxir/client/v4/cmix"
	"gitlab.com/elixxir/client/v4/cmix/rounds"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/id/ephemeral"
)

type Cmix interface {
	SendMany(messages []cmix.TargetedCmixMessage,
		p cmix.CMIXParams) (rounds.Round, []ephemeral.Id, error)
	IsHealthy() bool
	GetRoundResults(timeout time.Duration,
		roundCallback cmix.RoundEventCallback, roundList ...id.Round)
	SendManyWithAssembler(recipients []*id.ID, assembler cmix.ManyMessageAssembler,
		params cmix.CMIXParams) (rounds.Round, []ephemeral.Id, error)
	GetMaxMessageLength() int
}
