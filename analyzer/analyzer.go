package analyzer

import (
	"github.com/mjibson/go-dsp/fft"
	"github.com/softtacos/go-visualizer/util"
	"math"
)

type AnalyzerFunc func(amplitudes []float64)

func NewAnalyzer(ampIn <-chan []float64, out chan<- []float64, beatDetector *util.BeatDetector) *Analyzer {
	return &Analyzer{
		ampIn:        ampIn,
		out:          out,
		beatDetector: beatDetector,
	}
}

type Analyzer struct {
	ampIn        <-chan []float64
	out          chan<- []float64
	beatDetector *util.BeatDetector
}

var max float64

func (a *Analyzer) Start() {
	go func() {
		for {
			amplitudes := <-a.ampIn
			cf := fft.FFTReal(amplitudes)
			cf = cf[1 : len(cf)/2]
			frequencies := make([]float64, len(cf))
			for i := range cf {
				//f := math.Abs(real(cf[i]))
				//if f > max {
				//	max = f
				//}
				//f /= max
				frequencies[i] =  math.Abs(real(cf[i]))
			}
			a.out <- frequencies
		}
	}()
}

//func GenerateFftAnalyzerAnalyzer(out chan []float64) AnalyzerFunc {
//	return func(amplitudes []float64) {
//		// TODO: overflow
//		cf := fft.FFTReal(amplitudes)
//		cf = cf[0 : len(cf)/2]
//		frequencies := make([]float64, len(cf))
//		for i := range cf {
//			frequencies[i] = real(cf[i])
//		}
//		out <- frequencies
//	}
//}
//
//func GenerateBeatDetectorAnalyzer(b *util.BeatDetector) AnalyzerFunc {
//	return func(amplitudes []float64) {
//		b.Push(amplitudes)
//	}
//}
//
//func GenerateFftBeatDetectorAnalyzer(out chan []float64, b *util.BeatDetector) AnalyzerFunc {
//	return func(amplitudes []float64) {
//		cf := fft.FFTReal(amplitudes)
//		cf = cf[0 : len(cf)/2]
//		frequencies := make([]float64, len(cf))
//		for i := range cf {
//			frequencies[i] = real(cf[i])
//		}
//		out <- frequencies
//		b.Push(frequencies[0 : len(frequencies)/10])
//	}
//}
