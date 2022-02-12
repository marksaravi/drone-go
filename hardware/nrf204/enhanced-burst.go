package nrf204

import (
	"fmt"
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
	tr.enhancedBurstInit()
	tr.enhancedBurstReadConfigRegisters()
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
	tr.registers[ADDRESS_CONFIG] = bitEnable(tr.registers[ADDRESS_CONFIG], 0, on)
	tr.ebApplyRegister(ADDRESS_CONFIG)
}

func (tr *nrf204l01) TransmitterOn() {
	tr.receiverOn(false)
}

func (tr *nrf204l01) ReceiverOn() {
	tr.receiverOn(true)
	tr.ce.Out(gpio.High)
}

func (tr *nrf204l01) PowerOn() {
	tr.setPower(true)
}

func (tr *nrf204l01) PowerOff() {
	tr.setPower(false)
}

func (tr *nrf204l01) setPower(on bool) {
	tr.registers[ADDRESS_CONFIG] = bitEnable(tr.registers[ADDRESS_CONFIG], 1, on)
	tr.ebApplyRegister(ADDRESS_CONFIG)
}

func (tr *nrf204l01) Transmit(payload models.Payload) error {
	fmt.Println(payload)
	_, err := writeSPI(ADDRESS_W_TX_PAYLOAD, payload[:], tr.conn)
	if err != nil {
		return err
	}
	err = tr.ce.Out(gpio.High)
	if err != nil {
		return err
	}
	ts := time.Now()
	for time.Since(ts) < 5*time.Microsecond {
	}
	tr.ce.Out(gpio.Low)
	return err
}

func (tr *nrf204l01) Receive() (models.Payload, error) {
	payload := models.Payload{0, 0, 0, 0, 0, 0, 0, 0}
	data, err := readSPI(ADDRESS_W_TX_PAYLOAD, int(constants.RADIO_PAYLOAD_SIZE), tr.conn)
	if err != nil {
		return payload, err
	}
	copy(payload[:], data)
	return payload, err
}

func (tr *nrf204l01) enhancedBurstInit() {
	tr.ce.Out(gpio.Low)
	tr.setPower(false)
	time.Sleep(time.Millisecond)
	tr.ebSetRegisters()
	tr.setRxTxAddress()
	tr.setPower(true)
	time.Sleep(time.Millisecond)
}

func (tr *nrf204l01) enhancedBurstReadConfigRegisters() {
	config, _ := tr.ebReadConfigRegister(ADDRESS_CONFIG)
	enaa, _ := tr.ebReadConfigRegister(ADDRESS_EN_AA)
	enrxaddr, _ := tr.ebReadConfigRegister(ADDRESS_EN_RXADDR)
	setupaw, _ := tr.ebReadConfigRegister(ADDRESS_SETUP_AW)
	setupretr, _ := tr.ebReadConfigRegister(ADDRESS_SETUP_RETR)
	rfch, _ := tr.ebReadConfigRegister(ADDRESS_RF_CH)
	rfsetup, _ := tr.ebReadConfigRegister(ADDRESS_RF_SETUP)
	log.Printf(
		"\n	CONFIG: %b\n	EN_AA: %b\n	EN_RXADDR: %b\n	SETUP_AW: %b\n	SETUP_RETR: %b\n	RFCH: %b\n	RF_SETUP: %b",
		config, enaa, enrxaddr, setupaw, setupretr, rfch, rfsetup)
}

func (tr *nrf204l01) ebReadConfigRegister(address byte) ([]byte, error) {
	return readSPI(address, 1, tr.conn)
}

func (tr *nrf204l01) ebApplyRegister(address byte) {
	writeSPI(address, []byte{tr.registers[address]}, tr.conn)
}

func (tr *nrf204l01) ebSetRegisters() {
	for address := range tr.registers {
		tr.ebApplyRegister(address)
	}
}

func (tr *nrf204l01) setRxTxAddress() {
	writeSPI(ADDRESS_RX_ADDR_P0, tr.address, tr.conn)
	writeSPI(ADDRESS_TX_ADDR, tr.address, tr.conn)
}

func (tr *nrf204l01) UpdateStatus() {
	res, _ := writeSPI(0b11111111, []byte{0}, tr.conn)
	tr.status = res[0]
}

func (tr *nrf204l01) ReceiverDataReady(update bool) bool {
	if update {
		tr.UpdateStatus()
	}
	return tr.status&0b01000000 != 0
}

func (tr *nrf204l01) TransmitFailed(update bool) bool {
	if update {
		tr.UpdateStatus()
	}
	return tr.status&0b00100000 != 0
}

func (tr *nrf204l01) ResetStatus() {
	writeSPI(ADDRESS_STATUS, []byte{DEFAULT_STATUS}, tr.conn)
}
