package constants

import "github.com/marksaravi/drone-go/hardware/mcp3008"

const RADIO_PAYLOAD_SIZE uint8 = 8
const JOYSTICK_RESOLUTION = mcp3008.DIGITAL_MAX_VALUE

const (
	COMMAND_DUMMY uint8 = iota
	COMMAND_CONTROL
	COMMAND_SHUTDOWN_DRONE
)
