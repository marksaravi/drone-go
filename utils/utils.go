package utils

import "math"

func SignedMax(v, maxValue float64) float64 {
	if math.Abs(v) < maxValue {
		return v
	}
	if v < 0 {
		return -maxValue
	}
	return maxValue
}
