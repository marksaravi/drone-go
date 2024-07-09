package utils

import (
	"math"
)

type kalmanFilter struct {
	prevValue        float64
	decimalPrecision float64
	gain             float64
}

func NewKalmanFilter(gain float64, decimalPrecision int) *kalmanFilter {
	return &kalmanFilter{
		prevValue:        0,
		gain:             gain,
		decimalPrecision: math.Pow10(decimalPrecision),
	}
}

func (k *kalmanFilter) Value(value float64) float64 {

	value = math.Round(value*k.decimalPrecision) / k.decimalPrecision
	smoothValue := k.prevValue*(1-k.gain) + value*k.gain
	k.prevValue = smoothValue
	return smoothValue
}
