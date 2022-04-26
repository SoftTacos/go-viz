package analyzer

import "github.com/mjibson/go-dsp/fft"

func NewAnalyzer(ampIn <-chan []float64, ampOut, freqOut chan<- []float64, beatDetectedCallback func()) *analyzer {
	return &analyzer{
		ampIn:                ampIn,
		ampOut:               ampOut,
		freqOut:              freqOut,
		beatDetectedCallback: beatDetectedCallback,
	}
}

type analyzer struct {
	ampIn <-chan []float64
	//outputs []chan []float64
	ampOut               chan<- []float64
	freqOut              chan<- []float64
	beatDetectedCallback func()
}

func (a *analyzer) Start() {
	for {
		// TODO: figure out a way to make this pass-through to the visualizer
		amplitudes := <-a.ampIn

		cf := fft.FFTReal(amplitudes)
		frequencies := make([]float64, len(cf))
		for i := range cf {
			frequencies[i] = real(cf[i])
		}

		a.ampOut <- amplitudes
		a.freqOut <- frequencies
	}
}
