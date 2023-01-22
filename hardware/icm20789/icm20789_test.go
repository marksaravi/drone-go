package icm20789_test

import (
	"testing"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20789"
)

var gerr error
var gdata icm20789.Mems6DOFData

func BenchmarkReadSpeed(b *testing.B) {
	hardware.HostInitialize()

	imu := icm20789.NewICM20789()

	var err error
	var data icm20789.Mems6DOFData

	for i := 0; i < b.N; i++ {
		data, err = imu.Read()
	}
	gdata = data
	gerr = err
}
