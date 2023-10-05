package main

import "github.com/marksaravi/drone-go/apps/plotter"

func main() {
	p := plotter.NewPlotter(plotter.PlotterSettings{
		UDPAddress: ":4010",
	})
	p.StartPlotter()
}
