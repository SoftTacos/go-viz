package analyzer

import (
	"github.com/mjibson/go-dsp/fft"
	"github.com/softtacos/go-visualizer/util"
)

type AnalyzerFunc func(amplitudes []float64)

func NewAnalyzer(ampIn <-chan []float64, funcs ...AnalyzerFunc) *Analyzer {
	return &Analyzer{
		ampIn: ampIn,
		funcs: funcs,
	}
}

type Analyzer struct {
	ampIn <-chan []float64
	funcs []AnalyzerFunc
}

func (a *Analyzer) Start() {
	go func() {
		for {
			amplitudes := <-a.ampIn
			for _, fun := range a.funcs {
				fun(amplitudes)
			}
		}
	}()
}

func (a *Analyzer) Analyze(amplitudes []float64) {
	for _, fun := range a.funcs {
		fun(amplitudes)
	}
}

func GenerateFftAnalyzerAnalyzer(out chan []float64) AnalyzerFunc {
	return func(amplitudes []float64) {
		// TODO: overflow
		cf := fft.FFTReal(amplitudes)
		cf = cf[0 : len(cf)/2]
		frequencies := make([]float64, len(cf))
		for i := range cf {
			frequencies[i] = real(cf[i])
		}
		out <- frequencies
	}
}

func GenerateBeatDetectorAnalyzer(b *util.BeatDetector) AnalyzerFunc {
	return func(amplitudes []float64) {
		b.Push(amplitudes)
	}
}
