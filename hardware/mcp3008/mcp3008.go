package mcp3008

import (
	"periph.io/x/conn/v3/spi"
)

const (
	DIGITAL_MAX_VALUE uint16 = 1024
)

type mcp3008dev struct {
	spiConn     spi.Conn
	channel     int
	value       uint16
	midValue    uint16
	minValue    uint16
	maxValue    uint16
	outputRange int
}

func NewMCP3008(spiConn spi.Conn, channel int, midValue uint16, outputRange int) *mcp3008dev {
	inputRange := midValue
	upperRange := DIGITAL_MAX_VALUE - midValue
	if upperRange < inputRange {
		inputRange = upperRange
	}
	return &mcp3008dev{
		spiConn:     spiConn,
		channel:     channel,
		value:       0,
		midValue:    midValue,
		minValue:    midValue - inputRange,
		maxValue:    midValue + inputRange,
		outputRange: outputRange,
	}
}

func (dev *mcp3008dev) Read() uint16 {
	ch := byte(dev.channel)
	if ch > 7 {
		ch = 0
	}
	w := []byte{0b00000001, 0b10000000 | (ch << 4), 0b00000000}
	r := []byte{0, 0, 0}
	err := dev.spiConn.Tx(w, r)
	if err != nil {
		return dev.value
	}
	v := uint16(r[2]) | (uint16(r[1]) << 8 & 0b0000001100000000)
	return dev.ConvertToOutputRange(v)
}

func (dev *mcp3008dev) ConvertToOutputRange(v uint16) uint16 {
	if v < dev.minValue {
		v = dev.minValue
	}
	if v > dev.maxValue {
		v = dev.maxValue
	}
	x := float32(v)
	diff := x - float32(dev.midValue)

	y := diff*float32(dev.outputRange)/float32(dev.maxValue-dev.minValue) + float32(dev.outputRange)/2
	// fmt.Printf("%8.2f    %8.2f   %8.2f  %8.2f\n", diff, float32(dev.outputRange), float32(dev.maxValue-dev.minValue), y)
	return uint16(y)
}
