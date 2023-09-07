package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/apps/remote"
	pushbutton "github.com/marksaravi/drone-go/devices/push-button"
	"github.com/marksaravi/drone-go/hardware"
)

type pushButton interface {
	Start(ctx context.Context) <-chan bool
	Name() string
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	configs := remote.ReadConfigs("./remote-configs.yaml")
	buttons := make(map[string]pushButton)
	for name, gpioPin := range configs.Buttons {
		fmt.Printf("%s:%s\n", name, gpioPin)
		pin := hardware.NewPullupPushButton(gpioPin)
		buttons[name] = pushbutton.NewPushButton(name, pin)
	}
	fmt.Println()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	go func() {
		fmt.Scanln()
		cancel()
	}()

	for name, button := range buttons {
		wg.Add(1)
		go func(btnName string, pb pushButton) {
			defer wg.Done()
			running := true
			pressed := pb.Start(ctx)
			for running {
				select {
				case _, ok := <-pressed:
					if ok {
						log.Printf("%s pressed\n", pb.Name())
					} else {
						running = false
						log.Printf("%s closed\n", pb.Name())
					}
				}
			}
		}(name, button)
	}
	wg.Wait()
}
