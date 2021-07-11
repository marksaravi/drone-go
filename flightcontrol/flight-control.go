package flightcontrol

import (
	"sync"

	commands "github.com/MarkSaravi/drone-go/constants"
	"github.com/MarkSaravi/drone-go/modules/command"
	"github.com/MarkSaravi/drone-go/types"
)

type flightControl struct {
	imu    types.IMU
	pid    types.PID
	esc    types.ESC
	logger types.UdpLogger
}

func (fc *flightControl) Start() {
	var wg sync.WaitGroup

	commandChannel := command.CreateCommandChannel(&wg)

	fc.imu.ResetReadingTimes()
	var running bool = true

	fc.esc.MotorsOn()
	for running {
		if fc.imu.CanRead() {
			rotations, err := fc.imu.GetRotations()
			if err == nil {
				_ = fc.pid.Update(rotations)
				// fc.esc.SetThrottles(throttles)
				fc.logger.Send(rotations)
			}
		}
		select {
		case command := <-commandChannel:
			if command.Command == commands.COMMAND_END_PROGRAM {
				wg.Wait()
				fc.esc.MotorsOff()
				running = false
			}
		default:
		}
	}
}

func CreateFlightControl(
	imu types.IMU,
	pid types.PID,
	esc types.ESC,
	logger types.UdpLogger,
) *flightControl {
	return &flightControl{
		imu:    imu,
		pid:    pid,
		esc:    esc,
		logger: logger,
	}
}
