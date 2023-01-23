package logger

import (
	"bytes"
	"encoding/binary"
	"sync"

	"github.com/marksaravi/drone-go/devices/imu"
)

const DIGIT_FACTOR = 10
const DATA_SIZE = 6

type udpLogger struct {
	buffer     *bytes.Buffer
	dt         int16
	packetSize int16
	wg         *sync.WaitGroup
}

func NewUDPLogger(wg *sync.WaitGroup, numberOfData, dt int16) *udpLogger {
	return &udpLogger{
		buffer:     new(bytes.Buffer),
		packetSize: numberOfData * DATA_SIZE,
		dt:         dt,
		wg:         wg,
	}
}

func (l *udpLogger) SendRotation(rotations imu.Rotations) {
	if l.buffer.Len() == 0 {
		l.serialiseInt16(l.packetSize)
		l.serialiseInt16(l.dt)
	}
	l.serialiseFloat64(rotations.Roll)
	if l.buffer.Len() == int(l.packetSize) {
		l.transmit()
	}
}

func (l *udpLogger) serialiseFloat64(f float64) {
	l.serialiseInt16(int16(int(f * DIGIT_FACTOR)))
}

func (l *udpLogger) serialiseInt16(value int16) {
	binary.Write(l.buffer, binary.LittleEndian, value)
}

func (l *udpLogger) transmit() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
	}()
}
