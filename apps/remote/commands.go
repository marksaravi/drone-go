package remote

import (
	"time"
)

var roll byte = 0
var pitch byte = 10
var yaw byte = 20
var throttle byte = 30

func (r *remoteControl) ReadCommands() (commands, bool) {
	if time.Since(r.lastCommand) < time.Second/time.Duration(r.commandPerSecond) {
		return commands{}, false
	}
	r.lastCommand = time.Now()
	roll++
	pitch++
	yaw++
	throttle++
	return commands{roll: roll, pitch: pitch, yaw: yaw, throttle: throttle}, true
}
