package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	commands "github.com/MarkSaravi/drone-go/constants"
	"github.com/MarkSaravi/drone-go/types"
)

func createCommandChannel(wg *sync.WaitGroup) chan types.Command {
	command := make(chan types.Command, 1)
	go func() {
		wg.Add(1)
		var b []byte = make([]byte, 1)
		for range time.Tick(time.Millisecond * 100) {
			os.Stdin.Read(b)
			if b[0] == '\n' {
				break
			}
		}
		command <- types.Command{
			Command: commands.COMMAND_END_PROGRAM,
		}
		close(command)
		fmt.Println("Command stopped.")
		wg.Done()
		return
	}()
	return command
}
