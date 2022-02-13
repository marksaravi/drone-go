package nrf204

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/constants"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/models"
	"periph.io/x/periph/conn/gpio"
)

func NewNRF204EnhancedBurst(
	spiBusNum int,
	spiChipSelect int,
	spiChipEnabledGPIO string,
	rxTxAddress string,
) *nrf204l01 {
	radioSPIConn := hardware.NewSPIConnection(
		spiBusNum,
		spiChipSelect,
	)

	tr := nrf204l01{
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
	tr.readConfigurations()
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

func (tr *nrf204l01) receiverOn(on bool) {
	tr.ClearStatus()
	tr.ceLow()
	tr.registers[ADDRESS_CONFIG] = bitEnable(tr.registers[ADDRESS_CONFIG], 0, on)
	tr.writeRegister(ADDRESS_CONFIG)
	tr.flushRx()
	tr.flushTx()
	tr.PowerOn()
	tr.ceHigh()
}

func (tr *nrf204l01) TransmitterOn() {
	tr.receiverOn(false)
}

func (tr *nrf204l01) ReceiverOn() {
	tr.receiverOn(true)
}

func (tr *nrf204l01) PowerOn() {
	tr.setPower(true)
}

func (tr *nrf204l01) PowerOff() {
	tr.setPower(false)
}

func (tr *nrf204l01) ceHigh() {
	tr.ce.Out(gpio.High)
}

func (tr *nrf204l01) ceLow() {
	tr.ce.Out(gpio.Low)
}

func (tr *nrf204l01) setPower(on bool) {
	tr.registers[ADDRESS_CONFIG] = bitEnable(tr.registers[ADDRESS_CONFIG], 1, on)
	tr.writeRegister(ADDRESS_CONFIG)
	time.Sleep(time.Millisecond)
}

func (tr *nrf204l01) Transmit(payload models.Payload) error {
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

func (tr *nrf204l01) Receive() (models.Payload, error) {
	tr.ceLow()
	payload := models.Payload{0, 0, 0, 0, 0, 0, 0, 0}
	data, err := readSPI(ADDRESS_R_RX_PAYLOAD, int(constants.RADIO_PAYLOAD_SIZE), tr.conn)
	if err != nil {
		return payload, err
	}
	copy(payload[:], data)
	return payload, err
}

func (tr *nrf204l01) init() {
	tr.ceLow()
	tr.setPower(false)
	time.Sleep(time.Millisecond)
	tr.writeConfigRegisters()
	tr.setRxTxAddress()
	tr.setPower(true)
	time.Sleep(time.Millisecond)
}

func (tr *nrf204l01) readConfigurations() {
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

func (tr *nrf204l01) readRegister(address byte) ([]byte, error) {
	return readSPI(address, 1, tr.conn)
}

func (tr *nrf204l01) writeRegister(address byte) {
	tr.ceLow()
	writeSPI(address, []byte{tr.registers[address]}, tr.conn)
}

func (tr *nrf204l01) writeConfigRegisters() {
	for address := range tr.registers {
		tr.writeRegister(address)
	}
}

func (tr *nrf204l01) setRxTxAddress() {
	tr.ceLow()
	writeSPI(ADDRESS_RX_ADDR_P0, tr.address, tr.conn)
	writeSPI(ADDRESS_TX_ADDR, tr.address, tr.conn)
}

func (tr *nrf204l01) UpdateStatus() {
	res, _ := readSPI(ADDRESS_STATUS, 1, tr.conn)
	tr.status = res[0]
}

func (tr *nrf204l01) IsReceiverDataReady(update bool) bool {
	if update {
		tr.UpdateStatus()
	}

	return tr.status&0b01000000 != 0
}

func (tr *nrf204l01) IsTransmitFailed(update bool) bool {
	if update {
		tr.UpdateStatus()
	}
	return tr.status&0b00010000 != 0
}

func (tr *nrf204l01) ClearStatus() {
	tr.ceLow()
	writeSPI(ADDRESS_STATUS, []byte{DEFAULT_STATUS}, tr.conn)
}

func (tr *nrf204l01) flushRx() {
	writeSPI(ADDRESS_FLUSH_RX, []byte{0xFF}, tr.conn)
}

func (tr *nrf204l01) flushTx() {
	writeSPI(ADDRESS_FLUSH_TX, []byte{0xFF}, tr.conn)
}
