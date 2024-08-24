package ads1115

import (
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

const ADS1X15_REG_POINTER_MASK =byte(0x03)      ///< Point mask
const ADS1X15_REG_POINTER_CONVERT =byte(0x00)   ///< Conversion
const ADS1X15_REG_POINTER_CONFIG =byte(0x01)    ///< Configuration
const ADS1X15_REG_POINTER_LOWTHRESH =byte(0x02) ///< Low threshold
const ADS1X15_REG_POINTER_HITHRESH =byte(0x03)  ///< High threshold

const ADS1X15_REG_CONFIG_CQUE_1CONV = 0x0000
const ADS1X15_REG_CONFIG_MODE_SINGLE = 0x0100
const ADS1X15_REG_CONFIG_MUX_SINGLE_0 = 0x4000 ///< Single-ended AIN0
const ADS1X15_REG_CONFIG_MUX_SINGLE_1 = 0x5000 ///< Single-ended AIN1
const ADS1X15_REG_CONFIG_MUX_SINGLE_2 = 0x6000 ///< Single-ended AIN2
const ADS1X15_REG_CONFIG_MUX_SINGLE_3 = 0x7000 ///< Single-ended AIN3
const ADS1X15_REG_CONFIG_OS_SINGLE = 0x8000
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

func (d *ads1115AtoD) Read(channel int) (l, h byte) {
	value:=d.ReadADC_SingleEnded(byte(channel))
	l=byte(value & 0b0000000011111111)
	h=byte(value>>8)
	return l, h
}

func (d *ads1115AtoD) ReadADC_SingleEnded(channel byte) int16 {
	if (channel > 3) {
	  return 0;
	}
  
	d.startADCReading(MUX_BY_CHANNEL[channel]);
  
	// Wait for the conversion to complete
	for !d.conversionComplete() {}
  
	// Read the conversion results
	return d.getLastConversionResults();
}

func (d *ads1115AtoD) startADCReading(mux uint16) {
	// Start with default values
	// uint16_t config =
	// 	ADS1X15_REG_CONFIG_CQUE_1CONV |   // Set CQUE to any value other than
	// 									  // None so we can use it in RDY mode
	// 	ADS1X15_REG_CONFIG_CLAT_NONLAT |  // Non-latching (default val)
	// 	ADS1X15_REG_CONFIG_CPOL_ACTVLOW | // Alert/Rdy active low   (default val)
	// 	ADS1X15_REG_CONFIG_CMODE_TRAD;    // Traditional comparator (default val)
	config := uint16(0)
  
	config |= ADS1X15_REG_CONFIG_MODE_SINGLE;
  
	// Set PGA/voltage range
	config |= uint16(AMP_FS_2)
  
	// Set data rate
	config |= uint16(DATA_RATE_128SPS)
  
	// Set channels
	config |= mux
  
	// Set 'start single-conversion' bit
	config |= ADS1X15_REG_CONFIG_OS_SINGLE
  
	// Write config register to the ADC
	d.writeRegister(ADS1X15_REG_POINTER_CONFIG, config)
  
	// Set ALERT/RDY to RDY mode.
	d.writeRegister(ADS1X15_REG_POINTER_HITHRESH, 0x8000)
	d.writeRegister(ADS1X15_REG_POINTER_LOWTHRESH, 0x0000)
  }

func (d *ads1115AtoD) writeRegister(reg byte, value uint16) {
	w := []byte{reg, byte(value >> 8), byte(value & 0xFF)}
	d.i2cDev.Tx(w, nil)
}


func (d *ads1115AtoD) conversionComplete() bool {
	return (d.readRegister(ADS1X15_REG_POINTER_CONFIG) & 0x8000) != 0;
}

func (d *ads1115AtoD) readRegister(reg byte) uint16 {
	// buffer[0] = reg;
	// m_i2c_dev->write(buffer, 1);
	// m_i2c_dev->read(buffer, 2);
	// return ((buffer[0] << 8) | buffer[1]);
	r := make([]byte, 2)
	w := []byte{reg}
	d.i2cDev.Tx(w, r)
	return uint16(r[0]) <<8 | uint16(r[1])
}

func (d *ads1115AtoD) getLastConversionResults() int16 {
	const m_bitShift = 0
	res := d.readRegister(ADS1X15_REG_POINTER_CONVERT) >> m_bitShift
	if m_bitShift == 0 {
	  return int16(res);
	} else {
	  // Shift 12-bit results right 4 bits for the ADS1015,
	  // making sure we keep the sign bit intact
	  if res > 0x07FF {
		// negative number - extend the sign to 16th bit
		res |= 0xF000;
	  }
	  return int16(res);
	}
}
