package remote

import (
	"context"
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/apps/common"
)

type radioTransmiter interface {
	On()
	Transmit(payload []byte) error
}

type joystick interface {
	Read() uint16
}

type PushButton interface {
	Name() string
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
	rollMidValue           int
	pitchMidValue          int
	yawMidValue            int
	rotationRange          float64
	maxThrottle            float64
}

type RemoteSettings struct {
	Transmitter                radioTransmiter
	CommandPerSecond           int
	Roll, Pitch, Yaw, Throttle joystick
	PushButtons                []PushButton
	OLED                       oled
	DisplayUpdatePerSecond     int
	RollMidValue               int
	PitchMidValue              int
	YawMidValue                int
	RotationRange              float64
	MaxThrottle                float64
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
		rollMidValue:           settings.RollMidValue,
		pitchMidValue:          settings.PitchMidValue,
		yawMidValue:            settings.YawMidValue,
		rotationRange:          settings.RotationRange,
		maxThrottle:            settings.MaxThrottle,
	}
}

func (r *remoteControl) Start(ctx context.Context) {
	running := true
	r.transmitter.On()
	r.Initisplay()
	prevPayload := make([]byte, 10)
	for running {
		select {
		default:
			if r.ReadCommands() {
				continuesOutputButtons, pulseOutputButtons := r.PushButtonsPayloads()
				hRoll, lRoll := Uint16ToBytes(r.commands.roll)
				hPitch, lPitch := Uint16ToBytes(r.commands.pitch)
				hYaw, lYaw := Uint16ToBytes(r.commands.yaw)
				hThrottle, lThrottle := Uint16ToBytes(r.commands.throttle)
				payload := []byte{
					hRoll,
					lRoll,
					hPitch,
					lPitch,
					hYaw,
					lYaw,
					hThrottle,
					lThrottle,
					continuesOutputButtons,
					pulseOutputButtons,
				}
				if isChanged(payload, prevPayload) {
					fmt.Println(payload)
				}
				copy(prevPayload, payload)
				r.transmitter.Transmit(payload)
				r.UpdateDisplay(payload)
			}
		case <-ctx.Done():
			running = false
		}
	}
}

func (r *remoteControl) ReadCommands() bool {
	if time.Since(r.lastCommandRead) < time.Second/time.Duration(r.commandPerSecond) {
		return false
	}
	r.lastCommandRead = time.Now()
	r.ReadJoysticks()
	r.ReadButtons()
	return true
}

func (r *remoteControl) Initisplay() {
	r.oled.WriteString("Trottle:", 0, 0)
}

func (r *remoteControl) UpdateDisplay(payload []byte) {
	if time.Since(r.lastDisplayUpdate) < time.Second/time.Duration(r.displayUpdatePerSecond) {
		return
	}
	r.lastDisplayUpdate = time.Now()
	r.oled.WriteString(" ", 13, 0)
	r.oled.WriteString(fmt.Sprintf("%2.1f%%", common.CalcThrottleFromRawJoyStickRaw(payload[6:8], r.maxThrottle)), 8, 0)
}

func (r *remoteControl) PushButtonsPayloads() (byte, byte) {
	continuesOutputs := byte(0)
	coshift := 0
	pulseOutputs := byte(0)
	pshift := 0
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
		pressed := button.IsPressed()
		if pressed {
			r.buttonsPressed[i] = byte(1)
		} else {
			r.buttonsPressed[i] = byte(0)
		}
	}
}

func Uint16ToBytes(x uint16) (high, low byte) {
	low = byte(x)
	high = byte(x >> 8)
	return
}

var counter int = 0

func isChanged(payload, prevPayload []byte) bool {
	if payload[9] != prevPayload[9] || payload[8] != prevPayload[8] || counter > 20 {
		counter = 0
		return true
	}
	counter++
	return false
}
