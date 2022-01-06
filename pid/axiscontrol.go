package pid

type axisControl struct {
	limit    float64
	limitI   float64
	previous float64
	memory   float64
}

func NewPIDControl(limit, limitI float64) *axisControl {
	return &axisControl{
		limit:    limit,
		limitI:   limitI,
		previous: 0,
		memory:   0,
	}
}

func (c *axisControl) calc(rotation, targetRotation float64, gains *gains) float64 {
	rotationDiff := targetRotation - rotation
	p := gains.P * rotationDiff
	sum := p
	return sum
}
