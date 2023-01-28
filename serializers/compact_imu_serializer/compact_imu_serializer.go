/*
IMU compacted data packet serialization format

Time Stamp: duration from the start of sending packets (micro seconds)
Format Code: 16
Roll, Pitch, Yaw range: -360..360 with one decimal accuracy
Roll, Pitch, Yaw comact data range: -3600..3600
Decimal Precision: 1 digit
Roll, Pitch, Yaw math conversion: round(original value * 10)
Roll, Pitch, Yaw to byte conversion: LittleEndian
Time Interval range: Micro Seconds (max 60ms)

Packet Header (12 bytes 0..11):
┌────────────┬────────────┬─────────────────┬───────────────┬────────────────┬───────────────────────┬───────────────────────┐
│ packet len │ time stamp │ format code: 16 │ time interval │ number of data │ (Roll, Pitch, Yaw)[0] │ (Roll, Pitch, Yaw)[n] │
│ bytes 0..1 │ bytes 2..5 │ bytes 6..7      │ bytes 8..9    │ bytes 10..11   │ bytes 12..17          │ bytes 12+6*n..17+6*n  │
└────────────┴────────────┴─────────────────┴───────────────┴────────────────┴───────────────────────┴───────────────────────┘

*/

package compactimuserializer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
)

const DIGIT_FACTOR = 10
const FORMAT_CODE = 16
const HEADER_LEN = 12
const DATA_LEN = 6

type CompactSerialiserConfig struct {
	DataPerPacket int
	Interval      time.Duration
}

type compactSerialiser struct {
	config     CompactSerialiserConfig
	buffer     *bytes.Buffer
	packetSize int
}

func NewCompactSerialiser(config CompactSerialiserConfig) *compactSerialiser {
	packetSize := HEADER_LEN + config.DataPerPacket*DATA_LEN
	buffer := new(bytes.Buffer)
	buffer.Grow(packetSize)
	fmt.Println("PACKET SIZE: ", packetSize)

	return &compactSerialiser{
		config:     config,
		packetSize: packetSize,
		buffer:     buffer,
	}
}

func (c *compactSerialiser) Send(r imu.Rotations, timeStamp time.Duration) bool {
	if c.buffer.Len() == 0 {
		c.buffer.Reset()
		c.encodeHeader(timeStamp)
	}
	if c.buffer.Len() < c.packetSize {
		c.encodeRotations(r)
	}
	return c.buffer.Len() >= c.packetSize
}

func (c *compactSerialiser) Read() []byte {
	b := c.buffer.Bytes()
	c.buffer.Reset()
	return b
}

func compactAngle(value float64) int16 {
	return int16(int(math.Round(value * DIGIT_FACTOR)))
}

func (c *compactSerialiser) encodeHeader(timeStamp time.Duration) {
	binary.Write(c.buffer, binary.LittleEndian, uint16(c.packetSize))
	binary.Write(c.buffer, binary.LittleEndian, uint32(timeStamp.Microseconds()))
	binary.Write(c.buffer, binary.LittleEndian, uint16(FORMAT_CODE))
	binary.Write(c.buffer, binary.LittleEndian, uint16(c.config.Interval.Microseconds()))
	binary.Write(c.buffer, binary.LittleEndian, uint16(c.config.DataPerPacket))
}

func (c *compactSerialiser) encodeRotations(r imu.Rotations) {
	binary.Write(c.buffer, binary.LittleEndian, compactAngle(r.Roll))
	binary.Write(c.buffer, binary.LittleEndian, compactAngle(r.Pitch))
	binary.Write(c.buffer, binary.LittleEndian, compactAngle(r.Yaw))
}
