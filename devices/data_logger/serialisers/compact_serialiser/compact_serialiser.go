/*
Data Package format
Packet Format:

bytes: 0          1            2           3            4          5           6           ...n
┌────────────┬────────────┬───────────┬───────────┬───────────┬───────────┬─────────────┬─────────┐
│ packet len │ packet len │ packet id │ packet id │ packet id │ packet id │ format code │ data... │
│ high byte  │ low byte   │ byte #3   │ byte #2   │ byte #1   │ byte #0   │             │         │
└────────────┴────────────┴───────────┴───────────┴───────────┴───────────┴─────────────┴─────────┘


Type 16, Simple Roll, Pitch, Yaw data serialisation
packet len: 10 + 6 * (number of data)
Packet ID is in ms from start time
Format Code: 16
Roll, Pitch, Yaw range: -360..360
Roll, Pitch, Yaw data type: int16
Decimal Precision: 1 digit
Roll, Pitch, Yaw math conversion: round(original value * 10)
Roll, Pitch, Yaw to byte conversion: LittleEndian
Time Interval range: 1..255ms
Time Interval data type: byte

Packet Information (bytes 0..9):
bytes: 0          1            2           3            4          5        6   7         8                9             10               11                12           13
┌────────────┬────────────┬───────────┬───────────┬───────────┬───────────┬───┬──────┬───────────────┬──────────────┬───────────────┬───────────────┬────────────────┬────────────────┐
│ packet len │ packet len │ packet id │ packet id │ packet id │ packet id │ 16│ not  │ time interval │time interval │ time interval │ time interval │ number of data │ number of data │
│ high byte  │ low byte   │ byte #3   │ byte #2   │ byte #1   │ byte #0   │   │ used │ (us) byte #3  │ (us) byte #2 │ (us) byte #1  │ (us) byte #0  │ high byte      │ low byte       │
└────────────┴────────────┴───────────┴───────────┴───────────┴───────────┴───┴──────┴───────────────┴──────────────┴───────────────┴───────────────┴────────────────┴────────────────┘

Packet Data (bytes 14..9+6*(number of data)):
bytes: 13       14          15           16         17         18          ...                                                     14..9+6*(number of data)
┌───────────┬──────────┬───────────┬───────────┬───────────┬───────────┬───────────┬──────────┬───────────┬──────────┬───────────┬──────────┐
│ Roll      │ Roll     │ Pitch     │ Pitch     │ Yaw       │ Yaw       │ Roll      │ Roll     │ Pitch     │ Pitch    │ Yaw       │ Yaw      │
│ high byte │ low byte │ high byte │ low byte  │ high byte │ low byte  │ high byte │ low byte │ high byte │ low byte │ high byte │ low byte │
└───────────┴──────────┴───────────┴───────────┴───────────┴───────────┴───────────┴──────────┴───────────┴──────────┴───────────┴──────────┘

*/

package compactserialiser

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
const HEADER_SIZE = 14
const DATA_SIZE = 6

type CompactSerialiserConfig struct {
	DataPerPacket int
	IntervalMS    int
}

type compactSerialiser struct {
	config     CompactSerialiserConfig
	buffer     *bytes.Buffer
	packetSize int
	startTime  time.Time
}

func NewCompactSerialiser(config CompactSerialiserConfig) *compactSerialiser {
	packetSize := HEADER_SIZE + config.DataPerPacket*DATA_SIZE
	buffer := new(bytes.Buffer)
	buffer.Grow(packetSize)
	fmt.Println("PACKET SIZE: ", packetSize)

	return &compactSerialiser{
		config:     config,
		packetSize: packetSize,
		buffer:     buffer,
		startTime:  time.Now(),
	}
}

func (c *compactSerialiser) setHeader() {
	c.buffer.Reset()
	binary.Write(c.buffer, binary.LittleEndian, uint16(c.packetSize))
	t := time.Since(c.startTime).Milliseconds()
	binary.Write(c.buffer, binary.LittleEndian, int32(t))
	binary.Write(c.buffer, binary.LittleEndian, int16(FORMAT_CODE))
	binary.Write(c.buffer, binary.LittleEndian, int32(c.config.IntervalMS))
	binary.Write(c.buffer, binary.LittleEndian, uint16(c.config.DataPerPacket))
}

func (c *compactSerialiser) Send(r imu.Rotations) bool {
	if c.buffer.Len() == 0 {
		c.setHeader()
	}
	if c.buffer.Len() < c.packetSize {
		binary.Write(c.buffer, binary.LittleEndian, compactAngle(r.Roll))
		binary.Write(c.buffer, binary.LittleEndian, compactAngle(r.Pitch))
		binary.Write(c.buffer, binary.LittleEndian, compactAngle(r.Yaw))
	}
	return c.buffer.Len() >= c.packetSize
}

func compactAngle(value float64) int16 {
	return int16(int(math.Round(value * DIGIT_FACTOR)))
}
