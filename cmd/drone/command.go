package main

import (
	"fmt"

	commands "github.com/MarkSaravi/drone-go/constants"
	"github.com/MarkSaravi/drone-go/types"
)

func createCommandChannel() chan types.Command {
	command := make(chan types.Command, 1)
	go func() {
		var input = "no"
		for input != "end" {
			fmt.Printf("Waiting for command: ")
			fmt.Scanf("%s", &input)
		}
		command <- types.Command{
			Command: commands.COMMAND_END_PROGRAM,
		}
	}()
	return command
}
