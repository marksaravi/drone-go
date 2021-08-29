package drivers

import (
	"log"

	"periph.io/x/periph/host"
)

func InitHost() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
}
