package nrf204

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/models"
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

	tr := nrf204l01{
		address:    []byte(rxTxAddress),
		ce:         initPin(spiChipEnabledGPIO),
		conn:       radioSPIConn,
		powerDBm:   dbmStrToDBm(powerDb),
		isReceiver: true,
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

func (tr *nrf204l01) SetTransmitter(on bool) {
	registers[ADDRESS_CONFIG] = bitEnable(registers[ADDRESS_CONFIG], 0, on)
	tr.ebApplyRegister(ADDRESS_CONFIG)
}

func (tr *nrf204l01) setPower(on bool) {
	registers[ADDRESS_CONFIG] = bitEnable(registers[ADDRESS_CONFIG], 1, on)
	tr.ebApplyRegister(ADDRESS_CONFIG)
}

func (tr *nrf204l01) Transmit(payload models.Payload) error {
	_, err := tr.ebWriteRegisterBytes(W_TX_PAYLOAD, payload[:])
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
	return readSPI(address|R_REGISTER_MASK, 1, tr.conn)
}

func (tr *nrf204l01) ebApplyRegister(address byte) {
	writeSPI(address|W_REGISTER_MASK, []byte{registers[address]}, tr.conn)
}

func (tr *nrf204l01) ebWriteRegisterBytes(address byte, values []byte) ([]byte, error) {
	return writeSPI(address|W_REGISTER_MASK, values, tr.conn)
}

func (tr *nrf204l01) ebSetRegisters() {
	for address := ADDRESS_CONFIG; address <= ADDRESS_RF_SETUP; address++ {
		tr.ebApplyRegister(address)
	}
}

func (tr *nrf204l01) setRxTxAddress() {
	tr.ebWriteRegisterBytes(ADDRESS_RX_ADDR_P0, tr.address)
	tr.ebWriteRegisterBytes(ADDRESS_TX_ADDR, tr.address)
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
