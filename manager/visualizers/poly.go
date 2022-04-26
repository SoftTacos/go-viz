package visualizers

import (
	eb "github.com/hajimehoshi/ebiten/v2"
	m "github.com/softtacos/go-visualizer/model"
	"github.com/softtacos/go-visualizer/util"
	"image/color"
	"math"
	"sync"
	//"github.com/goccmack/godsp/peaks"
)

// TODO: bump to golang v1.18 and make a generic visualizer setup func
type PolyOption func(v *polyVisualizer)

const (
	defaultBufferSize = 3
	defaultPolyRadius = 500
	defaultPolySides  = 6
)

func SetBufferFunc(bufferSize, windowSize int) PolyOption {
	return func(v *polyVisualizer) {
		v.buffer = util.NewFrequencyBuffer(bufferSize, windowSize)
	}
}

func SetPolyFunc(sides int, radius float64) PolyOption {
	return func(v *polyVisualizer) {
		v.poly = NewPolygon(sides, radius)
	}
}

func NewPolyVisualizerOptionsConstructor(windowSize int, frequencyInput chan []float64, options ...PolyOption) func() m.Visualizer {
	return func() m.Visualizer {
		v := newPolyVisualizer(windowSize, frequencyInput)
		for _, o := range options {
			o(v)
		}
		return v
	}
}

func NewPolyVisualizerConstructor(windowSize int, frequencyInput chan []float64) m.Visualizer {
	return newPolyVisualizer(windowSize, frequencyInput)
}

func newPolyVisualizer(windowSize int, frequencyInput chan []float64) *polyVisualizer {
	var hi, med, low float64 = 0xff, 0x90, 0x50
	v := &polyVisualizer{
		buffer:         util.NewFrequencyBuffer(defaultBufferSize, windowSize),
		frequencyInput: frequencyInput,
		poly:           NewPolygon(defaultPolySides, defaultPolyRadius),
		colorFloats: [][]float64{
			{low, hi, low, 1}, // green
			{low, med, hi, 1}, // blue
			{hi, low, low, 1}, // red
			{med, low, hi, 1}, // purple
		},
		indexMutex: &sync.Mutex{},
	}
	//v.testBeatDetector = util.NewBeatDetector(48000, 200, v.beatCallback)
	//go v.listen()
	return v
}

//func NewPolyVisualizer(bufferSize, samples int, ampInput chan []float64) m.Visualizer {
//	var hi, med, low float64 = 0xff, 0x90, 0x50
//	v := &polyVisualizer{
//		buffer:   util.NewFrequencyBuffer(bufferSize, samples),
//		ampInput: ampInput,
//		poly: NewPolygon(6, 500),
//		colorFloats: [][]float64{
//			{low, hi, low, 1},  // green
//			{low, med, hi, 1}, // blue
//			{hi, low, low, 1}, // red
//			{med, low, hi, 1},  // purple
//		},
//		indexMutex: &sync.Mutex{},
//	}
//	v.testBeatDetector = util.NewBeatDetector(48000, 200, v.beatCallback)
//	//go v.listen()
//	return v
//}

type polyVisualizer struct {
	frequencyInput chan []float64
	buffer         *util.Buffer
	poly           *Polygon
	colorFloats    [][]float64
	colors         []color.Color
	r                float64
	indexMutex       *sync.Mutex
	colorIndex       int
}

//
//func (v *polyVisualizer) listen() {
//	for {
//		amplitudes := <-v.ampInput
//
//		cf := fft.FFTReal(amplitudes)
//		cf = cf[0 : len(cf)/2]
//		frequencies := make([]float64, len(cf))
//		for i := range cf {
//			frequencies[i] = math.Abs(real(cf[i]))
//		}
//		v.buffer.Push(frequencies)
//		v.testBeatDetector.Push(frequencies...)
//
//	}
//}

func (v *polyVisualizer) BeatCallback() {
	v.indexMutex.Lock()
	defer v.indexMutex.Unlock()
	v.colorIndex = (v.colorIndex + 1) % len(v.colorFloats)
}

func (v *polyVisualizer) Draw(screen *eb.Image) {
	var (
		//frequencies   = v.buffer.GetAverage()
		frequencies   = <-v.frequencyInput
		w, h          = screen.Size()
		width, height = float64(w), float64(h)
		groups        = util.GroupFrequencies(5, frequencies)
		//dGroups = util.GroupFrequencies(5, v.buffer.GetDerivative())
	)
	//test:=uint8(lazyClamp(dGroups[3]*1000000.0,255))
	//fmt.Println(dGroups[3],test)
	screen.Fill(color.RGBA{
		R: 0x00,
		G: 0x00,
		B: 0x00,
		A: 0xff,
	})
	ops := eb.DrawImageOptions{}
	var change float64 = groups[4]
	if change < 0 {
		change = 0
	}
	ops.ColorM.Scale(0, 0, 0, .9) //.1+change)
	v.indexMutex.Lock()
	ops.ColorM.Translate(v.colorFloats[v.colorIndex][0]/0xff, v.colorFloats[v.colorIndex][1]/0xff, v.colorFloats[v.colorIndex][2]/0xff, 0)
	v.indexMutex.Unlock()
	ops.GeoM.Translate(width/2, height/2)
	ops.Filter = eb.FilterLinear
	v.r += .015
	scale := 2.0
	rotStep, scaleStep := CalcPolyRotationScale(v.poly)
	SpiralNestPolygons(screen, v.poly, 22, scale, scaleStep, v.r, rotStep, ops)
	SpiralNestPolygons(screen, v.poly, 22, scale, scaleStep, v.r+math.Pi, rotStep, ops)
}

func SpiralNestPolygons(screen *eb.Image, poly *Polygon, depth int, scale, scaleStep, rotation, rotationStep float64, ops eb.DrawImageOptions) {
	if depth < 1 {
		return
	}

	cirleW, circleH := poly.GetImg().Size()
	local := eb.DrawImageOptions{
		ColorM: ops.ColorM,
	}
	local.GeoM.Translate(-float64(cirleW)/2, -float64(circleH)/2)
	local.GeoM.Rotate(rotation)
	local.GeoM.Scale(scale, scale)

	local.GeoM.Concat(ops.GeoM)
	screen.DrawImage(poly.GetImg(), &local)
	rotation += rotationStep
	scale *= scaleStep
	SpiralNestPolygons(screen, poly, depth-1, scale, scaleStep, rotation, rotationStep, ops)
}

func CalcPolyRotationScale(poly *Polygon) (rotation, scale float64) {
	scale = math.Sqrt(math.Pow((.5+.5*math.Cos(math.Pi-poly.GetTheta())), 2) + math.Pow(.5*math.Sin(math.Pi-poly.GetTheta()), 2))
	rotation = (math.Pi - poly.GetTheta()) / 2.0
	return
}

func NewPolygon(n int, radius float64) *Polygon {
	if n < 3 {
		n = 3
	}
	p := &Polygon{
		n:      float64(n),
		radius: radius,
		img:    CreatePolyImgHollow(radius, n, .96),
	}
	p.calculateTheta()
	p.calculateL()
	return p
}

type Polygon struct {
	radius float64
	n      float64
	theta  float64 // radians, interior angle
	l      float64 // length of one side

	img *eb.Image
}

func (p *Polygon) GetImg() *eb.Image {
	return p.img
}

func (p *Polygon) GetL() float64 {
	return p.l
}

func (p *Polygon) GetTheta() float64 {
	return p.theta
}

func (p *Polygon) SetRadius(radius float64) {
	p.radius = radius
	p.calculateTheta()
	p.calculateL()
}

func (p *Polygon) calculateTheta() {
	p.theta = (p.n - 2) * math.Pi / p.n
}

func (p *Polygon) calculateL() {
	p.l = 2 * p.radius * math.Cos(p.theta/2)
}

// fill is a ratio fill must be between 0 and 1
func CreatePolyImgHollow(radius float64, numPoly int, fill float64) (polyImg *eb.Image) {
	if fill <= 0 || fill > 1 {
		return CreatePolyImgFill(radius, numPoly)
	}

	px := 2 * (int(radius) + 1)
	emptyImage := eb.NewImage(px, px)
	emptyImage.Fill(color.White)
	polyImg = eb.NewImage(px, px)
	op := eb.DrawTrianglesOptions{}
	vertices := []eb.Vertex{}
	indices := []uint16{}
	for i := 0; i < numPoly; i++ {
		rate := float64(i) / float64(numPoly)
		x1 := float32(radius * fill * math.Cos(float64(2)*math.Pi*rate))
		x2 := float32(radius * math.Cos(float64(2)*math.Pi*rate))
		y1 := float32(radius * fill * math.Sin(float64(2)*math.Pi*rate))
		y2 := float32(radius * math.Sin(float64(2)*math.Pi*rate))
		vertices = append(vertices, eb.Vertex{
			DstX:   x1 + float32(radius),
			DstY:   y1 + float32(radius),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		}, eb.Vertex{
			DstX:   x2 + float32(radius),
			DstY:   y2 + float32(radius),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		})
	}
	var l = uint16(len(vertices))
	for i := uint16(0); i < l; i += 2 {
		indices = append(indices, i, (i)+1, (i+3)%l, (i), (i+3)%l, (i+2)%l)
	}

	//for i := uint16(1); int(i) < len(vertices)-1; i++ {
	//	indices = append(indices, 0, i, i+1)
	//}
	indices = append(indices, l-2, 1, 0)
	polyImg.DrawTriangles(vertices, indices, emptyImage, &op)
	return
}

func CreatePolyImgFill(radius float64, numPoly int) (polyImg *eb.Image) {
	px := 2 * (int(radius) + 1)
	emptyImage := eb.NewImage(px, px)
	emptyImage.Fill(color.White)
	polyImg = eb.NewImage(px, px)
	op := eb.DrawTrianglesOptions{}
	vertices := []eb.Vertex{}
	indices := []uint16{}
	vertices = append(vertices, eb.Vertex{
		DstX:   float32(radius),
		DstY:   float32(radius),
		ColorR: 1,
		ColorG: 1,
		ColorB: 1,
		ColorA: 1,
	})
	for i := 0; i < numPoly; i++ {
		rate := float64(i) / float64(numPoly)
		x := float32(radius * math.Cos(float64(2)*math.Pi*rate))
		y := float32(radius * math.Sin(float64(2)*math.Pi*rate))
		vertices = append(vertices, eb.Vertex{
			DstX:   x + float32(radius),
			DstY:   y + float32(radius),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		})
	}
	for i := uint16(1); int(i) < len(vertices)-1; i++ {
		indices = append(indices, 0, i, i+1)
	}
	indices = append(indices, 0, uint16(len(vertices)-1), 1)
	polyImg.DrawTriangles(vertices, indices, emptyImage, &op)
	return
}
