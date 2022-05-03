package main

import (
	pa "github.com/gordonklaus/portaudio"
	"github.com/hajimehoshi/ebiten/v2"
	man "github.com/softtacos/go-visualizer/manager"
	v "github.com/softtacos/go-visualizer/manager/visualizers"
	m "github.com/softtacos/go-visualizer/model"
	"log"
)

const (
	wWidth, wHeight = 1280, 960
)

func main() {
	var err error
	if err = pa.Initialize(); err != nil {
		log.Fatal("failed to initialize portaudio ", err)
		return
	}

	// vis manager currently just handles passing the draw calls to the visualizer
	// will eventually decide which visualizer should be set for the current vis
	// uses closures to store setup data for the visualizers so they can be created and destroyed easily at runtime
	karen := man.NewVisManager([]m.VisualizerConstructor{
		//v.NewForkingVisualizerConstructor,
		v.NewPolyVisualizerConstructor,
		//v.NewLazyPaddedVisualizer,
	})
	karen.Start()
	ebiten.SetWindowSize(wWidth, wHeight)
	ebiten.SetWindowResizable(true)
	//ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	ebiten.SetWindowTitle("go visualizer")
	if err = ebiten.RunGame(karen); err != nil {
		log.Fatal("error running game ", err)
	}
	karen.Stop()
}
