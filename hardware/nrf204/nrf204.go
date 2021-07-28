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
	EN_AA      byte = 0x01
	EN_RXADDR  byte = 0x02
	SETUP_AW   byte = 0x03
	SETUP_RETR byte = 0x04
	RF_CH      byte = 0x5
	RF_SETUP   byte = 0x06
	NRF_STATUS byte = 0x07
	RX_PW_P0   byte = 0x11
	RX_ADDR_P0 byte = 0x0A
	DYNPD      byte = 0x1C
	FEATURE    byte = 0x1D
	FLUSH_TX   byte = 0xE1
	FLUSH_RX   byte = 0xE2
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
	radio.initRadio()
}

func (radio *nrf204l01) initRadio() {
	radio.setRetries()
	radio.setDataRate()
	radio.writeRegisterByte(DYNPD, 0)
	radio.writeRegisterByte(EN_AA, 0x3F)
	radio.writeRegisterByte(EN_RXADDR, 3)
	radio.setPayloadSize()
	radio.setAddressWidth()
	radio.setChannel()
	radio.writeRegisterByte(NRF_STATUS, 112)
}

func (radio *nrf204l01) setRetries() {
	radio.writeRegisterByte(SETUP_RETR, 95)

	// radio.flushRx()
	// radio.flushTx()

	// // Clear CONFIG register:
	// //      Reflect all IRQ events on IRQ pin
	// //      Enable PTX
	// //      Power Up
	// //      16-bit CRC (CRC required by auto-ack)
	// // Do not write CE high so radio will remain in standby I mode
	// // PTX should use only 22uA of power
	// radio.writeRegisterByte(NRF_CONFIG, 12)
	// configReg, _ := radio.readRegisterByte(NRF_CONFIG)
	// fmt.Println("config reg: ", configReg)

	// radio.powerUp()
}

func (radio *nrf204l01) setDataRate() {
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
		radio.address[i] = b[i]
	}
	fmt.Println(radio.address)
}

func (radio *nrf204l01) SetPALevel(level byte, lnaEnable byte) {
	reg, _ := radio.readRegisterByte(RF_SETUP)
	setup := reg & 0xF8

	if level > 3 { // If invalid level, go to max PA
		level = (RF24_PA_MAX << 1) + lnaEnable // +1 to support the SI24R1 chip extra bit
	} else {
		level = (level << 1) + lnaEnable // Else set level as requested
	}

	nSetup := setup | level
	radio.writeRegisterByte(RF_SETUP, nSetup)
}

func (radio *nrf204l01) StartListening() {
	radio.powerUp()
	radio.writeRegisterByte(NRF_CONFIG, 15)
	radio.writeRegisterByte(NRF_STATUS, 112)
	radio.ce.Out(gpio.High)

	// Restore the pipe0 address, if exists
	radio.writeRegister(RX_ADDR_P0, radio.address)
}

func (radio *nrf204l01) IsAvailable(pipeNum byte) bool {
	// get implied RX FIFO empty flag from status byte
	status := radio.getStatus()
	pipe := (status >> 1) & 0x77
	fmt.Println(pipe)
	return pipe <= 5
}

func (radio *nrf204l01) getStatus() byte {
	// status, _ := radio.readRegisterByte(NRF_STATUS)
	r, _ := radio.writeRegisterByte(0xFF, 0xFF)
	return r[0]
}

func (radio *nrf204l01) powerUp() {
	radio.writeRegisterByte(NRF_CONFIG, 12)
}

// func (radio *nrf204l01) readRegister(address byte, datalen int) ([]byte, error) {
// 	return utils.ReadSPI(address&R_REGISTER, datalen, radio.conn)
// }

func (radio *nrf204l01) readRegisterByte(address byte) (byte, error) {
	b, err := utils.ReadSPI(address&R_REGISTER, 1, radio.conn)
	return b[0], err
}

func (radio *nrf204l01) writeRegister(address byte, data []byte) ([]byte, error) {
	return utils.WriteSPI((address&R_REGISTER)|W_REGISTER, data, radio.conn)
}

func (radio *nrf204l01) writeRegisterByte(address byte, data byte) ([]byte, error) {
	return utils.WriteSPI((address&R_REGISTER)|W_REGISTER, []byte{data}, radio.conn)
}

func (radio *nrf204l01) setPayloadSize() {
	var i byte
	for i = 0; i < 6; i++ {
		radio.writeRegisterByte(RX_PW_P0+i, 32)
	}
}

func (radio *nrf204l01) setAddressWidth() {
	radio.writeRegisterByte(SETUP_AW, 3)
}

func (radio *nrf204l01) setChannel() {
	radio.writeRegisterByte(RF_CH, 76)
}

func (radio *nrf204l01) flushRx() {
	radio.writeRegisterByte(FLUSH_RX, 0xFF)
}

func (radio *nrf204l01) flushTx() {
	radio.writeRegisterByte(FLUSH_TX, 0xFF)
}

func (radio *nrf204l01) Read(len byte) {

	// Fetch the payload
	// radio.readPayload(buf, len)

	//Clear the only applicable interrupt flags
	radio.writeRegisterByte(NRF_STATUS, 72)

}

func initPin(pinName string) gpio.PinIO {
	pin := gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	return pin
}
