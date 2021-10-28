// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
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
	fmt.Println(dev.Bounds())
	if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
		log.Fatal(err)
	}

	for y := 0; y < dev.Bounds().Dy(); y++ {
		for x := 0; x < dev.Bounds().Dx(); x++ {
			setPixel(x, y, dev.Bounds(), img)
		}
	}

	if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
		log.Fatal(err)
	}
}

func setPixel(x, y int, bounds image.Rectangle, img *image1bit.VerticalLSB) {
	absIndex := (y/8)*bounds.Dx() + x
	maskIndex := y % 8
	maskValue := byte(1) << byte(maskIndex)

	img.Pix[absIndex] = img.Pix[absIndex] | maskValue
}
