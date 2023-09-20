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

func addButton(tag string, index int, gpioPin string, buttons map[string]pushButton) {
	name := fmt.Sprintf("%s-%d", tag, index)
	pin := hardware.NewPullDownPushButton(gpioPin)
	buttons[name] = pushbutton.NewPushButton(name, pin)
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	configs := remote.ReadConfigs("./configs/remote-configs.json")
	buttons := make(map[string]pushButton)
	for i := 0; i < len(configs.PushButtons.Left); i++ {
		addButton("left", i, configs.PushButtons.Left[i], buttons)
		addButton("right", i, configs.PushButtons.Right[i], buttons)
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
			count := 0
			for running {
				select {
				case _, ok := <-pressed:
					if ok {
						count++
						log.Printf("%s pressed  %3d\n", pb.Name(), count)
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
