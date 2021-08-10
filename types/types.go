package types

// Logger is interface for the udpLogger
type UdpLogger interface {
	Send(ImuRotations)
}
