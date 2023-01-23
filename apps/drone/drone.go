package drone

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
)

type inertialMeasurementUnit interface {
	Read() (imu.Rotations, error)
}

type plotter interface {
	SendRotation(imu.Rotations)
}

type drone struct {
	imu           inertialMeasurementUnit
	plotter       plotter
	imuSampleRate int

	lastIMUReadingTime time.Time
}

func NewDrone(imu inertialMeasurementUnit) *drone {
	return &drone{
		imu:                imu,
		imuSampleRate:      2,
		lastIMUReadingTime: time.Now(),
	}
}

func (d *drone) Fly(ctx context.Context, wg *sync.WaitGroup) {
	log.Println("drone started")
	defer log.Println("drone stopped")
	running := true
	for running {
		select {
		case <-ctx.Done():
			log.Println("STOP SIGNAL RECEIVED")
			running = false
		default:
		}
		rotations, imuok := d.readRotations()
		if imuok {
			log.Println(rotations)
		}
	}
}

func (d *drone) readRotations() (imu.Rotations, bool) {
	if time.Since(d.lastIMUReadingTime) < time.Second/time.Duration(d.imuSampleRate) {
		return imu.Rotations{}, false
	}
	d.lastIMUReadingTime = time.Now()
	rotations, err := d.imu.Read()
	if err != nil {
		return rotations, false
	}
	return rotations, true
}
