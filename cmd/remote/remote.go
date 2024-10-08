package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/marksaravi/drone-go/apps/drone"
	"github.com/marksaravi/drone-go/apps/remote"
	pushbutton "github.com/marksaravi/drone-go/devices/push-button"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/ads1115"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
	"github.com/marksaravi/drone-go/hardware/ssd1306"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func main() {
	runAsService := false
	if len(os.Args) > 1 {
		runAsService = os.Args[1] == "run-as-service"
		if runAsService {
			fmt.Println("Running as Service")
		}
	}
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	log.Println("Starting RemoteControl")
	configs := remote.ReadConfigs("./configs/remote-configs.json")
	droneConfigs := drone.ReadConfigs("./configs/drone-configs.json")
	log.Println(configs)

	radioConfigs := configs.Radio
	radioLink := nrf24l01.NewNRF24L01EnhancedBurst(
		radioConfigs.SPI,
		radioConfigs.RxTxAddress,
	)

	radioTransmitter := radio.NewRadioTransmitter(radioLink)

	buttons := make([]remote.PushButton, 0, 10)
	buttonsCount := make([]int, 0, 10)
	for i := 0; i < len(configs.PushButtons); i++ {
		pin := hardware.NewPushButtonInput(configs.PushButtons[i].GPIO)
		buttons = append(buttons, pushbutton.NewPushButton(configs.PushButtons[i].Name, configs.PushButtons[i].Index, configs.PushButtons[i].IsPushButton, pin))
		buttonsCount = append(buttonsCount, 0)
	}

	i2cbus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer i2cbus.Close()
	displayi2c := &i2c.Dev{Addr: 0x3D, Bus: i2cbus}
	atodi2c := &i2c.Dev{Addr: 0x48, Bus: i2cbus}
	oled := ssd1306.NewSSD1306(displayi2c, ssd1306.DefaultOptions)
	err = oled.Init()
	if err != nil {
		log.Fatal(err)
	}
	atod := ads1115.NewADS1115(atodi2c)
	fmt.Println(configs.Joysticks.RollMin, configs.Joysticks.RollMid, configs.Joysticks.RollMax)
	remoteControl := remote.NewRemoteControl(remote.RemoteSettings{
		Transmitter:            radioTransmitter,
		CommandPerSecond:       configs.CommandsPerSecond,
		JoyStick:               atod,
		Roll:                   0,
		Pitch:                  1,
		Yaw:                    3,
		Throttle:               2,
		PushButtons:            buttons,
		OLED:                   oled,
		DisplayUpdatePerSecond: configs.DisplayUpdatePerSecond,
		RollMin:                configs.Joysticks.RollMin,
		RollMid:                configs.Joysticks.RollMid,
		RollMax:                configs.Joysticks.RollMax,
		PitchMin:               configs.Joysticks.PitchMin,
		PitchMid:               configs.Joysticks.PitchMid,
		PitchMax:               configs.Joysticks.PitchMax,
		YawMin:                 configs.Joysticks.YawMin,
		YawMid:                 configs.Joysticks.YawMid,
		YawMax:                 configs.Joysticks.YawMax,
		ThrottleMin:            configs.Joysticks.ThrottleMin,
		ThrottleMid:            configs.Joysticks.ThrottleMid,
		ThrottleMax:            configs.Joysticks.ThrottleMax,
		RotationRange:          droneConfigs.Commands.RotationRange,
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if !runAsService {
			fmt.Scanln()
			cancel()
		}
	}()

	remoteControl.Start(ctx)
}
