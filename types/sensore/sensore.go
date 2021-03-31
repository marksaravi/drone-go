package sensore

import "github.com/MarkSaravi/drone-go/types"

// ThreeAxisSensore is interface to a 3 Axis sensore
type ThreeAxisSensore interface {
	GetConfig() types.Config
	SetConfig(config types.Config)
	GetData() types.XYZ
	SetData(x, y, z float64)
	GetDiff() float64
}
