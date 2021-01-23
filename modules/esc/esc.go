package esc

const (
	MinPW float32 = 0.001
	MaxPW float32 = 0.002
	Frequency float32 = 400
)

//PWMDevice is the electronic board that generate PWM
type PWMDevice interface {
	Start(frequency float32) error
	SetPulseWidth(channel int, pulseWidth float32)
	SetPulseWidthAll(pulseWidth float32)
	StopAll()
	Halt() error
	Close()
}

//ESC is the PWM manager
type ESC struct {
	PWMDevice
}

//NewESC create an ESC
func NewESC(device PWMDevice) *ESC {
	return &ESC{
		PWMDevice: device,
	}
}
