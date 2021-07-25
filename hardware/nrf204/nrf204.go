package nrf204

import (
	"fmt"
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/spi"
)

const (
	R_REGISTER byte = 0x1F
	W_REGISTER byte = 0x20
)
const (
	CLEAR_CONFIG      = 0x00
	EN_CRC       byte = 0b00001000
	CRCO         byte = 0b00000100
)

const (
	NRF_CONFIG byte = 0x00
	SETUP_RETR byte = 0x04
	RF_SETUP   byte = 0x06
	NRF_STATUS byte = 0x07
	RX_ADDRESS byte = 0x0A
)
const (
	RF24_PA_MIN byte = iota
	RF24_PA_LOW
	RF24_PA_HIGH
	RF24_PA_MAX
	RF24_PA_ERROR
)

const (
	RF24_1MBPS byte = iota
	RF24_2MBPS
	RF24_250KBPS
)

const addressSize int = 5

type RadioMode int

const (
	Receiver RadioMode = iota
	Transmitter
)

type nrf204l01 struct {
	ce      gpio.PinOut
	address []byte
	conn    spi.Conn
}

func CreateNRF204(config types.RadioLinkConfig, conn spi.Conn) *nrf204l01 {
	return &nrf204l01{
		ce:      initPin(config.GPIO.CE),
		address: make([]byte, addressSize),
		conn:    conn,
	}
}

func (radio *nrf204l01) Init() {
	radio.ce.Out(gpio.Low)
	time.Sleep(time.Millisecond * 10)
	radio.setRetries(5, 15)
}

func (radio *nrf204l01) setRetries(delay byte, count byte) {
	retry := utils.Min(delay, 15)<<4 | utils.Min(count, 15)
	radio.writeRegisterByte(SETUP_RETR, retry)
	radio.serDataSize()
}

func (radio *nrf204l01) serDataSize() {
	radio.writeRegisterByte(RF_SETUP, 1)
}

func (radio *nrf204l01) OpenReadingPipe(rxAddress string) {
	// This implementation only supports the single pipe for now
	b := []byte(rxAddress)
	lenb := len(b)
	if lenb != len(radio.address) {
		log.Fatal("Rx Address for Radio link is incorrect")
	}
	for i := 0; i < lenb; i++ {
		radio.address[i] = b[i] - 48
	}
	fmt.Println(radio.address)
}

func (radio *nrf204l01) SetPALevel(level byte, lnaEnable byte) {
	// r, err := radio.readRegister(RF_SETUP, 1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// setup := r[0] & 0xF8
	// l := level
	// if l > 3 {
	// 	l = (RF24_PA_MAX << 1) + lnaEnable
	// } else {
	// 	l = (l << 1) + lnaEnable
	// }
	// data := []byte{setup | l}
	// radio.writeRegister(RF_SETUP, data)
}

func (radio *nrf204l01) StartListening() {
	radio.powerUp()
	// var config_reg byte = 0
	// var status_reg byte = 0

	// r.writeRegister(NRF_CONFIG, config_reg)
	// r.writeRegister(NRF_STATUS, status_reg)
	// r.writeRegister(RX_ADDRESS, r.address)
}

func (radio *nrf204l01) IsAvailable() bool {
	return false
}

func (radio *nrf204l01) Read() []byte {
	return []byte{0, 0, 0, 0}
}

func (radio *nrf204l01) powerUp() {
}

func (radio *nrf204l01) readRegister(address byte, datalen int) ([]byte, error) {
	return utils.ReadSPI(address&R_REGISTER, datalen, radio.conn)
}

func (radio *nrf204l01) readRegisterByte(address byte) (byte, error) {
	b, err := utils.ReadSPI(address&R_REGISTER, 1, radio.conn)
	return b[0], err
}

func (radio *nrf204l01) writeRegister(address byte, data []byte) error {
	return utils.WriteSPI((address&R_REGISTER)|W_REGISTER, data, radio.conn)
}

func (radio *nrf204l01) writeRegisterByte(address byte, data byte) error {
	return utils.WriteSPI((address&R_REGISTER)|W_REGISTER, []byte{data}, radio.conn)
}

func initPin(pinName string) gpio.PinIO {
	pin := gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	return pin
}
