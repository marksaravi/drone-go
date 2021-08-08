package types

// Config is the generic configuration
type Config interface {
}

type PID interface {
	Update(ImuRotations) []Throttle
}

// Logger is interface for the udpLogger
type UdpLogger interface {
	Send(ImuRotations)
}
