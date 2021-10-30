// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"image"
	"log"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/devices/ssd1306/image1bit"
	"periph.io/x/periph/host"
)

func main() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()
	ops := ssd1306.DefaultOpts
	ops.H = 64
	dev, err := ssd1306.NewI2C(b, &ops)
	if err != nil {
		log.Fatalf("failed to initialize ssd1306: %v", err)
	}

	img := image1bit.NewVerticalLSB(dev.Bounds())
	if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
		log.Fatal(err)
	}

	// for i := 0; i < 36; i++ {
	// 	writeChar(dev, img, i+65, i, i/12)
	// }

	writeString(dev, img, "Hello Mark!", 0, 0)
	writeString(dev, img, "Disconnected", 0, 1)
	writeString(dev, img, "Power: 54.3%", 0, 2)

	if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
		log.Fatal(err)
	}
}

func writeString(dev *ssd1306.Dev, img *image1bit.VerticalLSB, msg string, x, y int) {
	charCodes := []byte(msg)
	for i := 0; i < len(charCodes); i++ {
		writeChar(dev, img, int(charCodes[i]), x+i, y)
	}
}

func writeChar(dev *ssd1306.Dev, img *image1bit.VerticalLSB, charCode, x, y int) {
	const CHAR_W = 9
	const CHAR_H = 13
	var xOffset = CHAR_W*x + 4
	var yOffset = (CHAR_H + 10) * y
	char := monoFont[charCode]
	for row := 0; row < CHAR_H; row++ {
		for col := 0; col < CHAR_W; col++ {
			if char[row][col] > 0 {
				setPixel(row+yOffset, col+xOffset, dev.Bounds(), img)
			}
		}
	}

}

func setPixel(row, col int, bounds image.Rectangle, img *image1bit.VerticalLSB) {
	absIndex := (row/8)*bounds.Dx() + col
	maskIndex := row % 8
	maskValue := byte(1) << byte(maskIndex)

	img.Pix[absIndex] = img.Pix[absIndex] | maskValue
}
