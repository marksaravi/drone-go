package mcp3008

type spi interface{}

const (
	DIGITAL_MAX_VALUE uint16 = 1024
)

type mcp3008dev struct {
	spiConn spi.Conn
	channel int
	value   uint16
}

func NewMCP3008(spiConn spi.Conn, channel int) *mcp3008dev {
	return &mcp3008dev{
		spiConn: spiConn,
		channel: channel,
		value:   0,
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
	dev.value = uint16(r[2]) | (uint16(r[1]) << 8 & 0b0000001100000000)
	return dev.value
}
