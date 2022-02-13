package nrf204

import (
	"log"

	"github.com/marksaravi/drone-go/constants"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/spi"
)

const (
	ADDRESS_CONFIG      byte = 0x0
	ADDRESS_EN_AA       byte = 0x1
	ADDRESS_EN_RXADDR   byte = 0x2
	ADDRESS_SETUP_AW    byte = 0x3
	ADDRESS_SETUP_RETR  byte = 0x4
	ADDRESS_RF_CH       byte = 0x5
	ADDRESS_RF_SETUP    byte = 0x6
	ADDRESS_STATUS      byte = 0x7
	ADDRESS_OBSERVE_TX  byte = 0x8
	ADDRESS_CD          byte = 0x9
	ADDRESS_RX_ADDR_P0  byte = 0xA
	ADDRESS_RX_ADDR_P1  byte = 0xB
	ADDRESS_RX_ADDR_P2  byte = 0xC
	ADDRESS_RX_ADDR_P3  byte = 0xD
	ADDRESS_RX_ADDR_P4  byte = 0xE
	ADDRESS_RX_ADDR_P5  byte = 0xF
	ADDRESS_TX_ADDR     byte = 0x10
	ADDRESS_RX_PW_P0    byte = 0x11
	ADDRESS_RX_PW_P1    byte = 0x12
	ADDRESS_RX_PW_P2    byte = 0x13
	ADDRESS_RX_PW_P3    byte = 0x14
	ADDRESS_RX_PW_P4    byte = 0x15
	ADDRESS_RX_PW_P5    byte = 0x16
	ADDRESS_FIFO_STATUS byte = 0x17

	ADDRESS_W_TX_PAYLOAD byte = 0xA0
	ADDRESS_R_RX_PAYLOAD byte = 0x61

	ADDRESS_FLUSH_RX byte = 0xE1
	ADDRESS_FLUSH_TX byte = 0xE2
)

const (
	DEFAULT_CONFIG     byte = 0b00001000
	DEFAULT_EN_AA      byte = 0b00111111
	DEFAULT_EN_RXADDR  byte = 0b00000001
	DEFAULT_SETUP_AW   byte = 0b00000011
	DEFAULT_SETUP_RETR byte = 0b00000011
	DEFAULT_RF_CH      byte = 0b01001100
	DEFAULT_RF_SETUP   byte = 0b00001111
	DEFAULT_STATUS     byte = 0b01110000
	DEFAULT_RX_PW_P0   byte = constants.RADIO_PAYLOAD_SIZE
	DEFAULT_RX_PW_P1   byte = 0b00000000
	DEFAULT_RX_PW_P2   byte = 0b00000000
	DEFAULT_RX_PW_P3   byte = 0b00000000
	DEFAULT_RX_PW_P4   byte = 0b00000000
	DEFAULT_RX_PW_P5   byte = 0b00000000
)

type nrf204l01 struct {
	ce        gpio.PinOut
	address   []byte
	conn      spi.Conn
	status    byte
	registers map[byte]byte
}

func writeSPI(address byte, data []byte, conn spi.Conn) ([]byte, error) {
	datalen := len(data)
	w := make([]byte, datalen+1)
	r := make([]byte, datalen+1)
	w[0] = address | 0x20 // adding write bit

	for i := 0; i < datalen; i++ {
		w[i+1] = data[i]
	}

	err := conn.Tx(w, r)
	return r, err
}

func readSPI(address byte, datalen int, conn spi.Conn) ([]byte, error) {
	w := make([]byte, datalen+1)
	r := make([]byte, datalen+1)
	w[0] = address
	err := conn.Tx(w, r)
	return r[1:], err
}

func initPin(pinName string) gpio.PinIO {
	pin := gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	return pin
}
