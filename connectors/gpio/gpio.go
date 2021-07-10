package gpio

import (
	"errors"
	"os"
	"reflect"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
)

//Pins in J8
const (
	J8Pin27 = iota
	J8Pin28
	J8Pin03
	J8Pin05
	J8Pin07
	J8Pin29
	J8Pin31
	J8Pin26
	J8Pin24
	J8Pin21
	J8Pin19
	J8Pin23
	J8Pin32
	J8Pin33
	J8Pin08
	J8Pin10
	J8Pin36
	J8Pin11
	J8Pin12
	J8Pin35
	J8Pin38
	J8Pin40
	J8Pin15
	J8Pin16
	J8Pin18
	J8Pin22
	J8Pin37
	J8Pin13
	MaxGPIOPin
)

// GPIO aliases to J8 pins
const (
	GPIO01 = J8Pin28
	GPIO02 = J8Pin03
	GPIO03 = J8Pin05
	GPIO04 = J8Pin07
	GPIO05 = J8Pin29
	GPIO06 = J8Pin31
	GPIO07 = J8Pin26
	GPIO08 = J8Pin24
	GPIO09 = J8Pin21
	GPIO10 = J8Pin19
	GPIO11 = J8Pin23
	GPIO12 = J8Pin32
	GPIO13 = J8Pin33
	GPIO14 = J8Pin08
	GPIO15 = J8Pin10
	GPIO16 = J8Pin36
	GPIO17 = J8Pin11
	GPIO18 = J8Pin12
	GPIO19 = J8Pin35
	GPIO20 = J8Pin38
	GPIO21 = J8Pin40
	GPIO22 = J8Pin15
	GPIO23 = J8Pin16
	GPIO24 = J8Pin18
	GPIO25 = J8Pin22
	GPIO26 = J8Pin37
	GPIO27 = J8Pin13
)

const (
	memLen = 4096

	modeMask uint32 = 7 // pin mode is 3 bits wide
	pullMask uint32 = 3 // pull mode is 2 bits wide
	// BCM2835 pullReg is the same for all pins.
	pullReg2835 = 37
)

var (
	memoryLock sync.Mutex
	memory8bit []uint8
	memory     []uint32
	usedPins   [64]bool
)

//Level is High and Low
type Level bool
type Mode int
type Pull int

const (
	Input Mode = iota
	Output
)

const (
	Low  Level = false
	High Level = true
)

// Open opens gpio for read/write
func Open() (err error) {
	devmemFile, err := os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return
	}
	defer devmemFile.Close()
	memoryLock.Lock()
	defer memoryLock.Unlock()
	memory8bit, err = unix.Mmap(int(devmemFile.Fd()), 0, memLen, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		return
	}
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&memory8bit))
	//mapping to 32bit memory
	header.Len /= 4
	header.Cap /= 4
	memory = *(*[]uint32)(unsafe.Pointer(&header))

	return
}

// Close close the gpio
func Close() error {
	memoryLock.Lock()
	defer memoryLock.Unlock()
	memory = make([]uint32, 0)
	return unix.Munmap(memory8bit)
}

//Pin is a J8 pin
type Pin struct {
	pin         int
	fsel        int
	levelReg    int
	clearReg    int
	setReg      int
	pullReg2711 int
	bank        int
	mask        uint32
	shadow      Level
}

func (pin *Pin) setMode(mode Mode) {
	// shift for pin mode field within fsel Register.
	modeShift := uint(pin.pin%10) * 3

	memoryLock.Lock()
	defer memoryLock.Unlock()

	memory[pin.fsel] = memory[pin.fsel]&^(modeMask<<modeShift) | uint32(mode)<<modeShift
}

func (pin *Pin) write(level Level) {
	memoryLock.Lock()
	defer memoryLock.Unlock()
	if level == Low {
		memory[pin.clearReg] = pin.mask
	} else {
		memory[pin.setReg] = pin.mask
	}
	pin.shadow = level
}

//SetAsInput sets pin as input
func (pin *Pin) SetAsInput() {
	pin.setMode(Input)
}

//SetAsOutput sets pin as output
func (pin *Pin) SetAsOutput() {
	pin.setMode(Output)
}

//SetHigh sets pin to High
func (pin *Pin) SetHigh() {
	pin.write(High)
}

//SetLow sets pin to Low
func (pin *Pin) SetLow() {
	pin.write(Low)
}

//NewPin creates new pin
func NewPin(pin int) (*Pin, error) {
	if usedPins[pin] {
		return nil, errors.New("Pin is already taken")
	}
	usedPins[pin] = true
	if pin < 0 || pin >= MaxGPIOPin {
		return nil, errors.New("Invalid pin number")
	}
	fsel := pin / 10
	bank := pin / 32
	mask := uint32(1 << uint(pin&0x1f))
	levelReg := 13 + bank
	clearReg := 10 + bank
	setReg := 7 + bank
	pullReg := 57 + pin/16
	shadow := Low
	if memory[levelReg]&mask != 0 {
		shadow = High
	}
	return &Pin{
		pin:         pin,
		fsel:        fsel,
		bank:        bank,
		mask:        mask,
		levelReg:    levelReg,
		clearReg:    clearReg,
		pullReg2711: pullReg,
		setReg:      setReg,
		shadow:      shadow,
	}, nil
}
