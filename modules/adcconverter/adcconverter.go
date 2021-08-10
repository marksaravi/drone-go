package adcconverter

type AnalogToDigitalConverter interface {
	ReadInputVoltage(int, float32) (float32, error)
}
