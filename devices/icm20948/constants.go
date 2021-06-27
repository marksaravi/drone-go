package icm20948

const (
	BANK0 uint16 = 0 << 8
	BANK1 uint16 = 1 << 8
	BANK2 uint16 = 2 << 8
	BANK3 uint16 = 3 << 8
)

const (
	REG_BANK_SEL byte = 0x7F

	// BANK0
	WHO_AM_I     uint16 = BANK0 | 0x0
	LP_CONFIG    uint16 = BANK0 | 0x5
	PWR_MGMT_1   uint16 = BANK0 | 0x6
	PWR_MGMT_2   uint16 = BANK0 | 0x7
	INT_ENABLE_3 uint16 = BANK0 | 0x13
	ACCEL_XOUT_H uint16 = BANK0 | 0x2D
	ACCEL_XOUT_L uint16 = BANK0 | 0x2E
	ACCEL_YOUT_H uint16 = BANK0 | 0x2F
	ACCEL_YOUT_L uint16 = BANK0 | 0x30
	ACCEL_ZOUT_H uint16 = BANK0 | 0x31
	ACCEL_ZOUT_L uint16 = BANK0 | 0x32
	GYRO_XOUT_H  uint16 = BANK0 | 0x33
	GYRO_XOUT_L  uint16 = BANK0 | 0x34
	GYRO_YOUT_H  uint16 = BANK0 | 0x35
	GYRO_YOUT_L  uint16 = BANK0 | 0x36
	GYRO_ZOUT_H  uint16 = BANK0 | 0x37
	GYRO_ZOUT_L  uint16 = BANK0 | 0x38

	// BANK1
	XA_OFFS_H uint16 = BANK1 | 0x14
	XA_OFFS_L uint16 = BANK1 | 0x15
	YA_OFFS_H uint16 = BANK1 | 0x17
	YA_OFFS_L uint16 = BANK1 | 0x18
	ZA_OFFS_H uint16 = BANK1 | 0x1A
	ZA_OFFS_L uint16 = BANK1 | 0x1B

	// BANK2
	GYRO_SMPLRT_DIV uint16 = BANK2 | 0x0
	GYRO_CONFIG_1   uint16 = BANK2 | 0x1
	GYRO_CONFIG_2   uint16 = BANK2 | 0x2
	ZG_OFFS_USRH    uint16 = BANK2 | 0x7
	ACCEL_CONFIG    uint16 = BANK2 | 0x14
	ACCEL_CONFIG_2  uint16 = BANK2 | 0x15
	MOD_CTRL_USR    uint16 = BANK2 | 0x54
)

const (
	SENSITIVITY_2G  float64 = 16384
	SENSITIVITY_4G  float64 = 8192
	SENSITIVITY_8G  float64 = 4096
	SENSITIVITY_16G float64 = 2048

	SCALE_0 float64 = 131
	SCALE_1 float64 = 65.5
	SCALE_2 float64 = 32.8
	SCALE_3 float64 = 16.4
)

const (
	ACCELEROMETER = "accelerometer"
	GYROSCOPE     = "gyroscope"
	MAGNETOMETER  = "magnetometer"
)
