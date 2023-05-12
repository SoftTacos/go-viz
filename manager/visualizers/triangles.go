package visualizers

import (
	"image/color"

	eb "github.com/hajimehoshi/ebiten/v2"
	m "github.com/softtacos/go-visualizer/model"
	//"github.com/goccmack/godsp/peaks"
)

func NewRowsVisualizerConstructor(windowSize int, frequencyInput chan []float64) m.Visualizer {
	return newRowsVisualizer(windowSize, frequencyInput)
}

func newRowsVisualizer(windowSize int, frequencyInput chan []float64) *rowsVisualizer {
	const (
		triangleSideLength = 200
		numTriangles       = 10
		x, y               = 14, 7
	)
	var hi, med, low uint8 = 0xff, 0x90, 0x50
	v := &rowsVisualizer{
		frequencyInput: frequencyInput,
		shape:          generateTriangleShapeImage(triangleSideLength),
		colors: []color.RGBA{
			{low, hi, low, 1}, // green
			{low, med, hi, 1}, // blue
			{low, hi, hi, 1},  // cyan
			{hi, low, low, 1}, // red
			{hi, low, hi, 1},  // pink
		},
		shapeLocations: generateShapeLocations(triangleSideLength, x, y),
		// shapeGrid: generateGrid(10, 5),
	}
	return v
}

func generateTriangleShapeImage(px int) *eb.Image {
	emptyImage := eb.NewImage(px, px)
	emptyImage.Fill(color.Black)
	return eb.NewImage(px, px)
}

func generateShapeLocations(triangleSideLength, numX, numY int) (locations [][2]int) {
	locations = make([][2]int, numX)
	for i := 0; i < numX; i++ {
		for j := 0; j < numY; j++ {
			x:=
			locations[i*numY+j] = [2]int{x, y}
		}
	}
	return
}

type rowsVisualizer struct {
	frequencyInput chan []float64
	shape          *eb.Image
	colors         []color.RGBA
	shapeLocations [][2]int
	colorIndices   []int // index matches the index of the shape in shapeLocations, int is the index of the color in colors

}

func (v *rowsVisualizer) BeatCallback() {
	// TODO
}

func (v *rowsVisualizer) Draw(screen *eb.Image) {
	// var (
	// 	frequencies   = <-v.frequencyInput
	// 	w, h          = screen.Size()
	// 	width, height = float64(w), float64(h)
	// 	groups        = util.GroupFrequencies(5, frequencies)
	// )

	// screen.Fill(color.RGBA{
	// 	R: 0x00,
	// 	G: 0x00,
	// 	B: 0x00,
	// 	A: 0xff,
	// })
	// ops := eb.DrawImageOptions{}
	// var change float64 = groups[4]
	// if change < 0 {
	// 	change = 0
	// }
	// ops.ColorM.Scale(0, 0, 0, .9) //.1+change)
	// v.indexMutex.Lock()
	// ops.ColorM.Translate(v.colorFloats[v.colorIndex][0]/0xff, v.colorFloats[v.colorIndex][1]/0xff, v.colorFloats[v.colorIndex][2]/0xff, 0)
	// v.indexMutex.Unlock()
	// ops.GeoM.Translate(width/2, height/2)
	// ops.Filter = eb.FilterLinear
	// v.r += .015
	// scale := 2.0
	// rotStep, scaleStep := CalcPolyRotationScale(v.poly)
	// const depth = 25
	// SpiralNestPolygons(screen, v.poly, depth, scale, scaleStep, v.r, rotStep, ops)
	// SpiralNestPolygons(screen, v.poly, depth, scale, scaleStep, v.r+math.Pi, rotStep, ops)
}
