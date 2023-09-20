package remote

import (
	"time"
)

func (r *remoteControl) ReadCommands() (commands, bool) {
	if time.Since(r.lastCommand) < time.Second/time.Duration(r.commandPerSecond) {
		return commands{}, false
	}
	r.lastCommand = time.Now()

	roll := r.roll.Read()
	pitch := r.pitch.Read()
	yaw := r.yaw.Read()
	throttle := r.throttle.Read()

	cmd := commands{
		roll:     roll,
		pitch:    pitch,
		yaw:      yaw,
		throttle: throttle,
	}
	return cmd, true
}
