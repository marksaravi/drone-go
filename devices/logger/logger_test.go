package logger

import (
	"encoding/binary"
	"testing"
)

func TestSerialiseFloat64(t *testing.T) {
	const WANT float64 = 173.9
	l := NewUDPLogger(0, 0)
	l.serialiseFloat64(WANT)
	var i int16
	binary.Read(l.buffer, binary.LittleEndian, &i)
	got := float64(i) / DIGIT_FACTOR
	if got != WANT {
		t.Errorf("serialise error, want %f, got %f", WANT, got)
	}
}
