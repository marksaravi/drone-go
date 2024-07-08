package utils_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/marksaravi/drone-go/utils"
)

func TestFloat64Average(t *testing.T) {
	numOfSamples := 10
	float64Average := utils.NewAverage[float64](numOfSamples)
	values := make([]float64, numOfSamples)
	index := 0
	for i := 0; i < 1000000; i++ {
		values[index] = rand.Float64()
		float64Average.AddValue(values[index])
		index++
		if index == numOfSamples {
			index = 0
		}
	}
	sum := float64(0)
	for i := 0; i < numOfSamples; i++ {
		sum += values[i]
	}
	actualAverage := sum / float64(numOfSamples)
	if math.Abs(actualAverage-float64Average.Average()) > float64(1e-12) {
		t.Errorf("%f, %f", actualAverage, float64Average.Average())
	}
}

func TestIntAverage(t *testing.T) {
	numOfSamples := 10
	float64Average := utils.NewAverage[int](numOfSamples)
	values := make([]int, numOfSamples)
	index := 0
	for i := 0; i < 1000000; i++ {
		values[index] = int(rand.Intn(1000))
		float64Average.AddValue(values[index])
		index++
		if index == numOfSamples {
			index = 0
		}
	}
	sum := int(0)
	for i := 0; i < numOfSamples; i++ {
		sum += values[i]
	}
	actualAverage := sum / numOfSamples
	if actualAverage != float64Average.Average() {
		t.Errorf("%d, %d", actualAverage, float64Average.Average())
	}
}
