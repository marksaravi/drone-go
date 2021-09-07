package motors

const (
	NUM_OF_MOTORS uint8 = 4
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

func (mc *motorsControl) SetThrottles(throttles map[uint8]float32) {
	var motor uint8
	for motor = 0; motor < NUM_OF_MOTORS; motor++ {
		mc.setThrottle(motor, throttles[uint8(motor)])
	}
}

func (mc *motorsControl) setThrottle(motor uint8, throttle float32) {
	mc.esc.SetThrottle(mc.motorsEscMappings[int(motor)], throttle)
}

func (mc *motorsControl) allMotorsOff() {
	var i uint8
	for i = 0; i < NUM_OF_MOTORS; i++ {
		mc.setThrottle(i, 0)
	}
}
