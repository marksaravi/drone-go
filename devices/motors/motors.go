package motors

import (
	"log"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/drivers"
	"github.com/marksaravi/drone-go/drivers/pca9685"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

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

func NewESC() interface {
	On()
	Off()
	SetThrottles(map[uint8]float32)
} {
	flightControlConfigs := config.ReadFlightControlConfig()
	escConfigs := flightControlConfigs.Configs.ESC
	powerbreaker := drivers.NewGPIOOutput(flightControlConfigs.Configs.PowerBreaker)
	b, _ := i2creg.Open(escConfigs.I2CDev)
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	pwmDev, err := pca9685.NewPCA9685(pca9685.PCA9685Address, i2cConn, escConfigs.MaxThrottle)
	if err != nil {
		log.Fatal(err)
	}
	esc := NewMotorsControl(pwmDev, powerbreaker, escConfigs.MotorESCMappings)
	return esc
}
