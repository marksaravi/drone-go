package compactserialiser

import (
	"testing"

	"github.com/marksaravi/drone-go/devices/imu"
)

func TestCompactSerialiser(t *testing.T) {
	serialiser := NewCompactSerialiser(
		CompactSerialiserConfig{
			DataPerPacket: 20,
			IntervalMS:    25,
		},
	)
	packetSize := HEADER_SIZE + serialiser.config.DataPerPacket*DATA_SIZE
	t.Run("test capacity", func(t *testing.T) {
		if serialiser.buffer.Cap() < packetSize {
			t.Errorf("buffer capacity must be %d but is %d", 26, serialiser.buffer.Cap())
		}
	})
	t.Run("test empty len", func(t *testing.T) {
		if serialiser.buffer.Len() > 0 {
			t.Errorf("buffer len must be 0 but is %d", serialiser.buffer.Len())
		}
	})

	t.Run("test header len", func(t *testing.T) {
		serialiser.setHeader()
		if serialiser.buffer.Len() != HEADER_SIZE {
			t.Errorf("buffer len must be %d but is %d", HEADER_SIZE, serialiser.buffer.Len())
		}
	})

	t.Run("test after adding n len", func(t *testing.T) {
		const N = 3
		size := HEADER_SIZE
		for i := 0; i < N; i++ {
			serialiser.Send(imu.Rotations{})
			size += DATA_SIZE
			if serialiser.buffer.Len() != size {
				t.Errorf("buffer len must be %d but is %d", size, serialiser.buffer.Len())
			}
		}
	})
}
