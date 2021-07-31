package main

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/MarkSaravi/drone-go/hardware/nrf204"
	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/sysfs"
)

func byteArrayToInt16Array(ba []byte, size int) []int16 {
	type pInt16Array = *([]int16)
	var pi16 pInt16Array = pInt16Array(unsafe.Pointer(&ba))
	var ia []int16 = make([]int16, size/2)
	for i := 0; i < size/2; i++ {
		ia[i] = (*pi16)[i]
	}
	return ia
}

func main() {
	config := types.RadioLinkConfig{
		GPIO: types.RadioLinkGPIOPins{
			CE: "GPIO26",
		},
		BusNumber:  1,
		ChipSelect: 2,
		RxAddress:  "03896",
		PowerDBm:   nrf204.RF_POWER_MINUS_18dBm,
	}
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	spibus, err := sysfs.NewSPI(config.BusNumber, config.ChipSelect)
	if err != nil {
		log.Fatal(err)
	}
	spiconn, err := spibus.Connect(physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Start")
	receiver := nrf204.CreateNRF204(config, spiconn)
	receiver.Init()
	receiver.StartTransmitting()
	payload := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1}
	for range time.Tick(time.Millisecond * 10) {
		fmt.Println("send ", payload[0])
		err := receiver.WritePayload(payload)
		if err != nil {
			fmt.Println(err)
		}
		payload[0] = payload[0] + 1
	}
}
