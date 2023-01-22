package icm20789_test

import (
	"testing"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mems"
	"github.com/marksaravi/drone-go/hardware/mems/icm20789"
)

var gerr error
var gdata mems.Mems6DOFData

func BenchmarkReadSpeed(b *testing.B) {
	hardware.HostInitialize()

	imu := icm20789.NewICM20789(icm20789.Configs{
		Accelerometer: icm20789.InertialDeviceConfigs{
			FullScale: "2g",
			Offsets: icm20789.Offsets{
				X: 0,
				Y: 0,
				Z: 0,
			},
		},
		Gyroscope: icm20789.InertialDeviceConfigs{
			FullScale: "250dps",
			Offsets: icm20789.Offsets{
				X: 0,
				Y: 0,
				Z: 0,
			},
		},
	})

	var err error
	var data mems.Mems6DOFData

	for i := 0; i < b.N; i++ {
		data, err = imu.Read()
	}
	gdata = data
	gerr = err
}
