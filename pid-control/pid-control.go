package pidcontrol

import "github.com/MarkSaravi/drone-go/models"

type pidControl struct {
}

func NewPIDControl() *pidControl {
	return &pidControl{}
}

func (pid *pidControl) ApplyFlightCommands(flightCommands models.FlightCommands) {
}

func (pid *pidControl) ApplyRotations(rotations models.ImuRotations) {

}
