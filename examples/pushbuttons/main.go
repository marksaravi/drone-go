package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/marksaravi/drone-go/apps/remote"
	pushbutton "github.com/marksaravi/drone-go/devices/push-button"
	"github.com/marksaravi/drone-go/hardware"
)

type pushButton interface {
	Name()         string
	Index()        int
	IsPressed()    bool
	IsPushed()     bool
	IsPushButton() bool
	Update()
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	pushTest := true
	if len(os.Args)>1 {
		pushTest = os.Args[1] != "press"
		fmt.Println(pushTest)
	}

	log.Println("Starting the test...")
	hardware.HostInitialize()
	configs := remote.ReadConfigs("./configs/remote-configs.json")
	fmt.Println(configs.PushButtons)
	buttons := make([]pushButton, 0, 10)
	buttonsCount:=make([]int,0 , 10)
	for i := 0; i < len(configs.PushButtons); i++ {
		pin:=hardware.NewPushButtonInput(configs.PushButtons[i].GPIO)
		buttons = append(buttons, pushbutton.NewPushButton(configs.PushButtons[i].Name, configs.PushButtons[i].Index, configs.PushButtons[i].IsPushButton, pin))
		buttonsCount=append(buttonsCount, 0)
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
		for _, button := range buttons {
			button.Update()
			if pushTest {
				if button.IsPushed() && button.IsPushButton() {
					log.Printf("%s %d pushed\n", button.Name(), button.Index())
				}
			} else {
				if button.IsPressed() && !button.IsPushButton() {
					log.Printf("%s %d pressed\n", button.Name(), button.Index())
				}
			}
		}
	}
}
