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
	roll     float64
	pitch    float64
	yaw      float64
	throttle float64
}

type pidControl struct {
	pGain       float64
	iGain       float64
	dGain       float64
	conversions analogToDigitalConversions
	commands    pidCommands

	rotations     models.ImuRotations
	prevRotations models.ImuRotations
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
		throttles: map[uint8]float32{0: 0, 1: 0, 2: 0, 3: 0},
	}
}

func (pid *pidControl) ApplyFlightCommands(flightCommands models.FlightCommands) {
	commands := flightControlCommandToPIDCommand(flightCommands, pid.conversions)
	pid.calcThrottlesByCommands(commands)
}

func (pid *pidControl) ApplyRotations(rotations models.ImuRotations) {
	pid.calcThrottlesByFlightData(rotations)
}

func (pid *pidControl) calcThrottlesByFlightData(rotations models.ImuRotations) {
	pid.prevRotations = pid.rotations
	pid.rotations = rotations
}

func (pid *pidControl) calcThrottlesByCommands(commands pidCommands) {
	pid.commands = commands
	t := float32(pid.commands.throttle)
	pid.throttles = map[uint8]float32{
		0: t,
		1: t,
		2: t,
		3: t,
	}
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
