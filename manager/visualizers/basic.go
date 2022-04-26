package visualizers

import (
	"fmt"
	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/mjibson/go-dsp/fft"
	m "github.com/softtacos/go-visualizer/model"
	"image/color"
	"sync"
)

const thiccness float32 = 10

var (
	emptyDrawTriangleOptions = &eb.DrawTrianglesOptions{}
)

func NewBasicVisualizer(input chan []float64) m.Visualizer {
	return &basicVisualizer{
		input:        input,
		currentMutex: &sync.Mutex{},
	}
}

type basicVisualizer struct {
	input        chan []float64
	current      []float64
	currentMutex *sync.Mutex
}
var max float64

func (v *basicVisualizer) Draw(screen *eb.Image) {
	// get cf
	input := <-v.input
	cf := fft.FFTReal(input)
	cf = cf[0:len(cf)/2]
	l := len(cf)
	frequencies := make([]float64, l)
	for i := range cf {
		frequencies[i] = float64(real(cf[i]))
		if frequencies[i] > max{
			fmt.Println("MAX:",max)
			max = frequencies[i]
		}
	}

	//DrawFrequencies(screen,<-v.input)
	DrawFrequencies(screen, frequencies)
}

func DrawFrequencies(screen *eb.Image, frequencies []float64) {
	width, height := screen.Size()
	//frequencies := <-v.input
	//frequencies = frequencies[0 : len(frequencies)/3]
	l := len(frequencies)
	//fmt.Println(l)

	var (
		vertices  = make([]eb.Vertex, 0, (l+1)*2)
		midheight = float32(height) // /2
		xchunk    = float32(width) / float32(l)
	)
	vertices = append(vertices, eb.Vertex{
		DstX:   0,
		DstY:   midheight - thiccness,
		ColorR: 1,
		ColorG: 1,
		ColorB: 1,
		ColorA: 1,
	}, eb.Vertex{
		DstX:   0,
		DstY:   midheight + thiccness,
		ColorR: 1,
		ColorG: 1,
		ColorB: 1,
		ColorA: 1,
	})
	for i, f := range frequencies {
		x := (xchunk * float32(i+1))
		//x2:=xchunk*(i+1)

		fOffset := -float32(f) * 5.0//20000000 //float32(v.height)
		y1 := midheight + fOffset - 2*thiccness
		y2 := midheight + fOffset //+ thiccness

		vertices = append(vertices,
			eb.Vertex{
				DstX:   x,
				DstY:   y1,
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			},
			eb.Vertex{
				DstX:   x,
				DstY:   y2,
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			})
	}
	var indices = make([]uint16, 0, len(frequencies)*3)
	for i := uint16(0); i < uint16(len(frequencies)-1); i++ {
		indices = append(indices, i*2, (i*2)+1, (i*2)+3, (i * 2), (i*2)+3, (i*2)+2)
	}

	baseImg := eb.NewImage(width, height)
	baseImg.Fill(color.White)
	screen.DrawTriangles(vertices, indices, baseImg, emptyDrawTriangleOptions)
}
