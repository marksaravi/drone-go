package flightcontrol

import (
	"sync"

	"github.com/MarkSaravi/drone-go/modules/command"
	"github.com/MarkSaravi/drone-go/types"
)

type flightControl struct {
	imu              types.IMU
	motorsController types.MotorsController
	pid              types.PID
	logger           types.UdpLogger
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
	imu types.IMU,
	motorsController types.MotorsController,
	pid types.PID,
	logger types.UdpLogger,
) *flightControl {
	return &flightControl{
		imu:              imu,
		motorsController: motorsController,
		pid:              pid,
		logger:           logger,
	}
}
