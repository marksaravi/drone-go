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
	Read(channel int) (l, h byte)
}

type PushButton interface {
	Name()         string
	Index()        int
	Update()
	IsPressed()    bool
	IsPushed()     bool
	IsPushButton() bool
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
	joyStick               joystick
	rollChan               int
	pitchChan              int
	yawChan                int
	throttleChan           int
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
	JoyStick                   joystick
	Roll, Pitch, Yaw, Throttle int
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
		joyStick:               settings.JoyStick,
		rollChan:               settings.Roll,
		pitchChan:              settings.Pitch,
		yawChan:                settings.Yaw,
		throttleChan:               settings.Throttle,
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
	commandsUpdate := utils.WithDataPerSecond(r.commandPerSecond)
	fmt.Println("Commands per second: ", r.commandPerSecond)
	for running {
		select {
		default:
			r.updateButtons()
			if commandsUpdate.IsTime() {
				pressedButtons, pushButtons := r.readButtons()
				lRoll, hRoll := r.joyStick.Read(r.rollChan)
				lPitch, hPitch := r.joyStick.Read(r.pitchChan)
				lYaw, hYaw := byte(0), byte(0) //r.joyStick.Read(r.Chan)
				lThrottle, hThrottle := r.joyStick.Read(r.throttleChan)
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

func (r *remoteControl) Initisplay() {
	r.oled.WriteString("Trottle:", 0, 0)
}

func (r *remoteControl) UpdateDisplay(payload []byte) {
	if time.Since(r.lastDisplayUpdate) < time.Second/time.Duration(r.displayUpdatePerSecond) {
		return
	}
	r.lastDisplayUpdate = time.Now()
	r.oled.WriteString(" ", 13, 0)
	r.oled.WriteString(fmt.Sprintf("%2.1f%%", commons.CalcThrottleFromRawJoyStickRaw(payload[6:8], 100)), 8, 53)
}

func (r *remoteControl) readButtons() (pressedButtons byte, pushButtons byte) {
	pressedButtons = byte(0)
	pushButtons = byte(0)
	for _, button := range r.buttons {
		if button.IsPushed() && button.IsPushButton() {
			pushButtons = pushButtons | (byte(1)<<button.Index())  
		}			
		if button.IsPressed() && !button.IsPushButton() {
			pressedButtons = pressedButtons | (byte(1)<<button.Index())  
		}
	}
	return
}

func (r *remoteControl) updateButtons() {
	for _, button := range r.buttons {
		button.Update()
	}
}

func toInt(l, h byte) int {
	return int(uint16(l) | (uint16(h) << 8))
}