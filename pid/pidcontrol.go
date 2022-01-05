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
