package remote

import "time"

func (r *remoteControl) ReadCommands() (commands, bool) {
	if time.Since(r.lastCommand) < time.Second/time.Duration(r.commandPerSecond) {
		return commands{}, false
	}
	r.lastCommand = time.Now()
	return commands{}, true
}
