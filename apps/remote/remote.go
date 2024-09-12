package remote

import (
	"context"
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/constants"
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
	rollMin                    uint16
	pitchMin                   uint16
	yawMin                     uint16
	throttleMin                uint16
	rollMid                    uint16
	pitchMid                   uint16
	yawMid                     uint16
	throttleMid                uint16
	rollMax                    uint16
	pitchMax                   uint16
	yawMax                     uint16
	throttleMax                uint16
	rotationRange          float64
}

type RemoteSettings struct {
	Transmitter                radioTransmiter
	CommandPerSecond           int
	JoyStick                   joystick
	Roll, Pitch, Yaw, Throttle int
	PushButtons                []PushButton
	OLED                       oled
	DisplayUpdatePerSecond     int
	RollMin                    uint16
	PitchMin                   uint16
	YawMin                     uint16
	ThrottleMin                uint16
	RollMid                    uint16
	PitchMid                   uint16
	YawMid                     uint16
	ThrottleMid                uint16
	RollMax                    uint16
	PitchMax                   uint16
	YawMax                     uint16
	ThrottleMax                uint16
	RotationRange              float64
}

func NewRemoteControl(settings RemoteSettings) *remoteControl {
	rmin,rmid,rmax:=fixMax(settings.RollMin,settings.RollMid,settings.RollMax)
	pmin,pmid,pmax:=fixMax(settings.PitchMin,settings.PitchMid,settings.PitchMax)
	ymin,ymid,ymax:=fixMax(settings.YawMin,settings.YawMid,settings.YawMax)
	tmin,tmid,tmax:=fixMax(settings.ThrottleMin,settings.ThrottleMid,settings.ThrottleMax)


	return &remoteControl{
		transmitter:            settings.Transmitter,
		commandPerSecond:       settings.CommandPerSecond,
		joyStick:               settings.JoyStick,
		rollChan:               settings.Roll,
		pitchChan:              settings.Pitch,
		yawChan:                settings.Yaw,
		throttleChan:           settings.Throttle,
		buttons:                settings.PushButtons,
		oled:                   settings.OLED,
		displayUpdatePerSecond: settings.DisplayUpdatePerSecond,
		lastCommandRead:        time.Now(),
		lastDisplayUpdate:      time.Now(),

		rollMin:                rmin,
		rollMid:                rmid,
		rollMax:                rmax,
		
		pitchMin:               pmin,
		pitchMid:               pmid,
		pitchMax:               pmax,

		yawMin:                 ymin,
		yawMid:                 ymid,
		yawMax:                 ymax,

		throttleMin:            tmin,
		throttleMid:            tmid,
		throttleMax:            tmax,

		rotationRange:          settings.RotationRange,
		
	}
}

func (r *remoteControl) Start(ctx context.Context) {
	running := true
	r.transmitter.On()
	r.Initisplay()
	// displayUpdate := utils.WithDataPerSecond(3)
	commandsUpdate := utils.WithDataPerSecond(r.commandPerSecond)
	fmt.Println("Commands per second: ", r.commandPerSecond)
	for running {
		select {
		default:
			r.updateButtons()
			if commandsUpdate.IsTime() {
				pressedButtons, pushButtons := r.readButtons()
				l, h := r.joyStick.Read(r.rollChan)
				lRoll, hRoll := normilise(l, h, r.rollMin, r.rollMid, r.rollMax, constants.JOY_STICK_INPUT_RANGE)
				l, h = r.joyStick.Read(r.pitchChan)
				lPitch, hPitch := normilise(l, h, r.pitchMin, r.pitchMid, r.pitchMax, constants.JOY_STICK_INPUT_RANGE)
				lYaw, hYaw := byte(0), byte(0)
				l, h = r.joyStick.Read(r.throttleChan)
				lThrottle, hThrottle := normilise(l, h, r.throttleMin, r.throttleMid, r.throttleMax, constants.JOY_STICK_INPUT_RANGE)

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

func fixMax(min, mid, max uint16) (uint16, uint16, uint16) {
	fmt.Print(min,mid, max, " -> ")
	r:=mid-min
	max=r+mid
	return min, mid, max
}

func normiliseThrottle(l, h byte, min, mid, max , inputRange uint16) uint16 {
	v:=uint16(l) + uint16(h)<<8
	if v<min {
		v=min
	}
	if v>max {
		v=max
	}
	r:=mid-min
	f:=float64(inputRange/2)/float64(r)
	n:=uint16(float64(v-min)*f)
	return n
}

func normilise(l, h byte, min, mid, max , inputRange uint16) (byte, byte) {
	n:=normiliseThrottle(l, h, min, mid, max , inputRange)
	return byte(n & 0b0000000011111111), byte(n>>8)
}