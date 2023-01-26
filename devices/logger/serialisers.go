package logger

import (
	"github.com/marksaravi/drone-go/devices/imu"
)

const (
	ROLL_PITCH_YAW_TYPE_16_SERIALISATION = 16
)

func (l *udpLogger) rotationsToType16Serialisation(rotation imu.Rotations) bool {
	if l.buffer.Len() > 0 {

	} else {
		// now := time.Now().UnixMilli()
		// l.buffer.Write()
	}
	return l.buffer.Len() > 100
}

func (l *udpLogger) writeType16Header(rotation imu.Rotations) {
	// l.buffer.Write([]byte(time.Now().UnixMilli()))
}
