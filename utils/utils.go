package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var intervals map[string]time.Time = make(map[string]time.Time)

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

func WaitToAbortByESC(cancel context.CancelFunc, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Press ENTER to abort")
	go func(cancel context.CancelFunc) {
		defer log.Println("Aborting by user ENTER")
		defer wg.Done()
		// disable input buffering
		unixDisableInputBuffering()
		// do not display entered characters on the screen
		unixDisplayCharacterOnScreen(false)
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGQUIT)
		var b []byte = make([]byte, 1)
		for {
			select {
			case <-sigs:
				unixDisplayCharacterOnScreen(true)
				cancel()
				return
			default:
			}
			os.Stdin.Read(b)
			time.Sleep(50 * time.Millisecond)
			if b[0] == 27 {
				unixDisplayCharacterOnScreen(true)
				cancel()
				return
			}
		}
	}(cancel)
}

func ApplyLimits(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
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
