package remotecontrol

import (
	"context"
	"log"
	"time"
)

const buttonPressTimeout = time.Duration(time.Second * 3)

func (rc *remoteControl) shutdownPressed(pressed bool, cancel context.CancelFunc) {
	if !pressed {
		rc.shutdownCountdown = time.Now()
	}
	if time.Since(rc.shutdownCountdown) > buttonPressTimeout {
		log.Println("Exiting the Remote Control program...")
		cancel()
		rc.shutdownCountdown = time.Now().AddDate(1, 0, 0)
	}
}

func (rc *remoteControl) suppressLostConnectionPressed(pressed bool) {
	if !pressed {
		rc.suppressLostConnectionCountdown = time.Now()
	}
	if time.Since(rc.suppressLostConnectionCountdown) > buttonPressTimeout {
		rc.suppressLostConnectionCountdown = time.Now().AddDate(1, 0, 0)
		rc.radio.SuppressLostConnection()
	}
}
