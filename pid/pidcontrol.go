package pid

type pidControl struct {
	limit    float64
	previous float64
	memory   float64
}

func NewPIDControl(limit float64) *pidControl {
	return &pidControl{
		limit:    limit,
		previous: 0,
		memory:   0,
	}
}

func (c *pidControl) calc(rotation, targetRotation float64, gains *gains) float64 {
	rotationDiff := targetRotation - rotation
	p := gains.P * rotationDiff
	sum := p
	return sum
}
