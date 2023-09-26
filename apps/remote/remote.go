package remote

import (
	"context"
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/constants"
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

type oled interface {
	WriteString(string, int, int)
}

type remoteControl struct {
	transmitter            radioTransmiter
	roll                   joystick
	pitch                  joystick
	yaw                    joystick
	throttle               joystick
	buttons                []PushButton
	oled                   oled
	commandPerSecond       int
	lastCommandRead        time.Time
	buttonsPressed         []byte
	commands               commands
	displayUpdatePerSecond int
	lastDisplayUpdate      time.Time
}

type RemoteSettings struct {
	Transmitter                radioTransmiter
	CommandPerSecond           int
	Roll, Pitch, Yaw, Throttle joystick
	PushButtons                []PushButton
	OLED                       oled
	DisplayUpdatePerSecond     int
}

func NewRemoteControl(settings RemoteSettings) *remoteControl {
	return &remoteControl{
		transmitter:            settings.Transmitter,
		commandPerSecond:       settings.CommandPerSecond,
		roll:                   settings.Roll,
		pitch:                  settings.Pitch,
		yaw:                    settings.Yaw,
		throttle:               settings.Throttle,
		buttons:                settings.PushButtons,
		oled:                   settings.OLED,
		displayUpdatePerSecond: settings.DisplayUpdatePerSecond,
		buttonsPressed:         make([]byte, len(settings.PushButtons)),
		lastCommandRead:        time.Now(),
		lastDisplayUpdate:      time.Now(),
	}
}

func (r *remoteControl) Start(ctx context.Context) {
	running := true
	r.transmitter.On()
	r.Initisplay()
	for running {
		select {
		default:
			if r.ReadCommands() {
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
				fmt.Println(payload, r.JoystickToString())
				r.transmitter.Transmit(payload)
				r.UpdateDisplay()
			}
		case <-ctx.Done():
			running = false
		}
	}
}

func (r *remoteControl) ReadCommands() bool {
	if time.Since(r.lastCommandRead)<time.Second/time.Duration(r.commandPerSecond) {
		return false
	}
	r.lastCommandRead=time.Now()
	r.ReadJoysticks()
	r.ReadButtons()
	return true
}

func (r *remoteControl) Initisplay() {
	r.oled.WriteString("Trottle:", 0, 0)
}

func (r *remoteControl) UpdateDisplay() {
	if time.Since(r.lastDisplayUpdate)<time.Second/time.Duration(r.displayUpdatePerSecond) {
		return
	}
	r.lastDisplayUpdate=time.Now()
	r.oled.WriteString(" ", 13, 0)
	r.oled.WriteString(fmt.Sprintf("%2.1f%%", r.Throttle()), 9, 0)
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

func (r *remoteControl) ReadJoysticks() {
	roll := r.roll.Read()
	pitch := r.pitch.Read()
	yaw := r.yaw.Read()
	throttle := r.throttle.Read()

	r.commands = commands{
		roll:     roll,
		pitch:    pitch,
		yaw:      yaw,
		throttle: throttle,
	}
}

func (r *remoteControl) ReadButtons() {
	for i, button := range r.buttons {
		pressed:=button.IsPressed()
		if pressed {
			r.buttonsPressed[i]=byte(1)
		} else {
			r.buttonsPressed[i]=byte(0)
		}
	}
}

func (r *remoteControl) Throttle() float32 {
	return float32(r.commands.throttle)/constants.THROTTLE_MAX*100
}



func jsticktofloat(x uint16) float32 {
	const OUTPUT_RANGE_DEG = float32(90)
	return OUTPUT_RANGE_DEG/2 - float32(x)*OUTPUT_RANGE_DEG/constants.JOYSTICK_RANGE_DEG
}

func (r *remoteControl) Roll() float32 {
	return jsticktofloat(r.commands.roll)
}

func (r *remoteControl) Pitch() float32 {
	return jsticktofloat(r.commands.pitch)
}

func (r *remoteControl) Yaw() float32 {
	return jsticktofloat(r.commands.yaw)
}

func (r *remoteControl) JoystickToString() string {
	return fmt.Sprintf("%2.1f, %2.1f, %2.1f, %2.1f%%", r.Roll(), r.Pitch(), r.Yaw(), r.Throttle())
}