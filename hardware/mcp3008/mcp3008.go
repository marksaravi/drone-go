package mcp3008

import (
	"fmt"
	"periph.io/x/conn/v3/spi"
)

type mcp3008dev struct {
	spiConn     spi.Conn
	channel     int
	value       uint16
}

func NewMCP3008(spiConn spi.Conn, channel int) *mcp3008dev {
	return &mcp3008dev{
		spiConn:     spiConn,
		channel:     channel,
		value:       0,
	}
}

func (dev *mcp3008dev) Read() uint16 {
	ch := byte(dev.channel)
	w := []byte{0b00000001 , 0b00110000, 0b00000000}
	r := []byte{0, 0, 0}
	err := dev.spiConn.Tx(w, r)
	if err != nil {
		return dev.value
	}
	
	l:=uint16(r[2])
	h:=uint16(r[1])<<8
	v:= l | h

	if (ch == 3) {
		fmt.Printf("%b %b %b    %b %b %b %d \n", w[0], w[1], w[2], r[0], r[1], r[2], v)
	}
	
	return 0
}
