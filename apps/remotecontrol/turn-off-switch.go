package remotecontrol

import (
	"context"
	"log"
	"time"
)

func (rc *remoteControl) processShutdownPressed(pressed bool, cancel context.CancelFunc) {
	if !pressed {
		rc.shutdownCountdown = time.Now()
	}
	if time.Since(rc.shutdownCountdown) > time.Duration(3*time.Second) {
		log.Println("Exiting the Remote Control program...")
		cancel()
		rc.shutdownCountdown = time.Now().AddDate(1, 0, 0)
	}
}
