package radio

import (
	"github.com/marksaravi/drone-go/constants"
	"github.com/marksaravi/drone-go/models"
)

type radioLink interface {
	PowerOn()
	PowerOff()
	ClearStatus()
}

func StateToString(s models.ConnectionState) string {
	if s == constants.WAITING_FOR_CONNECTION {
		return "WAITING_FOR_CONNECTION"
	}
	if s == constants.DISCONNECTED {
		return "DISCONNECTED"
	}
	return "CONNECTED"
}
