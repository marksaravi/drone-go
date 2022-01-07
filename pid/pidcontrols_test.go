package pid

import (
	"testing"

	"github.com/marksaravi/drone-go/constants"
)

func TestJoystickToPidValue(t *testing.T) {
	var digitalValue uint16 = 64
	pidcontrols := NewPIDControls(PIDSettings{
		LimitRoll:               16,
		MaxJoystickDigitalValue: constants.JOYSTICK_RESOLUTION,
	})
	var want float64 = -7
	got := pidcontrols.joystickToPidValue(digitalValue, pidcontrols.roll.inputLimit)
	if got != want {
		t.Fatalf("Wanted %3.2f, got %3.2f", want, got)
	}

}

func TestThrottleToPidThrottle(t *testing.T) {
	var digitalValue uint16 = 64
	pidcontrols := NewPIDControls(PIDSettings{
		MaxJoystickDigitalValue: constants.JOYSTICK_RESOLUTION,
		ThrottleLimit:           16,
	})
	var want float64 = 1
	got := pidcontrols.throttleToPidThrottle(digitalValue)
	if got != want {
		t.Fatalf("Wanted %3.2f, got %3.2f", want, got)
	}
}
