package devices

import (
	"time"
)

type imudevice struct {
	readTime        time.Time
	readingInterval time.Duration
}

func NewIMU(readingInterval time.Duration) *imudevice {
	return &imudevice{
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
