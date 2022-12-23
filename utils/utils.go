package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var intervals map[string]time.Time = make(map[string]time.Time)

func WaitToAbortByESC(cancel context.CancelFunc) {
	log.Println("Press ESC to abort")
	go func(cancel context.CancelFunc) {
		// disable input buffering
		unixDisableInputBuffering()
		// do not display entered characters on the screen
		unixDisplayCharacterOnScreen(false)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM)
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			if b[0] == 27 {
				unixDisplayCharacterOnScreen(true)
				cancel()
				break
			}
		}
	}(cancel)
}

func unixDisableInputBuffering() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
}

func unixDisplayCharacterOnScreen(enable bool) {
	if enable {
		exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	} else {
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	}
}

func PrintIntervally(msg string, id string, interval time.Duration, useLog bool) {
	log.SetFlags(log.Lmicroseconds)
	now := time.Now()
	last, ok := intervals[id]
	if !ok {
		last = now
		intervals[id] = now
	}
	if time.Since(last) >= interval {
		intervals[id] = now
		if useLog {
			log.Print(msg)
		} else {
			fmt.Print(msg)
		}
	}
}
