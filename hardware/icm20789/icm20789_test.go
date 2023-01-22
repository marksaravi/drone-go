package icm20789_test

import (
	"testing"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20789"
	"github.com/marksaravi/drone-go/types"
)

var gerr error
var gdata types.IMUMems6DOFRawData

func BenchmarkReadSpeed(b *testing.B) {
	hardware.HostInitialize()

	imu := icm20789.NewICM20789(icm20789.ICM20789Configs{
		AccelerometerFullScale: "2g",
		GyroscopeFullScale:     "250dps",
	})

	var err error
	var data types.IMUMems6DOFRawData

	for i := 0; i < b.N; i++ {
		data, err = imu.Read()
	}
	gdata = data
	gerr = err
}
