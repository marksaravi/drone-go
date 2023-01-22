package drone

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/types"
)

type InertialMeasurementUnit interface {
	Read() (types.Rotations, error)
}

type drone struct {
	imu InertialMeasurementUnit

	imuSampleRate int

	lastIMUReadingTime time.Time
}

func NewDrone(imu InertialMeasurementUnit) *drone {
	return &drone{
		imu:                imu,
		imuSampleRate:      2,
		lastIMUReadingTime: time.Now(),
	}
}

func (d *drone) Start(ctx context.Context, wg *sync.WaitGroup) {
	log.Println("drone started")
	defer log.Println("drone stopped")
	d.controller(ctx, wg)
}

func (d *drone) controller(ctx context.Context, wg *sync.WaitGroup) {
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

func (d *drone) readRotations() (types.Rotations, bool) {
	if time.Since(d.lastIMUReadingTime) < time.Second/time.Duration(d.imuSampleRate) {
		return types.Rotations{}, false
	}
	d.lastIMUReadingTime = time.Now()
	rotations, err := d.imu.Read()
	if err != nil {
		return rotations, false
	}
	return rotations, true
}
