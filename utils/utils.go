package utils

import (
	"fmt"
	"log"
	"math/big"
	"os/exec"
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

func UInt64ToBytes2(n int64) []byte {
	big := new(big.Int)
	big.SetInt64(n)
	return big.Bytes()
}
