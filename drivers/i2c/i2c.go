package i2c

import (
	"os"
	"syscall"
	"unsafe"
)

// Connection is a connection to i2c device in /dev/i2c-x
type Connection struct {
	f *os.File
}

// Config contains i2c device name and address
type Config struct {
	DevName string
	Address uint8
}

// Open opens an i2c connection by /dev/devicename
func Open(devname string) (c *Connection, err error) {
	f, err := os.OpenFile(devname, syscall.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return &Connection{f: f}, nil
}

// Close closes the connection to i2c
func (c *Connection) Close() (err error) {
	return c.f.Close()
}

// ReadByte reads a byte from address and offset
func (c *Connection) ReadByte(slave uint8, offset uint8) (b uint8, err error) {
	buf := []uint8{0}
	msg := []message{
		{
			addr:  uint16(slave),
			flags: 0,
			len:   1,
			buf:   uintptr(unsafe.Pointer(&offset)),
		},
		{
			addr:  uint16(slave),
			flags: uint16(M_RD),
			len:   uint16(len(buf)),
			buf:   uintptr(unsafe.Pointer(&buf[0])),
		},
	}
	err = transfer(c.f, &msg[0], len(msg))
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

// ReadWord reads a word from address and offset
func (c *Connection) ReadWord(slave uint8, offset uint8) (w uint16, err error) {
	buf := []uint8{0, 0}
	msg := []message{
		{
			addr:  uint16(slave),
			flags: 0,
			len:   1,
			buf:   uintptr(unsafe.Pointer(&offset)),
		},
		{
			addr:  uint16(slave),
			flags: uint16(M_RD),
			len:   uint16(len(buf)),
			buf:   uintptr(unsafe.Pointer(&buf[0])),
		},
	}
	err = transfer(c.f, &msg[0], len(msg))
	if err != nil {
		return 0, err
	}
	w = (uint16(buf[0]) << 8) | uint16(buf[1])
	return
}

// WriteBytes writes a byte from address and offset
func (c *Connection) WriteBytes(slave uint8, offset uint8, bytes ...uint8) (err error) {
	buf := []uint8{offset}
	for _, b := range bytes {
		buf = append(buf, b)
	}

	msg := []message{
		{
			addr:  uint16(slave),
			flags: 0,
			len:   uint16(len(buf)),
			buf:   uintptr(unsafe.Pointer(&buf[0])),
		},
	}
	return transfer(c.f, &msg[0], len(msg))
}

// WriteByte writes a byte from address and offset
func (c *Connection) WriteByte(slave uint8, offset uint8, b uint8) (err error) {
	buf := []uint8{offset, b}
	msg := []message{
		{
			addr:  uint16(slave),
			flags: 0,
			len:   uint16(len(buf)),
			buf:   uintptr(unsafe.Pointer(&buf[0])),
		},
	}
	return transfer(c.f, &msg[0], len(msg))
}

const (
	RDWR                = 0x0707
	RDRW_IOCTL_MAX_MSGS = 42
	M_RD                = 0x0001
)

type message struct {
	addr    uint16
	flags   uint16
	len     uint16
	padding uint16
	buf     uintptr
}

type RDWRRIoCtlData struct {
	msgs  uintptr
	nmsgs uint32
}

func transfer(f *os.File, msgs *message, n int) (err error) {
	data := RDWRRIoCtlData{
		msgs:  uintptr(unsafe.Pointer(msgs)),
		nmsgs: uint32(n),
	}
	err = nil
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(f.Fd()),
		uintptr(RDWR),
		uintptr(unsafe.Pointer(&data)),
	)
	if errno != 0 {
		err = errno
	}
	return
}
