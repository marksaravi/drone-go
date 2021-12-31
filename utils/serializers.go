package utils

import (
	"github.com/marksaravi/drone-go/constants"
	"github.com/marksaravi/drone-go/models"
)

func shift2Bytes(x uint16, shift uint8) byte {
	return byte(x & 0b0000001100000000 >> (8 - shift))
}

func SerializeFlightCommand(flightCommands models.FlightCommands) models.Payload {

	var payloadType byte = flightCommands.PayloadType
	var reserverd byte = 0
	bottons := BoolArrayToByte([8]bool{
		flightCommands.ButtonFrontLeft,
		flightCommands.ButtonFrontRight,
		flightCommands.ButtonTopLeft,
		flightCommands.ButtonTopRight,
		flightCommands.ButtonBottomLeft,
		flightCommands.ButtonBottomRight,
		false,
		false,
	})

	var roll8bit byte = byte(flightCommands.Roll)
	var rollUpper2Bits byte = shift2Bytes(flightCommands.Roll, 0)
	var pitch8bit byte = byte(flightCommands.Pitch)
	var pitchUpper2Bits byte = shift2Bytes(flightCommands.Pitch, 2)
	var yaw8bit byte = byte(flightCommands.Yaw)
	var yawUpper2Bits byte = shift2Bytes(flightCommands.Yaw, 4)
	var throttle8bit byte = byte(flightCommands.Throttle)
	var throttleUpper2Bits byte = shift2Bytes(flightCommands.Throttle, 6)

	return [constants.RADIO_PAYLOAD_SIZE]byte{
		payloadType,
		reserverd,
		bottons,
		roll8bit,
		pitch8bit,
		yaw8bit,
		throttle8bit,
		rollUpper2Bits | pitchUpper2Bits | yawUpper2Bits | throttleUpper2Bits,
	}
}

func to10bits(lower8Bites byte, upper2Bits byte, pos uint8) uint16 {
	return uint16((upper2Bits>>pos)&0b00000011)<<8 | uint16(lower8Bites)
}

func DeserializeFlightCommand(payload models.Payload) models.FlightCommands {
	buttons := BoolArrayFromByte(payload[2])

	return models.FlightCommands{
		PayloadType:       payload[0],
		Roll:              to10bits(payload[3], payload[7], 0),
		Pitch:             to10bits(payload[4], payload[7], 2),
		Yaw:               to10bits(payload[5], payload[7], 4),
		Throttle:          to10bits(payload[6], payload[7], 6),
		ButtonFrontLeft:   buttons[0],
		ButtonFrontRight:  buttons[1],
		ButtonTopLeft:     buttons[2],
		ButtonTopRight:    buttons[3],
		ButtonBottomLeft:  buttons[4],
		ButtonBottomRight: buttons[5],
	}
}
