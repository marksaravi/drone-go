package hardware

// import (
// 	"log"
// )

// type gpio interface{}
// type gpioreg interface{}
// type i2c interface{}
// type i2creg interface{}

// type gpioswitch struct {
// 	pin gpio.PinIn
// }

// func (b *gpioswitch) Read() bool {
// 	return b.pin.Read() == gpio.Low
// }

// func NewGPIOSwitch(pinName string) *gpioswitch {
// 	var pin gpio.PinIn = gpioreg.ByName(pinName)
// 	if pin == nil {
// 		log.Fatal("Failed to find ", pinName)
// 	}
// 	pin.In(gpio.Float, gpio.NoEdge)
// 	return &gpioswitch{
// 		pin: pin,
// 	}
// }
