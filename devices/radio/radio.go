package radio

import "github.com/marksaravi/drone-go/models"

const (
	IDLE models.ConnectionState = iota
	WAITING_FOR_CONNECTION
	DISCONNECTED
	CONNECTED
)

type radioLink interface {
	PowerOn()
	PowerOff()
	ClearStatus()
}

func StateToString(s models.ConnectionState) string {
	if s == WAITING_FOR_CONNECTION {
		return "WAITING_FOR_CONNECTION"
	}
	if s == DISCONNECTED {
		return "DISCONNECTED"
	}
	return "CONNECTED"
}
