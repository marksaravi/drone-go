package command

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Command struct {
	Command string
}

func CreateCommandChannel(wg *sync.WaitGroup) chan Command {
	command := make(chan Command, 1)
	go func() {
		wg.Add(1)
		var b []byte = make([]byte, 1)
		for range time.Tick(time.Millisecond * 100) {
			os.Stdin.Read(b)
			if b[0] == '\n' {
				break
			}
		}
		command <- Command{
			Command: "COMMAND_END_PROGRAM",
		}
		close(command)
		fmt.Println("Command stopped.")
		wg.Done()
		return
	}()
	return command
}
