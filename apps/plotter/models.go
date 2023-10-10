package plotter

// Plotter rotations is representation of Roll, Pitch and Yaw value in range of 0..35999 which will be converted to float by range 0..359.99. 
// Every angle had 2 digit percision and can be converted to raange -179.99..179.99.
// This format is inteded to compact an angle while providing enough percisions
type rotations struct {
	Roll uint16
	Pitch uint16
	Yaw uint16
}

type PlotterRotationsDataPacket string {
	DurationFromDroneStart uint32
	Rotations rotations
	AccelerometerRotattions rotations
	GyroscopeRotattions rotations
	Throttle byte
}
