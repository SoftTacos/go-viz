package model

import eb "github.com/hajimehoshi/ebiten/v2"

type Visualizer interface {
	Draw(screen *eb.Image)
	//Start()
	//Stop()
}
