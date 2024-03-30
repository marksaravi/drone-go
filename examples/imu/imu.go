package main

import (
	"context"
	"fmt"
	"log"
	"time"

	plotterclient "github.com/marksaravi/drone-go/apps/plotter-client"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20789"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	plotterClient := plotterclient.NewPlotter(plotterclient.Settings{
		Active:  false,
		Address: "192.168.1.101:8000",
	})
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		fmt.Scanln()
		cancel()
	}(cancel)
	icm20789Configs := icm20789.ReadConfigs("./configs/hardware.json")
	mems := icm20789.NewICM20789(icm20789Configs)
	whoAmI, err := mems.WhoAmI()
	if err == nil {
		fmt.Printf("WHO AM I: %x\n", whoAmI)
	}
	imuConfigs := imu.ReadConfigs("./configs/drone-configs.json")
	fmt.Println(imuConfigs)
	imudev := imu.NewIMU(mems, imuConfigs)
	running := true

	plotterClient.SetStartTime(time.Now())
	printInterval := utils.WithDataPerSecond(4)
	readInterval := utils.WithMinInterval(imuConfigs.DataPerSecond, 1200)
	log.Println("Starting...")
	for running {
		select {
		case <-ctx.Done():
			running = false
		default:
			if readInterval.IsTime() {
				rot, acc, gyro, err := imudev.ReadAll()
				if err == nil {
					plotterClient.SendPlotterData(rot, acc, gyro)
				}
				if printInterval.IsTime() {
					printRotations(rot, acc, gyro)
				}
			}
		}
	}
}

func printRotations(rot, acc, gyro imu.Rotations) {
	// fmt.Printf("Roll: %6.2f, Pitch: %6.2f, Yaw: %6.2f,  Acc Roll: %6.2f, Pitch: %6.2f,  Gyro Roll: %6.2f, Pitch: %6.2f, Yaw: %6.2f\n", rot.Roll, rot.Pitch, rot.Yaw, acc.Roll, acc.Pitch, gyro.Roll, gyro.Pitch, gyro.Yaw)
	// fmt.Printf("Roll: %6.2f, Pitch: %6.2f\n", acc.Roll, acc.Pitch)
	// fmt.Printf("Acc Roll: %6.2f, Pitch: %6.2f,  Gyro Roll: %6.2f, Pitch: %6.2f, Yaw: %6.2f\n", acc.Roll, acc.Pitch, gyro.Roll, gyro.Pitch, gyro.Yaw)
	// fmt.Printf("Roll(%6.2f %6.2f %6.2f),    Pitch(%6.2f  %6.2f %6.2f), Yaw(%6.2f)\n", acc.Roll, gyro.Roll, rot.Roll, acc.Pitch, gyro.Pitch, rot.Pitch, gyro.Yaw)
	fmt.Printf("GYRO: %6.2f  %6.2f %6.2f\n", gyro.Roll, gyro.Pitch, gyro.Yaw)
	fmt.Printf(" ACC: %6.2f  %6.2f %6.2f\n", acc.Roll, acc.Pitch, acc.Yaw)
}
