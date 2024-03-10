package drone

import "fmt"

type commandHandler struct {
	escs escs
}

func NewCommandHandler(escs escs) *commandHandler {
	return &commandHandler{
		escs: escs,
	}
}

func (h *commandHandler) applyCommands(commands []byte) {
	h.offMotors(commands)
	h.onMotors(commands)
}

func (h *commandHandler) onMotors(commands []byte) {
	if commands[5] == 1 {
		h.escs.On()
		fmt.Println("ESC CONNECTED")
	}
}

func (h *commandHandler) offMotors(commands []byte) {
	if commands[5] == 16 {
		h.escs.Off()
		fmt.Println("ESC DISCONNECTED")
	}
}
