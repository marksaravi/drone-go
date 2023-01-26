package devices

type gpiooutput interface {
	SetHigh()
	SetLow()
}

type powerbreaker struct {
	gpiopin gpiooutput
}

func NewPowerBreaker(gpiopin gpiooutput) *powerbreaker {
	return &powerbreaker{
		gpiopin: gpiopin,
	}
}

func (pb *powerbreaker) Connect() {
	pb.gpiopin.SetHigh()
}

func (pb *powerbreaker) Disconnect() {
	pb.gpiopin.SetLow()
}
