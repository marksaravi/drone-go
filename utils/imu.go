package utils

import (
	"fmt"
	"math"
	"time"

	"github.com/MarkSaravi/drone-go/types"
)

var lastReading = time.Now()

func Print(v []float64, msInterval int) {
	if time.Since(lastReading) >= time.Millisecond*time.Duration(msInterval) {
		fmt.Println(v)
		lastReading = time.Now()
	}
}

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

func GyroRotations(dg types.RotationsChanges, gyroRotations types.Rotations) types.Rotations {
	return types.Rotations{
		Roll:  math.Mod(gyroRotations.Roll+dg.DRoll, 360),
		Pitch: math.Mod(gyroRotations.Pitch+dg.DPitch, 360),
		Yaw:   math.Mod(gyroRotations.Yaw+dg.DYaw, 360),
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

func applyFilter(pR float64, accR float64, gyroDR float64, lpfc float64) float64 {
	nR := math.Mod(pR+gyroDR, 360)
	if math.Abs(nR) >= 90 {
		// preventint acc correction for more than 90° as values can have different sign
		return nR
	}
	return lpfc*nR + (1-lpfc)*accR
}

func CalcRotations(pR types.Rotations, aR types.Rotations, dg types.RotationsChanges, lowPassFilterCoefficient float64) types.Rotations {
	roll := applyFilter(pR.Roll, aR.Roll, dg.DRoll, lowPassFilterCoefficient)
	pitch := applyFilter(pR.Pitch, aR.Pitch, dg.DPitch, lowPassFilterCoefficient)
	yaw := applyFilter(pR.Yaw, aR.Yaw, dg.DYaw, lowPassFilterCoefficient)
	return types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   yaw,
	}
}
