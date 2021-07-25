package utils

import (
	"periph.io/x/periph/conn/spi"
)

func ReadSPI(address byte, datalen int, conn spi.Conn) ([]byte, error) {
	w := make([]byte, datalen+1)
	r := make([]byte, datalen+1)
	w[0] = address
	err := conn.Tx(w, r)
	return r[1:], err
}

func WriteSPI(address byte, data []byte, conn spi.Conn) error {
	datalen := len(data)
	w := make([]byte, datalen+1)
	r := make([]byte, datalen+1)
	w[0] = address

	for i := 0; i < datalen; i++ {
		w[i+1] = data[i]
	}

	err := conn.Tx(w, r)
	return err
}
