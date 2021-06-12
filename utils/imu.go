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

func GyroRotations(dg types.RotationsChanges, gyroRotations types.Rotations) types.Rotations {
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

func applyFilter(pR float64, accR float64, gyroDR float64, lpfc float64) float64 {
	return lpfc*(pR+gyroDR) + (1-lpfc)*accR
}

func CalcRotations(pR types.Rotations, aR types.Rotations, dg types.RotationsChanges, lowPassFilterCoefficient float64) types.Rotations {
	//  nR = lpfc * (pR + gyroDPS * timeDelta ) + (1-lpfc) * accR;
	roll := applyFilter(pR.Roll, aR.Roll, dg.DRoll, lowPassFilterCoefficient)
	pitch := applyFilter(pR.Pitch, aR.Pitch, dg.DPitch, lowPassFilterCoefficient)
	yaw := applyFilter(pR.Yaw, aR.Yaw, dg.DYaw, lowPassFilterCoefficient)
	return types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   yaw,
	}
}
