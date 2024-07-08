package utils

type kalmanFilter struct {
	prevValue float64
	gain      float64
}

func NewKalmanFilter(gain float64) *kalmanFilter {
	return &kalmanFilter{
		prevValue: 0,
		gain:      gain,
	}
}

func (k *kalmanFilter) SmoothValue(value float64) float64 {
	smoothValue := k.prevValue*(1-k.gain) + value*k.gain
	k.prevValue = smoothValue
	return smoothValue
}
