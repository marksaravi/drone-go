package remotecontrol

import (
	"context"
	"log"
	"time"
)

const shutdownTimeout = time.Duration(time.Second * 3)
const alertResetTimeout = time.Duration(time.Second / 2)

func (rc *remoteControl) shutdownPressed(pressed bool, cancel context.CancelFunc) {
	if !pressed {
		rc.shutdownCountdown = time.Now()
	}
	if time.Since(rc.shutdownCountdown) > shutdownTimeout {
		log.Println("Exiting the Remote Control program...")
		cancel()
		rc.shutdownCountdown = time.Now().AddDate(1, 0, 0)
	}
}

func (rc *remoteControl) suppressLostConnectionPressed(pressed bool) {
	if !pressed {
		rc.suppressLostConnectionCountdown = time.Now()
	}
	if time.Since(rc.suppressLostConnectionCountdown) > alertResetTimeout {
		rc.suppressLostConnectionCountdown = time.Now().AddDate(1, 0, 0)
		log.Println("suppressing...")
	}
}
