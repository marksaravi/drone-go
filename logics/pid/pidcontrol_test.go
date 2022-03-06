package pid

import (
	"testing"
)

func createPIDControl() *pidControl {
	return &pidControl{
		pGain:         0,
		iGain:         0,
		dGain:         0,
		maxOutput:     0,
		maxI:          0,
		previousInput: 0,
		iMemory:       0,
	}
}
func TestGain(t *testing.T) {

}
