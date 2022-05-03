package visualizers

import (
	eb "github.com/hajimehoshi/ebiten/v2"
	m "github.com/softtacos/go-visualizer/model"
	"image/color"
	"math"
)

func NewForkingVisualizerConstructor(windowSize int, frequencyInput chan []float64) m.Visualizer {
	return newDefaultForkingVisualizer(windowSize, frequencyInput)
}

func newDefaultForkingVisualizer(windowSize int, frequencyInput chan []float64) *forkingVisualizer {
	v := &forkingVisualizer{
		frequencyInput: frequencyInput,
		fork:           NewFork(2,math.Pi / 6.0),
	}
	return v
}

type forkingVisualizer struct {
	frequencyInput chan []float64
	fork           *fork
}

func (v *forkingVisualizer) Draw(screen *eb.Image) {
	RecursiveForks(screen,v.fork)
}

func (v *forkingVisualizer) BeatCallback() {

}

func RecursiveForks(screen *eb.Image,fork *fork) {
	sX,sY:=screen.Size()
	_=sX
	imgX,imgY:=fork.GetImg().Size()
	_=imgX
	op :=eb.DrawImageOptions{}
	op.GeoM.Translate(0,float64(sY)/2)
	screen.DrawImage(fork.GetImg(),&op)

	fOp:=eb.DrawImageOptions{}
	fOp.GeoM.Translate(0,-float64(imgY)/2)
	fOp.GeoM.Rotate(-fork.GetTheta())
	fOp.GeoM.Translate(fork.rectX+fork.rectX*math.Cos(fork.GetTheta()),fork.rectX*math.Sin(fork.GetTheta()))
	fOp.GeoM.Concat(op.GeoM)
	screen.DrawImage(fork.GetImg(),&fOp)

	fOp=eb.DrawImageOptions{}
	fOp.GeoM.Translate(0,float64(imgY)/2)
	fOp.GeoM.Rotate(fork.GetTheta())
	fOp.GeoM.Translate(fork.rectX+fork.rectX*math.Cos(fork.GetTheta()),-fork.rectX*math.Sin(fork.GetTheta()))
	fOp.GeoM.Concat(op.GeoM)
	screen.DrawImage(fork.GetImg(),&fOp)
}

func NewFork(numEnds int,theta float64) *fork {
	return &fork{
		numEnds: numEnds,
		theta: theta,
		rectX: 40,
		rectY: 10,
		img:     NewForkImage(numEnds,theta),
	}
}

type fork struct {
	numEnds int
	theta float64 // angle between the midline and the fork
	rectX,rectY float64
	img *eb.Image
}

func (f *fork)GetImg()*eb.Image{
	return f.img
}

func (f *fork)GetTheta()float64{
	return f.theta
}

func NewForkImage(numForks int,theta float64) *eb.Image {
	baseX, baseY := 80, 80
	baseXf, baseYf := float64(baseX), float64(baseY)
	baseImg := eb.NewImage(baseX, baseY)
	baseImg.Fill(color.Transparent)
	_ = baseXf
	lineX, lineY := 40, 10
	lineXf, lineYf := float64(lineX), float64(lineY)
	_ = lineXf
	rect := Rect(lineX, lineY)
	op := eb.DrawImageOptions{}
	op.GeoM.Translate(0, baseYf/2-lineYf/2)
	baseImg.DrawImage(rect, &op)

	op.GeoM.Translate(lineXf, 0)
	rot := eb.DrawImageOptions{}
	rot.GeoM.Rotate(theta)
	rot.GeoM.Concat(op.GeoM)
	baseImg.DrawImage(rect, &rot)

	rot = eb.DrawImageOptions{}
	rot.GeoM.Translate(0,-lineYf)
	rot.GeoM.Rotate(-theta)
	rot.GeoM.Concat(op.GeoM)
	rot.GeoM.Translate(0,lineYf)
	baseImg.DrawImage(rect, &rot)

	return baseImg
}

func CalcForkOps(numForks int, theta float64)(ops []eb.DrawImageOptions){
	return
}

func Rect(x, y int) *eb.Image {
	i := eb.NewImage(x, y)
	i.Fill(color.White)
	return i
}
