package types

type Throttle struct {
	Motor int
	Value float32
}

type PowerBreaker interface {
	Connect()
	Disconnect()
}
