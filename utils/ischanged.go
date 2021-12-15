package utils

import (
	"math"

	"github.com/marksaravi/drone-go/models"
)

func compateFloats(v1, v2 float32, minChange float32) bool {
	return math.Abs(float64(v1-v2)) > float64(minChange)
}

func IsFlightCommandChaned(fc1, fc2 models.FlightCommands, minValue float32) bool {
	var isChanged bool = compateFloats(fc1.Roll, fc2.Roll, minValue) ||
		compateFloats(fc1.Pitch, fc2.Pitch, minValue) ||
		compateFloats(fc1.Yaw, fc2.Yaw, minValue) ||
		compateFloats(fc1.Throttle, fc2.Throttle, minValue) ||
		fc1.ButtonFrontLeft != fc2.ButtonFrontLeft ||
		fc1.ButtonTopLeft != fc2.ButtonTopLeft ||
		fc1.ButtonBottomLeft != fc2.ButtonBottomLeft ||
		fc1.ButtonFrontRight != fc2.ButtonFrontRight ||
		fc1.ButtonTopRight != fc2.ButtonTopRight ||
		fc1.ButtonBottomRight != fc2.ButtonBottomRight
	return isChanged
}
