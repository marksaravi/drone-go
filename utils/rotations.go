package utils

import (
	"math"

	"github.com/MarkSaravi/drone-go/types"
)

func RotationsRadToDeg(e types.Rotations) types.Rotations {
	return types.Rotations{
		Roll:  e.Roll / math.Pi * 180,
		Pitch: e.Pitch / math.Pi * 180,
	}
}

func AccelerometerToRotations(acc types.XYZ) (types.Rotations, error) {
	Roll := math.Atan2(acc.X, acc.Z)
	Pitch := math.Atan2(acc.Y, acc.Z)

	return RotationsRadToDeg(types.Rotations{
		Roll:  Roll,
		Pitch: Pitch,
	}), nil
}
