package utils

type valueTypes interface {
	float64 | int
}

type average[T valueTypes] struct {
	index  int
	values []T
	sum    T
}

func NewAverage[T valueTypes](numOfSamples int) *average[T] {
	return &average[T]{
		index:  0,
		values: make([]T, numOfSamples),
	}
}

func (a *average[T]) AddValue(v T) T {
	a.sum -= a.values[a.index]
	a.values[a.index] = v
	a.sum += v
	a.index++
	if a.index == len(a.values) {
		a.index = 0
	}
	return a.Average()
}

func (a *average[T]) Average() T {
	return a.sum / T(len(a.values))
}
