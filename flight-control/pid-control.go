package flightcontrol

import "github.com/MarkSaravi/drone-go/types"

type pid struct {
	power         float32
	currRotations types.ImuRotations
	prevRotations types.ImuRotations
}

func CreatePidController() types.PID {
	return &pid{}
}

func (c *pid) Update(r types.ImuRotations) []types.Throttle {
	c.prevRotations = c.currRotations
	c.currRotations = r
	return []types.Throttle{
		{Motor: 0, Value: 0},
		{Motor: 1, Value: 0},
		{Motor: 2, Value: 0},
		{Motor: 3, Value: 0},
	}
}
