package remote

import (
	"context"
	"fmt"
	"time"
)

type radioTransmiter interface {
	On()
	Transmit(payload []byte) error
}

type joystick interface {
	Read() uint16
}

type PushButton interface {
	Name()      string
	PulseMode() bool
	IsPressed() bool
}

type commands struct {
	roll     uint16
	pitch    uint16
	yaw      uint16
	throttle uint16
}

type remoteControl struct {
	transmitter      radioTransmiter
	roll             joystick
	pitch            joystick
	yaw              joystick
	throttle         joystick
	buttons          []PushButton
	commandPerSecond int
	buttonsPressed   []byte
	commands         commands
}

type RemoteSettings struct {
	Transmitter                radioTransmiter
	CommandPerSecond           int
	Roll, Pitch, Yaw, Throttle joystick
	PushButtons                []PushButton
}

func NewRemoteControl(settings RemoteSettings) *remoteControl {
	return &remoteControl{
		transmitter:      settings.Transmitter,
		commandPerSecond: settings.CommandPerSecond,
		roll:             settings.Roll,
		pitch:            settings.Pitch,
		yaw:              settings.Yaw,
		throttle:         settings.Throttle,
		buttons:          settings.PushButtons,
		buttonsPressed:   make([]byte, len(settings.PushButtons)),
	}
}

func (r *remoteControl) Start(ctx context.Context) {
	running := true
	r.transmitter.On()
	lastCommand:=time.Now()
	commandTimeout:=time.Second/time.Duration(r.commandPerSecond)
	for running {
		select {
		default:
			if time.Since(lastCommand)>=commandTimeout {
				r.ReadCommands()
				r.ReadButtons()
				continuesOutputButtons, pulseOutputButtons := r.PushButtonsPayloads()
				payload:= []byte {
					byte(r.commands.roll),
					byte(r.commands.pitch),
					byte(r.commands.yaw),
					byte(r.commands.throttle),
					continuesOutputButtons,
					pulseOutputButtons,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				}
				fmt.Println(payload)
				r.transmitter.Transmit(payload)
				
				lastCommand=time.Now()
			}
		case <-ctx.Done():
			running = false
		}
	}
}

func (r *remoteControl) PushButtonsPayloads() (byte, byte) {
	continuesOutputs:=byte(0)
	coshift:=0
	pulseOutputs:=byte(0)
	pshift:=0
	for i, bp := range r.buttonsPressed {
		if r.buttons[i].PulseMode() {
			pulseOutputs |= bp << pshift
			pshift++
		} else {
			continuesOutputs |= bp << coshift
			coshift++
		}
	}
	return continuesOutputs, pulseOutputs
}