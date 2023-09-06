package main

import (
	"context"
	"fmt"
	"log"

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
	for _, btn := range configs.Buttons {
		signals[btn.Name] = buttons[btn.Name].Start(ctx)
	}
	btnl0 := true
	btnl1 := true
	btnl2 := true
	btnl3 := true
	btnl4 := true
	btnr0 := true
	btnr1 := true
	btnr2 := true
	btnr3 := true
	btnr4 := true

	for btnl0 || btnl1 || btnl2 || btnl3 || btnl4 || btnr0 || btnr1 || btnr2 || btnr3 || btnr4 {
		select {
		case _, ok := <-signals["btn-l0"]:
			if ok {
				log.Println("signal btn-l0")

			} else {
				fmt.Println("closed btn-l0")
				btnl0 = false
			}
		case _, ok := <-signals["btn-l1"]:
			if ok {
				log.Println("signal btn-l1")

			} else {
				fmt.Println("closed btn-l1")
				btnl1 = false
			}
		case _, ok := <-signals["btn-l2"]:
			if ok {
				log.Println("signal btn-l2")

			} else {
				fmt.Println("closed btn-l2")
				btnl2 = false
			}
		case _, ok := <-signals["btn-l3"]:
			if ok {
				log.Println("signal btn-l3")

			} else {
				fmt.Println("closed btn-l3")
				btnl3 = false
			}
		case _, ok := <-signals["btn-l4"]:
			if ok {
				log.Println("signal btn-l4")

			} else {
				fmt.Println("closed btn-l4")
				btnl4 = false
			}
		case _, ok := <-signals["btn-r0"]:
			if ok {
				log.Println("signal btn-r0")

			} else {
				fmt.Println("closed btn-r0")
				btnr0 = false
			}
		case _, ok := <-signals["btn-r1"]:
			if ok {
				log.Println("signal btn-r1")

			} else {
				fmt.Println("closed btn-r1")
				btnr1 = false
			}
		case _, ok := <-signals["btn-r2"]:
			if ok {
				log.Println("signal btn-r2")

			} else {
				fmt.Println("closed btn-r2")
				btnr2 = false
			}
		case _, ok := <-signals["btn-r3"]:
			if ok {
				log.Println("signal btn-r3")

			} else {
				fmt.Println("closed btn-r3")
				btnr3 = false
			}
		case _, ok := <-signals["btn-r4"]:
			if ok {
				log.Println("signal btn-r4")

			} else {
				fmt.Println("closed btn-r4")
				btnr4 = false
			}
		default:
		}
	}
}
