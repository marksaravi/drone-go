package radio

type ConnectionState = int

const (
	IDLE ConnectionState = iota
	DISCONNECTED
	CONNECTED
	CONNECTION_LOST
)

type radioLink interface {
	PowerOn()
	PowerOff()
	ClearStatus()
}
