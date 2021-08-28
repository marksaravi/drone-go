package flightcontrol

import "fmt"

type flightControl struct {
}

func NewFlightControl() *flightControl {
	return &flightControl{}
}

func (fc *flightControl) Start() {
	fmt.Println("Starting Flight Control")
}
