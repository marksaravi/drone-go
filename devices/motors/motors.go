package motors

const (
	NUM_OF_MOTORS int = 4
)

type powerbreaker interface {
	SetHigh()
	SetLow()
}

type eschandler interface {
	SetThrottle(int, float32)
}

type motorsControl struct {
	esc               eschandler
	powerbreaker      powerbreaker
	motorsEscMappings map[int]int
}

func NewMotorsControl(esc eschandler, powerbreaker powerbreaker, motorsEscMappings map[int]int) *motorsControl {
	return &motorsControl{
		esc:               esc,
		powerbreaker:      powerbreaker,
		motorsEscMappings: motorsEscMappings,
	}
}

func (mc *motorsControl) On() {
	mc.allMotorsOff()
	mc.powerbreaker.SetHigh()
}

func (mc *motorsControl) Off() {
	mc.allMotorsOff()
	mc.powerbreaker.SetLow()
}

func (mc *motorsControl) SetThrottle(motor int, throttle float32) {
	mc.esc.SetThrottle(mc.motorsEscMappings[motor], throttle)
}

func (mc *motorsControl) allMotorsOff() {
	for i := 0; i < NUM_OF_MOTORS; i++ {
		mc.SetThrottle(i, 0)
	}
}
