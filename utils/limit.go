package utils

import "math"

func Limit(value, limit float64) float64 {
	if math.Abs(value) <= math.Abs(limit) {
		return value
	}
	if value < 0 {
		return -math.Abs(limit)
	}
	return math.Abs(limit)
}
