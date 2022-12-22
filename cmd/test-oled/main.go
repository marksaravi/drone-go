package main

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/hardware/ssd1306"
)

type host interface{}
type gpio interface{}
type gpioreg interface{}
type i2c interface{}
type i2creg interface{}

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()
	d := &i2c.Dev{Addr: 0x3D, Bus: b}
	oled := ssd1306.NewSSD1306(d, ssd1306.DefaultOptions)
	err = oled.Init()
	if err != nil {
		log.Fatal(err)
	}
	oled.Println("!000000000000!", 0)
	time.Sleep(time.Second)
	oled.Println("Hello Mark!", 0)
	oled.WriteString("Disconnected", 0, 1)
	oled.WriteString("T: 15.7%", 0, 2)
}
