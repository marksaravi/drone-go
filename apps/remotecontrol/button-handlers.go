package remotecontrol

import (
	"context"
	"time"
)

func (rc *remoteControl) actOnShutdownButtonState(pressed bool, cancel context.CancelFunc) {
	if !pressed {
		rc.shutdownCountdown = time.Now()
	}
	if time.Since(rc.shutdownCountdown) > time.Duration(time.Second*3) {
		cancel()
	}
}

func (rc *remoteControl) actOnSuppressLostConnectionButtonState(pressed bool) {
	if !pressed {
		rc.suppressLostConnectionCountdown = time.Now()
	}
	if time.Since(rc.suppressLostConnectionCountdown) > time.Duration(time.Second/2) {
		rc.radio.SuppressLostConnection()
	}
}
