package registers

const (
	BANK0 byte = 0
	BANK1 byte = 1
	BANK2 byte = 2
	BANK3 byte = 3
)

const (
	REG_BANK_SEL byte = 0x7F

	// BANK0
	WHO_AM_I     byte = 0x0
	LP_CONFIG    byte = 0x5
	PWR_MGMT_1   byte = 0x6
	PWR_MGMT_2   byte = 0x7
	INT_ENABLE_3 byte = 0x13
	ACCEL_ZOUT_H byte = 0x31
	ACCEL_ZOUT_L byte = 0x32
	GYRO_ZOUT_L  byte = 0x38

	// BANK1
	XA_OFFS_H byte = 0x14

	// BANK2
	GYRO_SMPLRT_DIV byte = 0x0
	GYRO_CONFIG_1   byte = 0x1
	GYRO_CONFIG_2   byte = 0x2
	ZG_OFFS_USRL    byte = 0x8
	ACCEL_CONFIG_2  byte = 0x15
	MOD_CTRL_USR    byte = 0x54
)
