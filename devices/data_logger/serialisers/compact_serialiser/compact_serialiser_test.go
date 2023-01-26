package compactserialiser

import "testing"

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
		serialiser.SetHeader()
		if serialiser.buffer.Len() != HEADER_SIZE {
			t.Errorf("buffer len must be %d but is %d", HEADER_SIZE, serialiser.buffer.Len())
		}
	})
}
