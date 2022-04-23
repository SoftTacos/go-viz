package visualizers

import (
	"fmt"
	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/mjibson/go-dsp/fft"
	m "github.com/softtacos/go-visualizer/model"
	"github.com/softtacos/go-visualizer/util"
	"image/color"
	"math"
	//"github.com/goccmack/godsp/peaks"
)

func NewCircleVisualizer(bufferSize int, ampInput chan []float64) m.Visualizer {
	var hi, low float64 = 0xff, 0x40 //,0x90
	v := &circleVisualizer{
		buffer:   util.NewFrequencyBuffer(bufferSize),
		ampInput: ampInput,
		//poly:     CreatePolyImgHollow(1000,6),
		poly: NewPolygon(6, 100),
		colorFloats: [][]float64{
			{low, hi, low, 1},
			{low, 0x90, hi, 1},
		},
		colors: []color.Color{
			color.RGBA{
				R: 255,
				A: 255,
			},
			color.RGBA{
				G: 255,
				A: 255,
			},
			color.RGBA{
				B: 255,
				A: 255,
			},
			color.RGBA{
				R: 255,
				G: 255,
				A: 255,
			},
			color.RGBA{
				R: 255,
				B: 255,
				A: 255,
			},
			color.RGBA{
				B: 255,
				G: 255,
				A: 255,
			},
		},
	}
	go v.listen()
	return v
}

func NewLazyCircleVisualizer(ampInput chan []float64) m.Visualizer {
	return NewCircleVisualizer(10, ampInput)
}

func (v *circleVisualizer) listen() {
	for {
		amplitudes := <-v.ampInput

		cf := fft.FFTReal(amplitudes)
		cf = cf[0 : len(cf)/2]
		frequencies := make([]float64, len(cf))
		for i := range cf {
			frequencies[i] = math.Abs(real(cf[i]))
		}
		v.buffer.Push(frequencies)
		//peaks.Get()
	}
}

type circleVisualizer struct {
	ampInput chan []float64
	buffer   *util.FrequenciesBuffer
	//poly        *eb.Image
	poly        *Polygon
	colorFloats [][]float64
	colors      []color.Color
}

func (v *circleVisualizer) Draw(screen *eb.Image) {
	var (
		frequencies   = v.buffer.GetAverage()
		w, h          = screen.Size()
		width, height = float64(w), float64(h)
		groups        = util.GroupFrequencies(3, frequencies)
		_             = groups
	)

	screen.Fill(color.RGBA{
		R: 0x00,
		G: 0x00,
		B: 0x00,
		A: 0xff,
	})
	ops := eb.DrawImageOptions{}
	ops.ColorM.Scale(0, 0, 0, 1)
	ops.ColorM.Translate(v.colorFloats[0][0]/0xff, v.colorFloats[0][1]/0xff, v.colorFloats[0][2]/0xff, 0)
	ops.GeoM.Translate(width/2, height/2)
	//ops = CenterOps(v.poly.GetImg(), ops)
	SpiralNestPolygons(screen, v.poly, 3, .9,1, ops)
	//DrawImgFromCenter(screen,v.poly.GetImg(),width/2,height/2,1,ops)
	/*
		for i, g :=range groups{
			//centerX,centerY := width/float64(len(groups))*float64(i),height/2
			centerX,centerY := width/2,height/2

			ops := eb.DrawImageOptions{}
			colors := v.colorFloats[i%len(v.colorFloats)]
			_=colors
			ops.ColorM.Scale(0,0,0,1)
			//ops.ColorM.Translate(float64(0xff)/0xff,float64(0x0)/0xff,float64(0xff)/0xff,0)
			ops.ColorM.Translate(colors[0]/0xff,colors[1]/0xff,colors[2]/0xff,0)
			//ops.ColorM.Apply(v.colors[i%len(v.colors)])
			//scale :=(normalizeFrequency(g,height)/2)*scaleCoefficient
			var	scale float64 = .5
			_=g
			//ops.GeoM.Scale(scale,scale)
			//cirleW, circleH :=v.poly.Size()
			//ops.GeoM.Translate(centerX-scale*float64(cirleW)/2,centerY-scale*float64(circleH)/2)

			//ops=CenterOps(v.poly,centerX,centerY,scale,ops)
			//screen.DrawImage(v.poly,&ops)
			//ops.ColorM.Scale(0,0,0,1)
			//ops=CenterOps(v.poly,centerX,centerY,scale*.9,ops)
			//screen.DrawImage(v.poly,&ops)

			DrawImgFromCenter(screen,v.poly,centerX,centerY,scale,ops)
			ops.ColorM.Scale(0,0,0,1)
			DrawImgFromCenter(screen,v.poly,centerX,centerY,scale*.9,ops)

			ops.ColorM.Translate(colors[0]/0xff,colors[1]/0xff,colors[2]/0xff,0)
			//ops.GeoM.Translate(?)
			ops.GeoM.Rotate(math.Pi*2/(6*2))
			DrawImgFromCenter(screen,v.poly,centerX,centerY,scale*.9,ops)
			//ops.GeoM.Scale(.9,.9)
			//ops.ColorM.Scale(0,0,0,1)
			//screen.DrawImage(v.poly,&ops)

		}
	*/
}

func DrawHollowImg(screen, img *eb.Image) {

}

func DrawImgFromCenter(screen, img *eb.Image, scale float64, ops eb.DrawImageOptions) {
	cirleW, circleH := img.Size()
	ops.GeoM.Scale(scale, scale)
	ops.GeoM.Translate(-scale*float64(cirleW)/2, -scale*float64(circleH)/2)
	screen.DrawImage(img, &ops)
}

// problem is:
// move half back
// rotate
// move half back again

// x and y are the coordinates WITHIN SCREEN that you want to draw
func DrawImgFromCenter_old(screen, img *eb.Image, x, y, scale float64, ops eb.DrawImageOptions) {
	cirleW, circleH := img.Size()
	ops.GeoM.Scale(scale, scale)
	ops.GeoM.Translate(x-scale*float64(cirleW)/2, y-scale*float64(circleH)/2)
	screen.DrawImage(img, &ops)
}

func CenterOps(img *eb.Image, ops eb.DrawImageOptions) eb.DrawImageOptions {
	cirleW, circleH := img.Size()
	local := eb.DrawImageOptions{
		ColorM: ops.ColorM,
	}
	local.GeoM.Translate(-float64(cirleW)/2, -float64(circleH)/2)
	//local.GeoM.Rotate(math.Pi/2)
	//local.GeoM.Scale(.9, .9)
	local.GeoM.Concat(ops.GeoM)
	return local
}

func SpiralNestPolygons(screen *eb.Image, poly *Polygon, depth int, scale,rotation float64, ops eb.DrawImageOptions) {
	if depth < 1 {
		return
	}

	r,s:=CalcPolyRotationScale(poly)
	rotation+=r
	scale*=s

	cirleW, circleH := poly.GetImg().Size()
	local := eb.DrawImageOptions{
		ColorM: ops.ColorM,
	}
	rot := eb.GeoM{}
	rot.Translate(-float64(cirleW)/2, -float64(circleH)/2)
	//rot.Translate(float64(cirleW)/2, float64(circleH)/2)
	rotation+=math.Pi / 4
	rot.Rotate(rotation)
	// translate in
	//local.GeoM.Translate(-float64(cirleW)/2, -float64(circleH)/2)
	// scale
	scale*=scale
	rot.Scale(scale, scale)

	// translate out

	local.GeoM.Concat(rot)
	local.GeoM.Concat(ops.GeoM)
	screen.DrawImage(poly.GetImg(), &local)
	//local.GeoM.Rotate(-math.Pi/4)
	//local.GeoM.Translate(float64(cirleW)/2, float64(circleH)/2)

	//ops.GeoM.Concat(local.GeoM)

	//// before the func, coordinates need to be localized
	//local := eb.DrawImageOptions{
	//	ColorM: ops.ColorM,
	//}
	//local.GeoM.Scale(scale,scale)
	//local.GeoM.Rotate(math.Pi/4)
	//local.GeoM.Concat(ops.GeoM)
	////local.GeoM.Translate(float64(cirleW)/2, float64(circleH)/2)

	SpiralNestPolygons(screen, poly, depth-1, scale,rotation, ops)
	//DrawImgFromCenter(screen, poly.GetImg(), 1, ops)
	//ops = RotatePolyOps(poly, ops)
	//SpiralNestPolygons(screen, poly, depth-1, ops)
}

func CalcPolyRotationScale(poly *Polygon)(rotation,scale float64){
	scale = math.Sqrt(math.Pow((.5+.5*math.Cos(poly.GetTheta())), 2) + math.Pow(.5*math.Sin(poly.GetTheta()), 2))
	rotation = (math.Pi - poly.GetTheta()) / 2 //math.Pi * poly.GetTheta() / 2
	fmt.Println(scale,rotation)
	return
}

func RotatePolyOps(poly *Polygon, ops eb.DrawImageOptions) eb.DrawImageOptions {
	newL := 1/math.Sqrt(math.Pow((.5+.5*math.Cos(poly.GetTheta())), 2) + math.Pow(.5*math.Sin(poly.GetTheta()), 2))
	fmt.Println(newL)
	rotateBy := (math.Pi - poly.GetTheta()) / 2 //math.Pi * poly.GetTheta() / 2
	imgX, imgY := poly.img.Size()

	deltaX := newL*math.Sin(rotateBy) - float64(imgX)/2
	deltaY := newL*math.Cos(rotateBy) - float64(imgY)/2
	//ops.GeoM.Rotate(rotateBy)
	ops.GeoM.Translate(deltaX, deltaY)

	local := eb.DrawImageOptions{}
	local.GeoM.Rotate(rotateBy)
	//ops.GeoM.Concat(local.GeoM)
	local.GeoM.Concat(ops.GeoM)

	//ops.GeoM.Scale(.9,.9)
	return ops
}

func NewPolygon(n int, radius float64) *Polygon {
	if n < 3 {
		return nil // ?
	}
	p := &Polygon{
		n:      float64(n),
		radius: radius,
		img:    CreatePolyImgHollow(radius, n, .7),
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
