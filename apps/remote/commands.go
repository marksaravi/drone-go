package remote

import (
	"fmt"
	"time"
)

var sr uint64 = 0
var sp uint64 = 0
var sy uint64 = 0
var counter uint64 = 0

func (r *remoteControl) ReadCommands() (commands, bool) {
	if time.Since(r.lastCommand) < time.Second/time.Duration(r.commandPerSecond) {
		return commands{}, false
	}
	r.lastCommand = time.Now()

	roll := r.roll.Read()
	pitch := r.pitch.Read()
	yaw := r.yaw.Read()
	throttle := r.throttle.Read()

	sr += uint64(roll)
	sp += uint64(pitch)
	sy += uint64(yaw)
	counter++
	fmt.Printf("%6d %6d %6d\n", sr/counter, sp/counter, sy/counter)
	cmd := commands{
		roll:     roll,
		pitch:    pitch,
		yaw:      yaw,
		throttle: throttle,
	}
	return cmd, true
}
