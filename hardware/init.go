package hardware

import (
	"log"
)

type host interface{}

func InitHost() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
}
