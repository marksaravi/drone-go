package nrf204

import (
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/constants"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/models"
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
	W_TX_PAYLOAD byte = 0xA0
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

const ADDRESS_SIZE int = 5

type radioMode int

const (
	Receiver radioMode = iota
	Transmitter
)

type nrf204l01 struct {
	ce         gpio.PinOut
	address    []byte
	conn       spi.Conn
	powerDBm   byte
	isReceiver bool
}

func NewNRF204(
	spiBusNum int,
	spiChipSelect int,
	spiChipEnabledGPIO string,
	rxTxAddress string,
	powerDb string,
) *nrf204l01 {
	radioSPIConn := hardware.NewSPIConnection(
		spiBusNum,
		spiChipSelect,
	)

	radio := nrf204l01{
		address:    []byte(rxTxAddress),
		ce:         initPin(spiChipEnabledGPIO),
		conn:       radioSPIConn,
		powerDBm:   dbmStrToDBm(powerDb),
		isReceiver: true,
	}
	radio.init()
	radio.receiverOn()
	return &radio
}

func (radio *nrf204l01) ReceivePayload() (models.Payload, bool) {
	if !radio.isReceiver {
		radio.receiverOn()
	}
	if !radio.isDataAvailable() {
		return models.Payload{}, false
	}
	payload := models.Payload{}
	data, err := readSPI(R_RX_PAYLOAD, int(constants.RADIO_PAYLOAD_SIZE), radio.conn)
	copy(payload[:], data)
	radio.resetDR()
	if err != nil {
		return models.Payload{}, false
	}
	return payload, true
}

func (radio *nrf204l01) TransmitPayload(payload models.Payload) error {
	if len(payload) < int(constants.RADIO_PAYLOAD_SIZE) {
		return fmt.Errorf("payload size error %d", len(payload))
	}
	if radio.isReceiver {
		radio.transmitterOn()
	}
	// radio.ce.Out(gpio.Low)
	_, err := radio.writeRegister(TX_ADDR, radio.address)
	if err == nil {
		_, err = writeSPI(W_TX_PAYLOAD, payload[:], radio.conn)
	}
	radio.ce.Out(gpio.High)
	time.Sleep(time.Millisecond)
	radio.receiverOn()
	return err
}

func (radio *nrf204l01) transmitterOn() {
	radio.isReceiver = false
	radio.ce.Out(gpio.Low)
	radio.setRetries(1, 0)
	radio.clearStatus()
	radio.setRx(OFF)
	radio.flushRx()
	radio.flushTx()
	radio.setPower(ON)
	// radio.ce.Out(gpio.High)
}

func (radio *nrf204l01) receiverOn() {
	radio.isReceiver = true
	radio.ce.Out(gpio.Low)
	radio.setPower(ON)
	radio.clearStatus()
	radio.setRx(ON)
	radio.flushRx()
	radio.flushTx()
	radio.ce.Out(gpio.High)
}

func dbmStrToDBm(dbm string) byte {
	switch dbm {
	case "-18dbm":
		return RF_POWER_MINUS_18dBm
	case "-12dbm":
		return RF_POWER_MINUS_12dBm
	case "-6dbm":
		return RF_POWER_MINUS_6dBm
	case "0dbm":
		return RF_POWER_0dBm
	default:
		return RF_POWER_MINUS_18dBm
	}
}

func (radio *nrf204l01) init() {
	radio.ce.Out(gpio.Low)
	radio.setPower(OFF)
	radio.setRetries(5, 15)
	radio.setPALevel(radio.powerDBm)
	radio.setDataRate(DATA_RATE_1MBPS)
	// disabling auto acknowlegment
	radio.writeRegisterByte(EN_AA, 0)
	radio.writeRegisterByte(EN_RXADDR, 3)
	radio.setPayloadSize()
	radio.setAddressWidth()
	radio.setChannel()
	radio.setCRCEncodingScheme()
	radio.enableCRC()
	radio.setAddress()
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

func (radio *nrf204l01) setAddress() {
	radio.writeRegister(RX_ADDR_P0, []byte{0, 0, 0, 0, 0})
	radio.writeRegister(TX_ADDR, []byte{0, 0, 0, 0, 0})
	radio.writeRegister(RX_ADDR_P0, radio.address)
	radio.writeRegister(TX_ADDR, radio.address)
}

func (radio *nrf204l01) setPALevel(rfPower byte) {
	setup, _ := radio.readRegisterByte(RF_SETUP)
	setup = (setup & 0b11110001) | (rfPower << 1)
	radio.writeRegisterByte(RF_SETUP, setup)
}

func (radio *nrf204l01) isDataAvailable() bool {
	// get implied RX FIFO empty flag from status byte
	status := radio.getStatus()
	// fmt.Println("Status: ", status)
	pipe := (status >> 1) & 0x07
	dr := status & 0b01000000
	return pipe <= 5 && dr == 64
}

func (radio *nrf204l01) getStatus() byte {
	status, _ := radio.readRegisterByte(NRF_STATUS)
	return status
}

func writeSPI(address byte, data []byte, conn spi.Conn) ([]byte, error) {
	datalen := len(data)
	w := make([]byte, datalen+1)
	r := make([]byte, datalen+1)
	w[0] = address

	for i := 0; i < datalen; i++ {
		w[i+1] = data[i]
	}

	err := conn.Tx(w, r)
	return r, err
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

func (radio *nrf204l01) resetDR() {
	status, _ := radio.readRegisterByte(NRF_STATUS)
	radio.writeRegisterByte(NRF_STATUS, status|0b01000000)
}

func readSPI(address byte, datalen int, conn spi.Conn) ([]byte, error) {
	w := make([]byte, datalen+1)
	r := make([]byte, datalen+1)
	w[0] = address
	err := conn.Tx(w, r)
	return r[1:], err
}

func (radio *nrf204l01) readRegisterByte(address byte) (byte, error) {
	b, err := readSPI(address&R_REGISTER, 1, radio.conn)
	return b[0], err
}

func (radio *nrf204l01) writeRegister(address byte, data []byte) ([]byte, error) {
	return writeSPI((address&R_REGISTER)|W_REGISTER, data, radio.conn)
}

func (radio *nrf204l01) writeRegisterByte(address byte, data byte) ([]byte, error) {
	return writeSPI((address&R_REGISTER)|W_REGISTER, []byte{data}, radio.conn)
}

func (radio *nrf204l01) setPayloadSize() {
	var i byte
	for i = 0; i < 6; i++ {
		radio.writeRegisterByte(RX_PW_P0+i, constants.RADIO_PAYLOAD_SIZE)
	}
}

func (radio *nrf204l01) setAddressWidth() {
	radio.writeRegisterByte(SETUP_AW, 3)
}

func (radio *nrf204l01) setChannel() {
	radio.writeRegisterByte(RF_CH, 76)
}

func (radio *nrf204l01) flushRx() {
	writeSPI(FLUSH_RX, []byte{0xFF}, radio.conn)
}

func (radio *nrf204l01) flushTx() {
	writeSPI(FLUSH_TX, []byte{0xFF}, radio.conn)
}

func initPin(pinName string) gpio.PinIO {
	pin := gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	return pin
}
