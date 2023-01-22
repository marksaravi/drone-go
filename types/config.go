package types

type IMUConfigs struct {
	AccelerometerFullScale string
	GyroscopeFullScale     string
	FilterCoefficient      float64 // Complementary Filter or Kalman filter Coefficient
}

type Configs struct {
	IMU IMUConfigs
}
