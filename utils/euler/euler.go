package euler

import (
	"math"

	"github.com/MarkSaravi/drone-go/types"
)

func EulerRadToDeg(e types.Euler) types.Euler {
	return types.Euler{
		Theta: e.Theta / math.Pi * 180,
	}
}

func AccelerometerToEulerAngles(acc types.XYZ) (types.Euler, error) {
	theta := math.Atan2(acc.X, acc.Z)

	return EulerRadToDeg(types.Euler{
		Theta: theta,
	}), nil
}
