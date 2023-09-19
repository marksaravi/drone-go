package remote_test

import (
	"testing"
	"time"

	"github.com/marksaravi/drone-go/apps/remote"
)

func TestReadCommandsBeforeTimeout(t *testing.T) {
	commandPerSecond := 20
	r := remote.NewRemote(remote.RemoteCongigs{CommandPerSecond: commandPerSecond})
	timeout := time.Second / time.Duration(commandPerSecond)
	time.Sleep(timeout / 2)
	_, ok := r.ReadCommands()
	if ok {
		t.Errorf("readcommand didn't wait for %v for reading", timeout)
	}
}

func TestReadCommandsAfterTimeout(t *testing.T) {
	commandPerSecond := 20
	r := remote.NewRemote(remote.RemoteCongigs{CommandPerSecond: commandPerSecond})
	timeout := time.Second / time.Duration(commandPerSecond)
	time.Sleep(timeout / 2)
	time.Sleep(timeout)
	_, ok := r.ReadCommands()
	if !ok {
		t.Errorf("readcommand didn't read after %v timeout", timeout)
	}

}
