package types

// Config is the generic configuration
type Config interface {
}

// Logger is interface for the udpLogger
type UdpLogger interface {
	Send(ImuRotations)
}
