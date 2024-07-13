package remote

import (
	"context"
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/apps/commons"
	"github.com/marksaravi/drone-go/utils"
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
	Update()
	IsPressed() bool
	IsPushed() bool
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
	displayUpdate := utils.WithDataPerSecond(3)
	for running {
		select {
		default:
			r.updateButtons()
			if r.readCommands() {
				pressedButtons, pushButtons := r.PushButtonsPayloads()
				lRoll, hRoll := commons.Uint16ToBytes(r.commands.roll)
				lPitch, hPitch := commons.Uint16ToBytes(r.commands.pitch)
				lYaw, hYaw := commons.Uint16ToBytes(r.commands.yaw)
				lThrottle, hThrottle := commons.Uint16ToBytes(r.commands.throttle)
				payload := []byte{
					lRoll,
					hRoll,
					lPitch,
					hPitch,
					lYaw,
					hYaw,
					lThrottle,
					hThrottle,
					pressedButtons,
					pushButtons,
				}
				r.transmitter.Transmit(payload)
				if displayUpdate.IsTime() {
					r.UpdateDisplay(payload)
				}
			}
		case <-ctx.Done():
			running = false
		}
	}
}

func (r *remoteControl) readCommands() bool {
	if time.Since(r.lastCommandRead) < time.Second/time.Duration(r.commandPerSecond) {
		return false
	}
	r.lastCommandRead = time.Now()
	r.readJoysticks()
	r.readButtons()
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
	r.oled.WriteString(fmt.Sprintf("%2.1f%%", commons.CalcThrottleFromRawJoyStickRaw(payload[6:8], 100)), 8, 0)
}

func (r *remoteControl) PushButtonsPayloads() (pressedButtons byte, pushButtons byte) {
	pressedButtons = byte(0)
	pushButtons = byte(0)
	for i, button := range r.buttons {
		if button.IsPushed() {
			pushButtons = pushButtons | (byte(1)<<i)  
		}			
		if button.IsPressed() {
			pressedButtons = pressedButtons | (byte(1)<<i)  
		}
	}
	return
}

func (r *remoteControl) readJoysticks() {
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

func (r *remoteControl) readButtons() {

}

func (r *remoteControl) updateButtons() {
	for _, button := range r.buttons {
		button.Update()
	}
}
