package types

type IMUConfigs struct {
	AccelerometerFullScale string
	GyroscopeFullScale     string
}

type Configs struct {
	IMU IMUConfigs
}
