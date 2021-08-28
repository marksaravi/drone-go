package oldflightcontrol

import (
	"sync"

	"github.com/MarkSaravi/drone-go/modules/command"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/modules/motors"
	"github.com/MarkSaravi/drone-go/modules/radiolink"
	"github.com/MarkSaravi/drone-go/modules/udplogger"
)

type EscConfig struct {
	UpdateFrequency int     `yaml:"update_frequency"`
	MaxThrottle     float32 `yaml:"max_throttle"`
}

type PidConfig struct {
	ProportionalGain float32 `yaml:"proportional–gain"`
	IntegralGain     float32 `yaml:"integral-gain"`
	DerivativeGain   float32 `yaml:"derivative-gain"`
}

type FlightConfig struct {
	PID PidConfig     `yaml:"pid"`
	Imu imu.ImuConfig `yaml:"imu"`
	Esc EscConfig     `yaml:"esc"`
}

type flightControl struct {
	imu              imu.IMU
	motorsController motors.MotorsController
	radio            radiolink.RadioLink
	logger           udplogger.UdpLogger
}

func (fc *flightControl) Start() {
	var wg sync.WaitGroup

	commandChannel := command.CreateCommandChannel(&wg)

	fc.imu.ResetReadingTimes()
	var running bool = true

	for running {
		if fc.imu.CanRead() {
			rotations, err := fc.imu.GetRotations()
			if err == nil {
				fc.logger.Send(rotations)
			}
		}
		select {
		case command := <-commandChannel:
			if command.Command == "COMMAND_END_PROGRAM" {
				wg.Wait()
				running = false
			}
		default:
		}
	}
}

func CreateFlightControl(
	imu imu.IMU,
	motorsController motors.MotorsController,
	radio radiolink.RadioLink,
	logger udplogger.UdpLogger,
) *flightControl {
	return &flightControl{
		imu:              imu,
		motorsController: motorsController,
		radio:            radio,
		logger:           logger,
	}
}
