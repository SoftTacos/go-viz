package visualizers

import (
	"fmt"
	eb "github.com/hajimehoshi/ebiten/v2"
	m "github.com/softtacos/go-visualizer/model"
)

func NewDummyVisualizer(windowSize int,ampInput chan []float64) m.Visualizer {
	return &dummyVisualizer{
		ampInput: ampInput,
	}
}

type dummyVisualizer struct {
	ampInput   chan []float64
}

func (v *dummyVisualizer) Draw(screen *eb.Image) {
	if len(v.ampInput) > 0{
		<-v.ampInput
	}
}

func (v *dummyVisualizer) BeatCallback() {
	fmt.Println("BEAT")
}