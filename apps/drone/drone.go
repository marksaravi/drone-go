package drone

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
)

type imuDevice interface {
	Setup()
	ReadInertialDevice() (imu.Rotations, bool)
}

type drone struct {
	imu imuDevice

	imuSampleRate int

	lastIMUReadingTime time.Time
}

func NewDrone(imu imuDevice) *drone {
	return &drone{
		imu:                imu,
		imuSampleRate:      2,
		lastIMUReadingTime: time.Now(),
	}
}

func (d *drone) Start(ctx context.Context, wg *sync.WaitGroup) {
	log.Println("drone started")
	defer log.Println("drone stopped")
	d.imu.Setup()
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

func (d *drone) readRotations() (imu.Rotations, bool) {
	if time.Since(d.lastIMUReadingTime) < time.Second/time.Duration(d.imuSampleRate) {
		return imu.Rotations{}, false
	}
	d.lastIMUReadingTime = time.Now()
	return d.imu.ReadInertialDevice()
}
