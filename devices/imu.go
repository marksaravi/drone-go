package devices

import (
	"time"

	"github.com/MarkSaravi/drone-go/models"
)

type imuMems interface {
	Read() (acc models.XYZ, gyro models.XYZ)
}

type imudevice struct {
	imuMems         imuMems
	readTime        time.Time
	readingInterval time.Duration
}

func NewIMU(imuMems imuMems, readingInterval time.Duration) *imudevice {
	return &imudevice{
		imuMems:         imuMems,
		readTime:        time.Now(),
		readingInterval: readingInterval,
	}
}

func (imu *imudevice) Read() (canRead bool) {
	if time.Since(imu.readTime) < imu.readingInterval {
		canRead = false
		return
	}
	imu.readTime = time.Now()
	canRead = true
	return
}
