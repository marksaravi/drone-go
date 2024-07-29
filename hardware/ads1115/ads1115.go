package ads1115

import (
	"periph.io/x/conn/v3/i2c"
)

const (
	CONVERSION_REGISTER_ADDRESS = byte(0)
	CONFIG_RREGISTER_ADDRESS =    byte(1)
	LO_THRESH_REGISTER_ADDRESS =  byte(2)
	HI_THRESH_REGISTER_ADDRESS =  byte(3)
)

type ads1115AtoD struct {
	i2cDev  *i2c.Dev
}

func NewADS1115(i2cDev  *i2c.Dev) *ads1115AtoD {
	return &ads1115AtoD {
		i2cDev: i2cDev,
	}
}

func (d *ads1115AtoD) Read(channel int) int {
	b,_ := d.readConversion()
	uv := uint16(b[0]) | uint16(b[1])<<8
	return int(uv)
}

func (d *ads1115AtoD) readConversion() ([]byte, error) {
	r := make([]byte, 2)
	w := []byte{CONVERSION_REGISTER_ADDRESS}
	err := d.i2cDev.Tx(w, r)
	return r, err
}

// func (d *ads1115AtoD) writeByte(address byte) error {
// 	write := []byte{offset, b}
// 	_, err := d.i2cDev.Write(write)
// 	return err
// }