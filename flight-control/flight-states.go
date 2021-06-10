package flightcontrol

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/types"
)

type FlightStates struct {
	Config       types.FlightConfig
	imuRotations types.ImuRotations
}

func (fs *FlightStates) Update(imuRotations types.ImuRotations) {
	fs.imuRotations = imuRotations
}

func (fs *FlightStates) ImuDataToJson() string {
	return fmt.Sprintf(`{"a":{"r":%0.2f,"p":%0.2f,"y":%0.2f},"g":{"r":%0.2f,"p":%0.2f,"y":%0.2f},"r":{"r":%0.2f,"p":%0.2f,"y":%0.2f},"t":%d,"dt":%d}`,
		fs.imuRotations.Accelerometer.Roll,
		fs.imuRotations.Accelerometer.Pitch,
		fs.imuRotations.Accelerometer.Yaw,
		fs.imuRotations.Gyroscope.Roll,
		fs.imuRotations.Gyroscope.Pitch,
		fs.imuRotations.Gyroscope.Yaw,
		fs.imuRotations.Rotations.Roll,
		fs.imuRotations.Rotations.Pitch,
		fs.imuRotations.Rotations.Yaw,
		fs.imuRotations.ReadTime,
		fs.imuRotations.ReadInterval,
	)
}
