package utils

import (
	"github.com/MarkSaravi/drone-go/models"
)

func SerializeFlightCommand(fc models.FlightCommands) []byte {
	return []byte{}
}

func DeserializeFlightCommand(data []byte) models.FlightCommands {
	return models.FlightCommands{}
}

const RADIO_PAYLOAD_SIZE int = 32

// func FlightStatusToBytes(flightCommands models.FlightCommands) []byte {
// 	isRemote := BoolToShiftedByte(flightCommands.IsRemoteControl, 0)
// 	isDrone := BoolToShiftedByte(flightCommands.IsDrone, 1)
// 	isMotorsEngaged := BoolToShiftedByte(flightCommands.IsMotorsEngaged, 2)
// 	return []byte{
// 		isRemote | isDrone | isMotorsEngaged,
// 		0,
// 	}
// }

// func flightDataToPayload(flightCommands models.FlightCommands) []byte {
// 	payload := append([]byte{}, UInt32ToBytes(flightCommands.Id)...)
// 	payload = append(payload, Float32ToBytes(flightCommands.Roll)...)
// 	payload = append(payload, Float32ToBytes(flightCommands.Pitch)...)
// 	payload = append(payload, Float32ToBytes(flightCommands.Yaw)...)
// 	payload = append(payload, Float32ToBytes(flightCommands.Throttle)...)
// 	payload = append(payload, Float32ToBytes(flightCommands.Altitude)...)
// 	payload = append(payload, ([]byte{0, 0, 0, 0, 0, 0})...)
// 	payload = append(payload, FlightStatusToBytes(flightCommands)...)
// 	return payload
// }

// func payloadToFlightData(payload []byte) models.FlightCommands {
// 	return models.FlightCommands{
// 		Id:       UInt32fromBytes(payload[0:4]),
// 		Roll:     Float32fromBytes(payload[4:8]),
// 		Pitch:    Float32fromBytes(payload[8:12]),
// 		Yaw:      Float32fromBytes(payload[12:16]),
// 		Throttle: Float32fromBytes(payload[16:20]),
// 		Altitude: Float32fromBytes(payload[20:24]),
// 	}
// }
