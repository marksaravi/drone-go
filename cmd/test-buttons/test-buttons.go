package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/apps/remote"
	pushbutton "github.com/marksaravi/drone-go/devices/push-button"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/host/v3"
)

type btn interface {
	Start(ctx context.Context) <-chan bool
	Name() string
	GPIO() string
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := driverreg.Init(); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(state)
	configs := remote.ReadConfigs("./configs.yaml")
	buttons := make(map[string]btn)
	for _, btn := range configs.Buttons {
		fmt.Println(btn.Name)
		buttons[btn.Name] = pushbutton.NewPushButton(btn.GPIO, btn.Name, true)
	}
	fmt.Println()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		fmt.Scanln()
		cancel()
	}()

	signals := make(map[string]<-chan bool)
	btnl0 := true
	for _, btn := range configs.Buttons {
		signals[btn.Name] = buttons[btn.Name].Start(ctx)
	}

	prev := time.Now()
	for btnl0 {
		select {
		case _, ok := <-signals["btn-l0"]:
			if ok {
				log.Println("signal btn-l0", time.Since(prev).Milliseconds())
				prev = time.Now()
			} else {
				fmt.Println("closed btn-l0")
				btnl0 = false
			}

		default:
		}
	}
}
