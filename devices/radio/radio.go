package radio

import (
	"github.com/marksaravi/drone-go/constants"
)

type radioLink interface {
	PowerOn()
	PowerOff()
	ClearStatus()
}

func StateToString(s int) string {
	if s == constants.WAITING_FOR_CONNECTION {
		return "WAITING_FOR_CONNECTION"
	}
	if s == constants.DISCONNECTED {
		return "DISCONNECTED"
	}
	if s == constants.IDLE {
		return "IDLE"
	}
	return "CONNECTED"
}
