package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/hardware/nrf204"
	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/sysfs"
)

func main() {
	fa := []float32{366.34, -180.24, 0, -144.32, 22.22}
	ba := utils.FloatArrayToByteArray(fa)
	fa2 := utils.ByteArrayToFloat32Array(ba)
	fmt.Println(ba)
	fmt.Println(fa2)
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
	for range time.Tick(time.Millisecond * 1000) {
		fmt.Println("send ", payload[0])
		err := receiver.WritePayload(payload)
		if err != nil {
			fmt.Println(err)
		}
		payload[0] = payload[0] + 1
	}
}
