package adcconverter

type AnalogToDigitalConverter interface {
	ReadInputVoltage(int) float32
}

type AnalogToDigitalDevice interface {
	ReadInput(int) int
}

type adc struct {
	dev AnalogToDigitalDevice
}

func NewADCConverter(dev AnalogToDigitalDevice) AnalogToDigitalConverter {
	return &adc{
		dev: dev,
	}
}

func (a *adc) ReadInputVoltage(channel int) float32 {
	digitalValue := a.dev.ReadInput(channel)
	return float32(digitalValue)
}
