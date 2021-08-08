package types

type Throttle struct {
	Motor int
	Value float32
}

type PowerBreaker interface {
	Connect()
	Disconnect()
}

type ESC interface {
	SetThrottle(int, float32)
}

type MotorsController interface {
	On()
	Off()
	SetThrottles(map[int]float32)
}
