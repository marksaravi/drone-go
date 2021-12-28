package utils

import (
	"github.com/marksaravi/drone-go/constants"
	"github.com/marksaravi/drone-go/models"
)

func SerializeFlightCommand(flightCommands models.FlightCommands) models.Payload {
	var payloadType byte = flightCommands.PayloadType
	var id byte = byte(flightCommands.Id)
	var roll byte = flightCommands.Roll
	var pitch byte = flightCommands.Pitch
	var yaw byte = flightCommands.Yaw
	var throttle byte = flightCommands.Throttle
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
	return [constants.RADIO_PAYLOAD_SIZE]byte{
		payloadType,
		id,
		roll,
		pitch,
		yaw,
		throttle,
		bottons,
		0,
	}
}

func DeserializeFlightCommand(payload models.Payload) models.FlightCommands {
	buttons := BoolArrayFromByte(payload[6])
	return models.FlightCommands{
		PayloadType:       payload[0],
		Id:                payload[1],
		Roll:              payload[2],
		Pitch:             payload[3],
		Yaw:               payload[4],
		Throttle:          payload[5],
		ButtonFrontLeft:   buttons[0],
		ButtonFrontRight:  buttons[1],
		ButtonTopLeft:     buttons[2],
		ButtonTopRight:    buttons[3],
		ButtonBottomLeft:  buttons[4],
		ButtonBottomRight: buttons[5],
	}
}
