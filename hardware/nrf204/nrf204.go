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
	ON  bool = true
	OFF bool = false
)

const (
	CLEAR_CONFIG      = 0x00
	EN_CRC       byte = 0b00001000
	CRCO         byte = 0b00000100
)

const (
	NRF_CONFIG   byte = 0x00
	EN_AA        byte = 0x01
	EN_RXADDR    byte = 0x02
	SETUP_AW     byte = 0x03
	SETUP_RETR   byte = 0x04
	RF_CH        byte = 0x5
	RF_SETUP     byte = 0x06
	NRF_STATUS   byte = 0x07
	RX_PW_P0     byte = 0x11
	RX_ADDR_P0   byte = 0x0A
	TX_ADDR      byte = 0x10
	DYNPD        byte = 0x1C
	FEATURE      byte = 0x1D
	R_RX_PAYLOAD byte = 0x61
	FLUSH_TX     byte = 0xE1
	FLUSH_RX     byte = 0xE2
)

const (
	RF_POWER_MINUS_18dBm byte = iota // -18 dBm
	RF_POWER_MINUS_12dBm             //-12 dBm
	RF_POWER_MINUS_6dBm              // -6 dBm
	RF_POWER_0dBm                    // 0 dBm
)

const (
	DATA_RATE_1MBPS byte = iota
	DATA_RATE_2MBPS
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
	radio.setPower(OFF)
	radio.setRetries(5, 15)
	radio.setDataRate(DATA_RATE_1MBPS)
	radio.writeRegisterByte(EN_AA, 0x3F)
	radio.writeRegisterByte(EN_RXADDR, 3)
	radio.setPayloadSize()
	radio.setAddressWidth()
	radio.setChannel()
	radio.setCRCEncodingScheme()
	radio.enableCRC()
}

func (radio *nrf204l01) setRetries(delay byte, numRetransmit byte) {
	nr := numRetransmit
	if nr > 15 {
		nr = 15
	}
	d := delay
	if d > 15 {
		d = 5
	}
	setup := nr | (d >> 4)
	radio.writeRegisterByte(SETUP_RETR, setup)
}

func (radio *nrf204l01) setDataRate(dataRate byte) {
	dr := dataRate
	if dr > DATA_RATE_2MBPS {
		dr = DATA_RATE_2MBPS
	}
	setup, _ := radio.readRegisterByte(RF_SETUP)

	radio.writeRegisterByte(RF_SETUP, (setup&0b11110111)|(dr<<3))
}

func (radio *nrf204l01) saveAddress(rxAddress string) {
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

func (radio *nrf204l01) SetAddress(rxAddress string) {
	radio.saveAddress(rxAddress)
	radio.writeRegister(RX_ADDR_P0, []byte{0, 0, 0, 0, 0})
	radio.writeRegister(TX_ADDR, []byte{0, 0, 0, 0, 0})
	radio.writeRegister(RX_ADDR_P0, radio.address)
	radio.writeRegister(TX_ADDR, radio.address)
}

func (radio *nrf204l01) SetPALevel(rfPower byte) {
	setup, _ := radio.readRegisterByte(RF_SETUP)
	setup = (setup & 0b11110001) | (rfPower << 1)
	radio.writeRegisterByte(RF_SETUP, setup)
}

func (radio *nrf204l01) startTransmitting() {
	radio.ce.Out(gpio.Low)
	time.Sleep(10 * time.Millisecond)
	radio.flushTx()
	radio.writeRegisterByte(NRF_CONFIG, 14)
	radio.writeRegisterByte(EN_RXADDR, 3)
}

func (radio *nrf204l01) StartListening() {
	radio.setPower(ON)
	radio.clearStatus()
	radio.setRx(ON)
	radio.flushRx()
	radio.flushTx()
	radio.ce.Out(gpio.High)
}

func (radio *nrf204l01) IsAvailable(pipeNum byte) bool {
	// get implied RX FIFO empty flag from status byte
	status := radio.getStatus()
	// fmt.Println("Status: ", status)
	pipe := (status >> 1) & 0x07
	dr := status & 0b01000000
	return pipe <= 5 && dr == 64
}

func (radio *nrf204l01) getStatus() byte {
	status, _ := radio.writeRegisterByte(0xFF, 0xFF)
	return status[0]
}

func (radio *nrf204l01) ReadPayload() []byte {
	payload, _ := utils.ReadSPI(R_RX_PAYLOAD, 32, radio.conn)
	radio.writeRegisterByte(NRF_STATUS, 64)
	return payload
}

func (radio *nrf204l01) configOnOff(on bool, bitmask byte) {
	config, _ := radio.readRegisterByte(NRF_CONFIG)
	if on {
		radio.writeRegisterByte(NRF_CONFIG, config|bitmask)
	} else {
		radio.writeRegisterByte(NRF_CONFIG, config&(^bitmask))
	}
}

func (radio *nrf204l01) setPower(isOn bool) {
	radio.configOnOff(isOn, 0b00000010)
}

func (radio *nrf204l01) setRx(isOn bool) {
	radio.configOnOff(isOn, 0b00000001)
}

func (radio *nrf204l01) setCRCEncodingScheme() {
	radio.configOnOff(ON, 0b00000100)
}

func (radio *nrf204l01) enableCRC() {
	radio.configOnOff(ON, 0b00001000)
}

func (radio *nrf204l01) clearStatus() {
	status, _ := radio.readRegisterByte(NRF_STATUS)
	radio.writeRegisterByte(NRF_STATUS, status|0b01110000)
}

func (radio *nrf204l01) readRegister(address byte, len int) ([]byte, error) {
	b, err := utils.ReadSPI(address&R_REGISTER, len, radio.conn)
	return b, err
}

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
	utils.WriteSPI(FLUSH_RX, []byte{0xFF}, radio.conn)
}

func (radio *nrf204l01) flushTx() {
	utils.WriteSPI(FLUSH_TX, []byte{0xFF}, radio.conn)
}

func initPin(pinName string) gpio.PinIO {
	pin := gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	return pin
}
