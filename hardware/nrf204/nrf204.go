package nrf204

const (
	RF_SETUP      byte = 0x06
	RF24_PA_MIN   byte = 0
	RF24_PA_LOW   byte = 1
	RF24_PA_HIGH  byte = 2
	RF24_PA_MAX   byte = 3
	RF24_PA_ERROR byte = 4
)

const addressSize int = 5

type nrf204l01 struct {
	address []byte
}

func CreateNRF204() *nrf204l01 {
	return &nrf204l01{
		address: make([]byte, addressSize),
	}
}

func (r *nrf204l01) Init() {
}

func (r *nrf204l01) OpenReadingPipe() {
}

func (r *nrf204l01) SetPALevel(level byte, lnaEnable byte) {
	setup := r.readRegister(RF_SETUP) & 0xF8
	l := level
	if l > 3 {
		l = (RF24_PA_MAX << 1) + lnaEnable
	} else {
		l = (l << 1) + lnaEnable
	}
	r.writeRegister(RF_SETUP, setup|l)
}

func (r *nrf204l01) StartListening() {
}

func (r *nrf204l01) IsAvailable() bool {
	return false
}

func (r *nrf204l01) Read() []byte {
	return []byte{0, 0, 0, 0}
}

func (r *nrf204l01) readRegister(address byte) byte {
	return 0
}
func (r *nrf204l01) writeRegister(address byte, value byte) {
}
