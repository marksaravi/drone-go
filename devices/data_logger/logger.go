package datalogger

import (
	"bytes"
	"sync"

	"github.com/marksaravi/drone-go/devices/imu"
)

const DIGIT_FACTOR = 10
const DATA_SIZE = 6

type dataLogger struct {
	buffer     *bytes.Buffer
	dt         int16
	packetSize int16
	wg         *sync.WaitGroup
}

func NewDataLogger(wg *sync.WaitGroup, numberOfData, dt int16) *dataLogger {
	return &dataLogger{
		buffer:     new(bytes.Buffer),
		packetSize: numberOfData * DATA_SIZE,
		dt:         dt,
		wg:         wg,
	}
}

func (l *dataLogger) SendRotation(rotations imu.Rotations) {
}
