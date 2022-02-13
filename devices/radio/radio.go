package radio

type ConnectionState = int

const (
	IDLE ConnectionState = iota
	DISCONNECTED
	CONNECTED
)

type radioLink interface {
	PowerOn()
	PowerOff()
	ClearStatus()
}
