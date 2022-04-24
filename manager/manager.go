package manager

import (
	eb "github.com/hajimehoshi/ebiten/v2"
	m "github.com/softtacos/go-visualizer/model"
	"math"
	"math/cmplx"
)


type VisualizerConstructor func(chan []float32)m.Visualizer

func NewVisManager(audio chan []float32,constructors []VisualizerConstructor)*manager {
	return &manager{
		constructors:constructors,
		currentVisualizer: constructors[0](audio),
	}
}

type manager struct {
	audio chan []float64
	vIndex int
	currentVisualizer m.Visualizer
	constructors []VisualizerConstructor
}

func (v *manager)Update() (err error){
	// TODO: listen for input event to change visualizer

	//l := len(v.input)
	//frequencies :=make([]complex128,l)
	//DFFT64(<-v.input,frequencies,l,1)
	//v.currentMutex.Lock()
	//v.current = make([]float64,l)
	//v.currentMutex.Unlock()
	return
}

func (v *manager)Draw(screen *eb.Image) {
	v.currentVisualizer.Draw(screen)
}

func (v *manager)Layout(outsideWidth,outsideHeight int) (int,int){
	return outsideWidth,outsideHeight
}


func DFFT64(x []float64, y []complex128, n, s int) {
	if n == 1 {
		y[0] = complex(x[0], 0)
		return
	}
	DFFT64(x, y, n/2, 2*s)
	DFFT64(x[s:], y[n/2:], n/2, 2*s)
	for k := 0; k < n/2; k++ {
		tf := cmplx.Rect(1, -2*math.Pi*float64(k)/float64(n)) * y[k+n/2]
		y[k], y[k+n/2] = y[k]+tf, y[k]-tf
	}
}
