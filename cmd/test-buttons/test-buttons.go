package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/apps/remote"
	"github.com/marksaravi/drone-go/devices/button"
	"periph.io/x/host/v3"
)

type btn interface {
	Read() (level bool, pressed bool)
	Name() string
}

func main() {
	state, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(state)
	configs := remote.ReadConfigs("./configs.yaml")
	buttons := make([]btn, len(configs.Buttons))
	for i, btn := range configs.Buttons {
		fmt.Println(btn.Name)
		buttons[i] = button.NewButton(btn.Name)
	}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		fmt.Scanln()
		cancel()
	}()

	running := true
	for running {
		select {
		case <-ctx.Done():
			running = false
		default:
			time.Sleep(20 * time.Millisecond)
		}
		for i := 0; i < len(buttons); i++ {
			level, pressed := buttons[i].Read()
			if level && pressed {
				fmt.Println(buttons[i].Name())
			}
		}
	}
}
