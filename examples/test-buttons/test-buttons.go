package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/apps/remote"
	pushbutton "github.com/marksaravi/drone-go/devices/push-button"
	"github.com/marksaravi/drone-go/hardware"
)

type pushButton interface {
	Name() string
	Read() bool
}

// func addButton(tag string, index int, gpioPin string, buttons []pushButton, buttonsStates map[string]bool) []pushButton {
// 	name := fmt.Sprintf("%s-%d", tag, index)
// 	pin:=hardware.NewPushButtonInput(gpioPin)
// 	buttons = append(buttons, pushbutton.NewPushButton(name, pin))
// 	buttonsStates[name]=false
// 	return buttons
// }

func main() {
	log.SetFlags(log.Lmicroseconds)
	log.Println("Starting the test...")
	hardware.HostInitialize()
	configs := remote.ReadConfigs("./configs/remote-configs.json")
	fmt.Println(configs.PushButtons)
	buttons := make([]pushButton, 0, 10)
	buttonsStates:=make(map[string]bool)
	buttonsCount:=make([]int,0 , 10)
	for i := 0; i < len(configs.PushButtons); i++ {
		pin:=hardware.NewPushButtonInput(configs.PushButtons[i].GPIO)
		buttons = append(buttons, pushbutton.NewPushButton(configs.PushButtons[i].Name, pin))
		buttonsCount=append(buttonsCount, 0)
		buttonsStates[configs.PushButtons[i].Name]=true
	}
	fmt.Println()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		fmt.Scanln()
		cancel()
	}()

	const DATA_PER_SECOND = 25
	timeout:=time.Second/DATA_PER_SECOND

	lastRead:=time.Now()
	running:=true
	for running {
		select {
		case <-ctx.Done():
			running=false
		default:
		}
		if time.Since(lastRead)<timeout {
			continue
		}
		lastRead=time.Now()
		for i, button := range buttons {
			pressed:=button.Read()
			if pressed!=buttonsStates[button.Name()] {
				if pressed {
					buttonsCount[i]++
					log.Printf("%s pressed  (%3d)\n", button.Name(), buttonsCount[i])
				}
				buttonsStates[button.Name()]=pressed
			}
		}
	}
}
