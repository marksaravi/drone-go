package utils

import (
	"math"

	"github.com/MarkSaravi/drone-go/types"
)

func goDurToDt(d int64) float64 {
	return float64(d) / 1e9
}

func GyroChanges(imuSensorsData types.ImuSensorsData) types.RotationsChanges {
	dt := goDurToDt(imuSensorsData.ReadInterval)
	return types.RotationsChanges{
		DRoll:  imuSensorsData.Gyro.Data.X * dt,
		DPitch: imuSensorsData.Gyro.Data.Y * dt,
		DYaw:   imuSensorsData.Gyro.Data.Z * dt,
	}
}

func GyroRotations(imuSensorsData types.ImuSensorsData, gyroRotations types.Rotations) types.Rotations {
	dg := GyroChanges(imuSensorsData)
	return types.Rotations{
		Roll:  gyroRotations.Roll + dg.DRoll,
		Pitch: gyroRotations.Pitch + dg.DPitch,
		Yaw:   gyroRotations.Yaw + dg.DYaw,
	}
}

func AccelerometerDataRotations(data types.XYZ) types.Rotations {
	roll := RadToDeg(math.Atan2(data.Y, data.Z))
	pitch := -RadToDeg(math.Atan2(data.X, data.Z))
	return types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   0,
	}
}
