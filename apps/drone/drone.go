package drone

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/hardware/icm20789"
)

type imu interface {
	Setup()
	ReadIMUData() ([]byte, error)
}

type drone struct {
	imu imu

	imuSampleRate int

	lastIMUReadingTime time.Time
}

func NewDrone() *drone {
	return &drone{
		imu: icm20789.NewICM20789(),
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
		imuraw, imuok := d.readImuRawData()
		if imuok {
			log.Println(imuraw)
		}
	}
}

func (d *drone) readImuRawData() ([]byte, bool) {
	if time.Since(d.lastIMUReadingTime) < time.Second/time.Duration(d.imuSampleRate) {
		return nil, false
	}
	d.lastIMUReadingTime = time.Now()
	data, err := d.imu.ReadIMUData()
	return data, err == nil
}
