package ads1115

import (
	"fmt"
	"periph.io/x/conn/v3/i2c"
)

const (
	CONVERSION_REGISTER_ADDRESS = byte(0x0)
	CONFIG_RREGISTER_ADDRESS =    byte(0x1)
	LO_THRESH_REGISTER_ADDRESS =  byte(0x2)
	HI_THRESH_REGISTER_ADDRESS =  byte(0x3)
	CHANNEL_0 =                   byte(0b01000000)
	CHANNEL_1 =                   byte(0b01010000)
	CHANNEL_2 =                   byte(0b01100000)
	CHANNEL_3 =                   byte(0b01110000)
	AMP_FS_6  =                   byte(0b00000000)
	AMP_FS_4  =                   byte(0b00000010)
	AMP_FS_2  =                   byte(0b00000100)
	AMP_FS_1  =                   byte(0b00000110)
	DATA_RATE_8SPS   =            byte(0b00000000)
	DATA_RATE_16SPS  =            byte(0b00100000)
	DATA_RATE_32SPS  =            byte(0b01000000)
	DATA_RATE_64SPS  =            byte(0b01100000)
	DATA_RATE_128SPS =            byte(0b10000000)
	COMP_QUE =                    byte(0b00000011)
)


const ADS1X15_REG_CONFIG_MUX_SINGLE_0 = 0x4000 ///< Single-ended AIN0
const ADS1X15_REG_CONFIG_MUX_SINGLE_1 = 0x5000 ///< Single-ended AIN1
const ADS1X15_REG_CONFIG_MUX_SINGLE_2 = 0x6000 ///< Single-ended AIN2
const ADS1X15_REG_CONFIG_MUX_SINGLE_3 = 0x7000 ///< Single-ended AIN3

var MUX_BY_CHANNEL []uint16 = []uint16 {
    ADS1X15_REG_CONFIG_MUX_SINGLE_0, ///< Single-ended AIN0
    ADS1X15_REG_CONFIG_MUX_SINGLE_1, ///< Single-ended AIN1
    ADS1X15_REG_CONFIG_MUX_SINGLE_2, ///< Single-ended AIN2
    ADS1X15_REG_CONFIG_MUX_SINGLE_3,  ///< Single-ended AIN3
}; 

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
	fmt.Println(r, err)
	return r, err
}

func (d *ads1115AtoD) Setup(channel int) error{
	config := COMP_QUE | DATA_RATE_128SPS
	w := []byte{CONFIG_RREGISTER_ADDRESS, config}
	err := d.i2cDev.Tx(w, nil)
	// fmt.Println(err)
	return err
}

func (d *ads1115AtoD) WriteConfigs(channel int) error{
	// config0 := COMP_QUE | 
	w := []byte{CONFIG_RREGISTER_ADDRESS, 0, 0}
	err := d.i2cDev.Tx(w, nil)
	// fmt.Println(err)
	return err
}

func (d *ads1115AtoD) ReadConfigs() ([]byte, error) {
	r := make([]byte, 2)
	w := []byte{CONFIG_RREGISTER_ADDRESS}
	err := d.i2cDev.Tx(w, r)
	// fmt.Println(r[0], r[1], err)
	return r, err
}

func (d *ads1115AtoD) ReadThreshold() ([]byte, error) {
	r := make([]byte, 2)
	w := []byte{LO_THRESH_REGISTER_ADDRESS}
	err := d.i2cDev.Tx(w, r)
	fmt.Println(r, err)
	w = []byte{HI_THRESH_REGISTER_ADDRESS}
	err = d.i2cDev.Tx(w, r)
	fmt.Println(r, err)	
	return r, err
}


// func (d *ads1115AtoD) writeByte(address byte) error {
// 	write := []byte{offset, b}
// 	_, err := d.i2cDev.Write(write)
// 	return err
// }

// int16_t Adafruit_ADS1X15::readADC_SingleEnded(uint8_t channel) {
// 	if (channel > 3) {
// 	  return 0;
// 	}
  
// 	startADCReading(MUX_BY_CHANNEL[channel], /*continuous=*/false);
  
// 	// Wait for the conversion to complete
// 	while (!conversionComplete())
// 	  ;
  
// 	// Read the conversion results
// 	return getLastConversionResults();
//   }

func (d *ads1115AtoD) ReadADC_SingleEnded(channel byte) uint16 {
	if (channel > 3) {
	  return 0;
	}
	return 1
  
	// startADCReading(MUX_BY_CHANNEL[channel], /*continuous=*/false);
  
	// // Wait for the conversion to complete
	// while (!conversionComplete())
	//   ;
  
	// // Read the conversion results
	// return getLastConversionResults();
  }