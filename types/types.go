package types

import "github.com/MarkSaravi/drone-go/modules/imu"

// Logger is interface for the udpLogger
type UdpLogger interface {
	Send(imu.ImuRotations)
}
