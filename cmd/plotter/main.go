package main

import "github.com/marksaravi/drone-go/apps/plotter"

func main() {
	p := plotter.NewPlotter(plotter.PlotterSettings{
		UDPServerAddress:  ":4014",
		HTTPServerAddress: ":3000",
	})
	p.StartPlotter()
}
