package remote

import "time"

func (r *remote) ReadCommands() (bool, commands) {
	if time.Since(r.lastCommand) < time.Second/time.Duration(r.commandPerSecond) {
		return false, commands{}
	}
	r.lastCommand = time.Now()
	return true, commands{}
}
