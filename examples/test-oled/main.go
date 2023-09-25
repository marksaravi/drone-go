package main

import (
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/ssd1306"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func main() {
	hardware.HostInitialize()
	b, err := i2creg.Open("")
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
	max:=float32(100)
	min:=float32(0)
	dt:=float32(0.1)
	throttle:=float32(0)
	for {
		time.Sleep(time.Second/4)
		oled.WriteString("            ", 0, 2)
		s:=fmt.Sprintf("T: %0.1f%s", throttle, "%")
		fmt.Println(s)
		oled.WriteString(s, 0, 2)
		throttle += dt
		if throttle>max || throttle <min {
			dt = -dt
		}
	}
}
