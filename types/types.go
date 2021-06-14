package types

// Config is the generic configuration
type Config interface {
}

type FlightConfig struct {
	ImuDataPerSecond         int     `yaml:"imu_data_per_second"`
	LowPassFilterCoefficient float64 `yaml:"lowpass_filter_coefficient"`
	AccelerometerOffsets     Offsets
	GyroscopeOffsets         Offsets
}

type UdpLoggerConfig struct {
	Enabled          bool   `yaml:"enabled"`
	IP               string `yaml:"ip"`
	Port             int    `yaml:"port"`
	PacketsPerSecond int    `yaml:"udp_packets_per_second"`
	PrintIntervalMs  int    `yaml:"print_interval_ms"`
}

// XYZ is X, Y, Z data
type XYZ struct {
	X, Y, Z float64
}

type Rotations struct {
	Roll, Pitch, Yaw float64
}

type RotationsChanges struct {
	DRoll, DPitch, DYaw float64
}

type SensorData struct {
	Error error
	Data  XYZ
}

type ImuSensorsData struct {
	Acc, Gyro, Mag SensorData
	ReadTime       int64
	ReadInterval   int64
}

type ImuRotations struct {
	PrevRotations Rotations
	Accelerometer Rotations
	Gyroscope     Rotations
	Rotations     Rotations
	ReadTime      int64
	ReadInterval  int64
}

// Sensor is devices that read data in x, y, z format
type Sensor struct {
	Type   string
	Config Config
}

// CommandParameters is parameters for the command
type CommandParameters interface {
}

type Command struct {
	Command    string
	Parameters CommandParameters
}

// GetConfig reads the config
func (a *Sensor) GetConfig() Config {
	return a.Config
}

// SetConfig sets the config
func (a *Sensor) SetConfig(config Config) {
	a.Config = config
}

type Offsets struct {
	X float64 `yaml:"X"`
	Y float64 `yaml:"Y"`
	Z float64 `yaml:"Z"`
}

// ImuMems is interface for the imu mems
type ImuMems interface {
	Close()
	InitDevice() error
	ReadSensorsRawData() ([]byte, error)
	ReadSensors() (acc SensorData, gyro SensorData, mag SensorData, err error)
	WhoAmI() (string, byte, error)
}

// IMU is interface to imu mems
type IMU interface {
	Close()
	GetRotations() (ImuRotations, error)
}
