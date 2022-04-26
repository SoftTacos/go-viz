package model

import eb "github.com/hajimehoshi/ebiten/v2"

type VisualizerConstructor func(windowSize int, frequencyInput chan []float64) Visualizer

type Visualizer interface {
	Draw(screen *eb.Image)
	BeatCallback()
	//Start()
	//Stop()
}
