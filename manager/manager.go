package manager

import (
	pa "github.com/gordonklaus/portaudio"
	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	a "github.com/softtacos/go-visualizer/analyzer"
	m "github.com/softtacos/go-visualizer/model"
	s "github.com/softtacos/go-visualizer/stream"
	"github.com/softtacos/go-visualizer/util"
	"log"
	"math"
	"math/cmplx"
	"os"
	"os/signal"
)

const (
	windowSize    = 256 // number of samples per []float64
	defaultBuffer = 2
)

// TODO: this is hideous, clean it up
func NewVisManager(constructors []m.VisualizerConstructor) *manager {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	streamOutput := make(chan []float64, defaultBuffer)
	frequencyOutput := make(chan []float64, defaultBuffer)

	streamer := s.NewStreamer(windowSize, streamOutput, defaultBuffer)
	viz := constructors[0](windowSize, frequencyOutput)
	var beatDetector *util.BeatDetector
	beatDetector = util.NewBeatDetector(48000, 190, nil)
	var analyzer *a.Analyzer = a.NewAnalyzer(streamOutput,
		a.GenerateFftAnalyzerAnalyzer(frequencyOutput),
		a.GenerateBeatDetectorAnalyzer(beatDetector))

	m := &manager{
		analyzer:          analyzer,
		constructors:      constructors,
		currentVisualizer: viz,
		streamer:          streamer,
		streamerOutput:    streamOutput,
	}

	beatDetector.SetCallback(m.BeatCallback)
	return m
}

type manager struct {
	audio             chan []float64
	analyzer          *a.Analyzer
	streamer          *s.Streamer
	streamerOutput    chan []float64
	vIndex            int
	currentVisualizer m.Visualizer
	constructors      []m.VisualizerConstructor
}

func (v *manager) Start() {
	v.streamer.Start()
	v.analyzer.Start()
}

func (v *manager) Stop() {
	var err error
	log.Println("attempting to terminate")
	if err = pa.Terminate(); err != nil {
		log.Println("failed to terminate:", err)
	}
	log.Println("shutting down")
}

func (v *manager) BeatCallback() {
	v.currentVisualizer.BeatCallback()
}

func (v *manager) Update() (err error) {
	// TODO: listen for input event to change visualizer
	if inpututil.IsKeyJustPressed(eb.KeyLeft) {
		v.ChangeVisualizers(-1)
	} else if inpututil.IsKeyJustPressed(eb.KeyRight) {
		v.ChangeVisualizers(1)
	}
	//l := len(v.input)
	//frequencies :=make([]complex128,l)
	//DFFT64(<-v.input,frequencies,l,1)
	//v.currentMutex.Lock()
	//v.current = make([]float64,l)
	//v.currentMutex.Unlock()
	return
}

func (v *manager) ChangeVisualizers(dir int) {
	newViz := v.constructors[(v.vIndex+dir)%len(v.constructors)](v.streamer.GetWindowSize(), v.streamerOutput)
	v.currentVisualizer = newViz
}

func (v *manager) Draw(screen *eb.Image) {
	v.currentVisualizer.Draw(screen)
}

func (v *manager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
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
