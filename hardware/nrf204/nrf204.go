package nrf204

import (
	"fmt"
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/spi"
)

const (
	NRF_CONFIG    byte = 0x00
	RF_SETUP      byte = 0x06
	NRF_STATUS    byte = 0x07
	RF24_PA_MIN   byte = 0
	RF24_PA_LOW   byte = 1
	RF24_PA_HIGH  byte = 2
	RF24_PA_MAX   byte = 3
	RF24_PA_ERROR byte = 4
	RX_ADDRESS    byte = 0x0A
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
	// SPI1 only supports Mode0
	//to enable SPI1 in raspberry pi follow instructions here https://docs.rs/rppal/0.8.1/rppal/spi/index.html
	// or add "dtoverlay=spi1-2cs" to /boot/config.txt
	return &nrf204l01{
		ce:      initPin(config.GPIO.CE),
		address: make([]byte, addressSize),
		conn:    conn,
	}
}

func (radio *nrf204l01) Init() {
	radio.setRadioMode(Receiver)
	time.Sleep(time.Millisecond * 200)
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
	r, err := radio.readRegister(RF_SETUP, 1)
	if err != nil {
		log.Fatal(err)
	}
	setup := r[0] & 0xF8
	l := level
	if l > 3 {
		l = (RF24_PA_MAX << 1) + lnaEnable
	} else {
		l = (l << 1) + lnaEnable
	}
	data := []byte{setup | l}
	radio.writeRegister(RF_SETUP, data)
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

func (radio *nrf204l01) initRadio() {
}

func (radio *nrf204l01) setRadioMode(mode RadioMode) {
	if mode == Receiver {
		radio.ce.Out(gpio.Low)
	} else {
		radio.ce.Out(gpio.High)
	}
}

func (radio *nrf204l01) powerUp() {
}

func (radio *nrf204l01) readRegister(address byte, datalen int) ([]byte, error) {
	w := make([]byte, datalen+1)
	r := make([]byte, datalen+1)
	w[0] = address
	err := radio.conn.Tx(w, r)
	return r[1:], err
}
func (radio *nrf204l01) writeRegister(address byte, data []byte) error {
	datalen := len(data)
	w := make([]byte, datalen+1)
	w[0] = address
	for i := 0; i < datalen; i++ {
		w[i+1] = data[i]
	}
	err := radio.conn.Tx(w, nil)
	return err
}

func initPin(pinName string) gpio.PinIO {
	pin := gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	return pin
}
