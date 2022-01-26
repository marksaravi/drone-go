package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/pca9685"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.InitHost()
	motor := flag.Int("motor", 0, "motor")
	flag.Parse()

	configs := config.ReadConfigs().FlightControl
	b, _ := i2creg.Open(configs.ESC.I2CDev)
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	pwmDev, _ := pca9685.NewPCA9685(pca9685.PCA9685Settings{
		Connection:      i2cConn,
		MaxThrottle:     configs.MaxThrottle,
		ChannelMappings: configs.ESC.PwmDeviceToESCMappings,
	})

	pwmDev.SetThrottle(*motor, 0, true)
	fmt.Scanln()
	log.Println("finished")
}
