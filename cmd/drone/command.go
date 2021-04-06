package main

import (
	"fmt"
	"sync"

	commands "github.com/MarkSaravi/drone-go/constants"
	"github.com/MarkSaravi/drone-go/types"
)

func createCommandChannel(wg *sync.WaitGroup) chan types.Command {
	command := make(chan types.Command, 1)
	go func() {
		wg.Add(1)
		var input = "no"
		for input != "end" {
			fmt.Scanf("%s", &input)
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
