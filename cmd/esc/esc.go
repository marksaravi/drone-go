package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/devices/pca9685"
	"github.com/MarkSaravi/drone-go/drivers/i2c"
	"github.com/MarkSaravi/drone-go/modules/esc"
)

func main() {
	channel := flag.Int("ch", 0, "ESC channel")
	frequency := flag.Float64("freq", 400, "Frequency")
	pulseWidth := flag.Float64("pw", 0.001, "Pulse Width")
	flag.Parse()
	fmt.Println(*channel, *frequency, *pulseWidth)
	i2cConnection, err := i2c.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println(err)
		return
	}
	pca9685, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection)
	esc := esc.NewESC(pca9685)
	if err != nil {
		fmt.Println(err)
		return
	}

	esc.Start(float32(*frequency))
	defer esc.Close()

	loop := 20000
	start := time.Now()
	for i := 1; i < loop; i++ {
		esc.SetPulseWidth(0, float32(*pulseWidth))
		esc.SetPulseWidth(1, float32(*pulseWidth))
		esc.SetPulseWidth(2, float32(*pulseWidth))
		esc.SetPulseWidth(3, float32(*pulseWidth))
	}
	elapsed := time.Since(start)
	fmt.Println(float64(loop) / elapsed.Seconds())

}
