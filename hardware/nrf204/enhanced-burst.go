package nrf204

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"periph.io/x/periph/conn/gpio"
)

func NewNRF204EnhancedBurst(
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
	radio.enhancedBurstInit()
	radio.enhancedBurstReadConfigRegisters()
	return &radio
}
func (tr *nrf204l01) enhancedBurstInit() {
	tr.ce.Out(gpio.Low)
	tr.setPower(OFF)
	time.Sleep(time.Millisecond)
	tr.ebSetRegisters()
	tr.setTransmitterAddress()
	tr.setReceiverAddress()
	tr.setPower(ON)
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
	return readSPI(address|R_REGISTER_MASK, 1, tr.conn)
}

func (tr *nrf204l01) ebSetRegisters() {

	for address := ADDRESS_CONFIG; address <= ADDRESS_RF_SETUP; address++ {
		writeSPI(address|W_REGISTER_MASK, []byte{configRegisters[address]}, tr.conn)
	}
}

func (radio *nrf204l01) setTransmitterAddress() {
	radio.writeRegister(ADDRESS_TX_ADDR, radio.address)
}

func (radio *nrf204l01) setReceiverAddress() {
	radio.writeRegister(ADDRESS_RX_ADDR_P0, radio.address)
}

func (tr *nrf204l01) ebCommitConfigRegister(address byte) {
	writeSPI(address|W_REGISTER_MASK, []byte{configRegisters[address]}, tr.conn)
}

func (tr *nrf204l01) ebReadRxPayload() {

}

func (tr *nrf204l01) ebWriteTxPayload() {

}

func (tr *nrf204l01) ebFlushRx() {

}

func (tr *nrf204l01) ebFlushTx() {

}

func (tr *nrf204l01) ebReadStatusRegister() {

}
