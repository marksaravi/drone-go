package compactimuserializer

import (
	"testing"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
)

func TestCompactSerialiser(t *testing.T) {
	timeStamp := time.Duration(0)
	interval := time.Duration(time.Millisecond)
	const dataPerPacket = 20
	const packetSize = HEADER_LEN + dataPerPacket*DATA_LEN
	config := CompactSerialiserConfig{
		DataPerPacket: dataPerPacket,
		Interval:      interval,
	}
	t.Run("test capacity", func(t *testing.T) {
		serialiser := NewCompactSerialiser(config)
		if serialiser.buffer.Cap() < packetSize {
			t.Errorf("buffer capacity must be >= %d but is %d", packetSize, serialiser.buffer.Cap())
		}
	})
	t.Run("test empty len", func(t *testing.T) {
		serialiser := NewCompactSerialiser(config)
		if serialiser.buffer.Len() > 0 {
			t.Errorf("buffer len must be 0 but is %d", serialiser.buffer.Len())
		}
	})

	t.Run("test header len", func(t *testing.T) {
		serialiser := NewCompactSerialiser(config)
		serialiser.encodeHeader(timeStamp)
		if serialiser.buffer.Len() != HEADER_LEN {
			t.Errorf("buffer len must be %d but is %d", HEADER_LEN, serialiser.buffer.Len())
		}
	})

	t.Run("test after adding n len", func(t *testing.T) {
		serialiser := NewCompactSerialiser(config)
		expectedLen := HEADER_LEN
		for i := 0; i < serialiser.config.DataPerPacket-1; i++ {
			ok := serialiser.Send(imu.Rotations{}, timeStamp)
			expectedLen += DATA_LEN
			if serialiser.buffer.Len() != expectedLen {
				t.Errorf("buffer len must be %d but is %d", expectedLen, serialiser.buffer.Len())
			}
			if ok {
				t.Error("serialiser must not be full yet")
			}
			timeStamp += interval
		}
		ok := serialiser.Send(imu.Rotations{}, timeStamp)
		if !ok {
			t.Error("serialiser must be full")
		}
	})
	t.Run("test after reset", func(t *testing.T) {
		serialiser := NewCompactSerialiser(config)
		for i := 0; i < serialiser.config.DataPerPacket; i++ {
			serialiser.Send(imu.Rotations{}, timeStamp)
		}
		buffer := serialiser.Read()
		if len(buffer) != packetSize {
			t.Errorf("buffer size must be %d but is %d", packetSize, len(buffer))
		}
		if serialiser.buffer.Len() != 0 {
			t.Error("buffer must be empty after reset")
		}
	})

}
