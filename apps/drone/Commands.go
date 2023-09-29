package drone

import (
	"time"
)

func (d *droneApp) ReceiveCommand() ([]byte, bool) {
	if time.Since(d.lastCommand) < time.Second/time.Duration(2*d.commandsPerSecond) {
		return nil, false
	}
	d.lastCommand = time.Now()
	return d.receiver.Receive()
}
