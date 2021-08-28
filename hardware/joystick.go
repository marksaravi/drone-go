package hardware

type analogToDigitalConverter interface {
	ReadInputVoltage(channel int, vRef float32, zeroValue float32) (float32, error)
}

type joystick struct {
	dev       analogToDigitalConverter
	channel   int
	zeroValue float32
	vRef      float32
	value     float32
}

func (j *joystick) Read() float32 {
	value, error := j.dev.ReadInputVoltage(j.channel, j.vRef, j.zeroValue)
	if error == nil {
		j.value = value
	}
	return j.value
}

func NewJoystick(dev analogToDigitalConverter, channel int, zeroValue float32, vRef float32) *joystick {
	return &joystick{
		dev:       dev,
		channel:   channel,
		zeroValue: zeroValue,
		vRef:      vRef,
	}
}
