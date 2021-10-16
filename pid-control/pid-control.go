package pidcontrol

import (
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/models"
)

type analogToDigitalConversion struct {
	ratio  float64
	offset float64
}
type analogToDigitalConversions struct {
	roll     analogToDigitalConversion
	pitch    analogToDigitalConversion
	yaw      analogToDigitalConversion
	throttle analogToDigitalConversion
}

type pidCommands struct {
	throttle float64
	roll     float64
	pitch    float64
	yaw      float64
}

type pidControl struct {
	pGain       float64
	iGain       float64
	dGain       float64
	conversions analogToDigitalConversions
	commands    pidCommands

	rotations     models.ImuRotations
	prevRotations models.ImuRotations
	throttle      float64
	iThrottle     float64
	throttles     map[uint8]float32
}

func NewPIDControl() *pidControl {
	configs := config.ReadFlightControlConfig().Configs.PID
	return &pidControl{
		pGain: configs.PGain,
		iGain: configs.IGain,
		dGain: configs.DGain,
		conversions: analogToDigitalConversions{
			roll: analogToDigitalConversion{
				ratio:  configs.AnalogInputToRoll.Ratio,
				offset: configs.AnalogInputToRoll.Offset,
			},
			pitch: analogToDigitalConversion{
				ratio:  configs.AnalogInputToPitch.Ratio,
				offset: configs.AnalogInputToPitch.Offset,
			},
			yaw: analogToDigitalConversion{
				ratio:  configs.AnalogInputToYaw.Ratio,
				offset: configs.AnalogInputToYaw.Offset,
			},
			throttle: analogToDigitalConversion{
				ratio:  configs.AnalogInputToThrottle.Ratio,
				offset: configs.AnalogInputToThrottle.Offset,
			},
		},
	}
}

func (pid *pidControl) ApplyFlightCommands(flightCommands models.FlightCommands) {
	pid.commands = flightControlCommandToPIDCommand(flightCommands, pid.conversions)
}

func (pid *pidControl) ApplyRotations(rotations models.ImuRotations) {
	pid.prevRotations = pid.rotations
	pid.rotations = rotations
	if t, err := pid.calcMotorsThrottles(); err == nil {
		pid.throttles = t
	}
}

func (pid *pidControl) calcMotorsThrottles() (map[uint8]float32, error) {
	throttle := float32(pid.commands.throttle)
	return map[uint8]float32{
		0: throttle,
		1: throttle,
		2: throttle,
		3: throttle,
	}, nil
}

func (pid *pidControl) Throttles() map[uint8]float32 {
	return pid.throttles
}

func convert(value float32, conversion analogToDigitalConversion) float64 {
	return float64(value)*conversion.ratio + conversion.offset
}

func flightControlCommandToPIDCommand(c models.FlightCommands, conversions analogToDigitalConversions) pidCommands {
	return pidCommands{
		roll:     convert(c.Roll, conversions.roll),
		pitch:    convert(c.Pitch, conversions.pitch),
		yaw:      convert(c.Yaw, conversions.yaw),
		throttle: convert(c.Throttle, conversions.throttle),
	}
}
