package nrf204

import (
	"log"

	"periph.io/x/periph/conn/gpio"
)

func (tr *nrf204l01) enhancedBurstInit() {
	tr.ce.Out(gpio.Low)
	tr.ebCommitConfigRegister(ADDRESS_CONFIG)
	tr.ebCommitConfigRegister(ADDRESS_EN_AA)
	tr.ebCommitConfigRegister(ADDRESS_EN_RXADDR)
	tr.ebCommitConfigRegister(ADDRESS_SETUP_AW)
	tr.ebCommitConfigRegister(ADDRESS_SETUP_RETR)
	tr.ebCommitConfigRegister(ADDRESS_RF_CH)
	tr.ebCommitConfigRegister(ADDRESS_RF_SETUP)
	tr.setPayloadSize()
	tr.setTransmitterAddress()
	tr.setReceiverAddress()
	tr.setPALevel(tr.powerDBm)
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
