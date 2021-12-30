package devices

import "math"

type analogtodigital interface {
	Read() uint16
	GetDigitalMaxValue() uint16
}

type joystickInput struct {
	input        analogtodigital
	aCoefficient float64
	bCoefficient float64
}

func (js *joystickInput) Read() int {
	digitalValue := js.input.Read()
	return int(digitalValue)
}

func NewJoystick(input analogtodigital, digitakMidValue uint16) *joystickInput {
	calcCoefficients(digitakMidValue, input.GetDigitalMaxValue())

	return &joystickInput{
		input: input,
	}
}

func calcCoefficients(digitalMidValue uint16, digitalMaxValue uint16) (float64, float64) {
	dMidValue := int(digitalMidValue)
	dMaxValue := int(digitalMaxValue)
	midValue := dMaxValue / 2
	x1 := float64(dMidValue)
	y1 := float64(midValue)
	x2 := float64(dMaxValue)
	y2 := float64(dMaxValue)
	x1_2 := x1 * x1
	x2_2 := x2 * x2
	k1 := y1 / x1_2
	bCoefficientRaw := (y2 - x2_2*k1) / (x2 - x2_2/x1)
	aCoefficientRaw := (y1 - bCoefficientRaw*x1) / x1_2
	const ROUND_FACTOR = 1000000000
	aCoefficient := math.Round(aCoefficientRaw*ROUND_FACTOR) / ROUND_FACTOR
	bCoefficient := math.Round(bCoefficientRaw*ROUND_FACTOR) / ROUND_FACTOR
	return float64(aCoefficient), float64(bCoefficient)
}

func calcValue(digitalValue uint16, aCoefficient float64, bCoefficient float64) uint16 {
	x := float64(digitalValue)
	value := aCoefficient*x*x + bCoefficient*x
	return uint16(value)
}
