package radio

type ConnectionState = int

const (
	IDLE ConnectionState = iota
	WAITING_FOR_CONNECTION
	DISCONNECTED
	CONNECTED
)

type radioLink interface {
	PowerOn()
	PowerOff()
	ClearStatus()
}

func StateToString(s ConnectionState) string {
	if s == WAITING_FOR_CONNECTION {
		return "WAITING_FOR_CONNECTION"
	}
	if s == DISCONNECTED {
		return "DISCONNECTED"
	}
	return "CONNECTED"
}
