package esc

import "github.com/MarkSaravi/drone-go/types"

type esc interface {
	Start(frequency float32) error
	SetPulseWidth(channel int, pulseWidth float32)
	SetPulseWidthAll(pulseWidth float32)
	StopAll()
	Halt() error
	Close()
}

type escsHandler struct {
	esc       esc
	throttles []types.Throttle
}

func (h *escsHandler) SetThrottles(throttles []types.Throttle) {
	h.throttles = throttles
	for _, throttle := range h.throttles {
		h.esc.SetPulseWidth(throttle.Motor, throttle.Value)
	}
}

func NewESCsHandler() *escsHandler {
	return &escsHandler{}
}
