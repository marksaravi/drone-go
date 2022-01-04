package pid

type pidControl struct {
	pGain    float64
	iGain    float64
	dGain    float64
	limit    float64
	previous float64
	memory   float64
}

func NewPIDControl(pGain float64, iGain float64, dGain float64, limit float64) *pidControl {
	return &pidControl{
		pGain:    pGain,
		iGain:    iGain,
		dGain:    dGain,
		limit:    limit,
		previous: 0,
		memory:   0,
	}
}
