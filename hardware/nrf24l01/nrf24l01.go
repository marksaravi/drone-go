package nrf24l01

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/constants"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/models"
)

type gpio interface{}
type gpioreg interface{}
type spi interface{}

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

	ADDRESS_FLUSH_TX byte = 0xE1
	ADDRESS_FLUSH_RX byte = 0xE2
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

type nrf24l01dev struct {
	ce        gpio.PinOut
	address   []byte
	conn      spi.Conn
	status    byte
	registers map[byte]byte
}

func NewNRF24L01EnhancedBurst(
	spiBusNum int,
	spiChipSelect int,
	spiChipEnabledGPIO string,
	rxTxAddress string,
) *nrf24l01dev {
	radioSPIConn := hardware.NewSPIConnection(
		spiBusNum,
		spiChipSelect,
	)

	tr := nrf24l01dev{
		address: []byte(rxTxAddress),
		ce:      initPin(spiChipEnabledGPIO),
		conn:    radioSPIConn,
		registers: map[byte]byte{
			ADDRESS_CONFIG:     DEFAULT_CONFIG,
			ADDRESS_EN_AA:      DEFAULT_EN_AA,
			ADDRESS_EN_RXADDR:  DEFAULT_EN_RXADDR,
			ADDRESS_SETUP_AW:   DEFAULT_SETUP_AW,
			ADDRESS_SETUP_RETR: DEFAULT_SETUP_RETR,
			ADDRESS_RF_CH:      DEFAULT_RF_CH,
			ADDRESS_RF_SETUP:   DEFAULT_RF_SETUP,
			ADDRESS_RX_PW_P0:   DEFAULT_RX_PW_P0,
			ADDRESS_STATUS:     DEFAULT_STATUS,
		},
		status: 0,
	}
	tr.init()
	return &tr
}

func bitEnable(value byte, bit byte, enable bool) byte {
	var mask byte = 0b00000001
	mask = mask << bit
	if enable {
		return value | mask
	}
	return value & ^mask

}

func (tr *nrf24l01dev) selectRadioMode(isRx bool) {
	tr.ceLow()
	tr.ClearStatus()
	tr.registers[ADDRESS_CONFIG] = bitEnable(tr.registers[ADDRESS_CONFIG], 0, isRx)
	tr.writeRegister(ADDRESS_CONFIG)
	tr.flushRx()
	tr.flushTx()
}

func (tr *nrf24l01dev) Listen() {
	tr.ceHigh()
}

func (tr *nrf24l01dev) TransmitterOn() {
	tr.selectRadioMode(false)
}

func (tr *nrf24l01dev) ReceiverOn() {
	tr.selectRadioMode(true)
}

func (tr *nrf24l01dev) PowerOn() {
	tr.setPower(true)
}

func (tr *nrf24l01dev) PowerOff() {
	tr.setPower(false)
}

func (tr *nrf24l01dev) ceHigh() {
	tr.ce.Out(gpio.High)
}

func (tr *nrf24l01dev) ceLow() {
	tr.ce.Out(gpio.Low)
}

func (tr *nrf24l01dev) setPower(on bool) {
	tr.registers[ADDRESS_CONFIG] = bitEnable(tr.registers[ADDRESS_CONFIG], 1, on)
	tr.writeRegister(ADDRESS_CONFIG)
	time.Sleep(time.Millisecond)
}

func (tr *nrf24l01dev) Transmit(payload models.Payload) error {
	_, err := writeSPI(ADDRESS_W_TX_PAYLOAD, payload[:], tr.conn)
	if err != nil {
		return err
	}
	tr.ceHigh()
	ts := time.Now()
	for time.Since(ts) < 5*time.Microsecond {
	}
	tr.ceLow()
	return err
}

func (tr *nrf24l01dev) Receive() (models.Payload, error) {
	tr.ceLow()
	payload := models.Payload{0, 0, 0, 0, 0, 0, 0, 0}
	data, err := readSPI(ADDRESS_R_RX_PAYLOAD, int(constants.RADIO_PAYLOAD_SIZE), tr.conn)
	if err != nil {
		return payload, err
	}
	copy(payload[:], data)
	tr.ClearStatus()
	return payload, err
}

func (tr *nrf24l01dev) init() {
	tr.PowerOff()
	tr.ceLow()
	time.Sleep(time.Millisecond)
	tr.writeConfigRegisters()
	tr.setRxTxAddress()
}

func (tr *nrf24l01dev) PrintConfigurations() {
	config, _ := tr.readRegister(ADDRESS_CONFIG)
	enaa, _ := tr.readRegister(ADDRESS_EN_AA)
	enrxaddr, _ := tr.readRegister(ADDRESS_EN_RXADDR)
	setupaw, _ := tr.readRegister(ADDRESS_SETUP_AW)
	setupretr, _ := tr.readRegister(ADDRESS_SETUP_RETR)
	rfch, _ := tr.readRegister(ADDRESS_RF_CH)
	rfsetup, _ := tr.readRegister(ADDRESS_RF_SETUP)
	rxpw0, _ := tr.readRegister(ADDRESS_RX_PW_P0)
	rxadd, _ := readSPI(ADDRESS_RX_ADDR_P0, 5, tr.conn)
	txadd, _ := readSPI(ADDRESS_TX_ADDR, 5, tr.conn)
	log.Printf(
		"\n	CONFIG: %b\n	EN_AA: %b\n	EN_RXADDR: %b\n	SETUP_AW: %b\n	SETUP_RETR: %b\n	RFCH: %b\n	RF_SETUP: %b\n	RX_PW0: %b\n	rx-add: %v\n	tx-add: %v",
		config, enaa, enrxaddr, setupaw, setupretr, rfch, rfsetup, rxpw0, rxadd, txadd)
}

func (tr *nrf24l01dev) readRegister(address byte) ([]byte, error) {
	return readSPI(address, 1, tr.conn)
}

func (tr *nrf24l01dev) writeRegister(address byte) {
	tr.ceLow()
	writeSPI(address, []byte{tr.registers[address]}, tr.conn)
}

func (tr *nrf24l01dev) writeConfigRegisters() {
	for address := range tr.registers {
		tr.writeRegister(address)
	}
}

func (tr *nrf24l01dev) setRxTxAddress() {
	tr.ceLow()
	writeSPI(ADDRESS_RX_ADDR_P0, tr.address, tr.conn)
	writeSPI(ADDRESS_TX_ADDR, tr.address, tr.conn)
}

func (tr *nrf24l01dev) updateStatus() {
	res, _ := readSPI(ADDRESS_STATUS, 1, tr.conn)
	tr.status = res[0]
}

func (tr *nrf24l01dev) IsReceiverDataReady(update bool) bool {
	if update {
		tr.updateStatus()
	}

	return tr.status&0b01000000 != 0
}

func (tr *nrf24l01dev) IsTransmitFailed(update bool) bool {
	if update {
		tr.updateStatus()
	}
	return tr.status&0b00010000 != 0
}

func (tr *nrf24l01dev) ClearStatus() {
	writeSPI(ADDRESS_STATUS, []byte{DEFAULT_STATUS}, tr.conn)
}

func (tr *nrf24l01dev) flushRx() {
	writeSPI(ADDRESS_FLUSH_RX, []byte{0x0}, tr.conn)
}

func (tr *nrf24l01dev) flushTx() {
	writeSPI(ADDRESS_FLUSH_TX, []byte{0x0}, tr.conn)
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
