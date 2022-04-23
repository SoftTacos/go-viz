package visualizers

import (
	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/mjibson/go-dsp/fft"
	m "github.com/softtacos/go-visualizer/model"
	"github.com/softtacos/go-visualizer/util"
	"math"
	"sync"

	//"github.com/goccmack/godsp/peaks"
)

func NewPaddedVisualizer(bufferSize int,ampInput chan []float64)m.Visualizer{
	v:=&paddedVisualizer{
		buffer: util.NewFrequencyBuffer(bufferSize),
		ampInput:ampInput,
		ampMutex: &sync.Mutex{},
	}
	go v.listen()
	return v
}

func NewLazyPaddedVisualizer(ampInput chan []float64)m.Visualizer{
	return NewPaddedVisualizer(7,ampInput)
}

func (v *paddedVisualizer)listen(){
	for {
		// TODO: extract out into a pipeline once we can test this out on a pi
		amplitudes:=<-v.ampInput

		v.ampMutex.Lock()
		v.currentAmp = amplitudes
		v.ampMutex.Unlock()

		cf:=fft.FFTReal(amplitudes)
		frequencies := make([]float64,len(cf))
		for i:=range cf {
			frequencies[i] = math.Abs(real(cf[i]))
		}
		v.buffer.Push(frequencies)
		//peaks.Get()
	}
}

type paddedVisualizer struct {
	ampInput chan []float64
	currentAmp []float64
	ampMutex *sync.Mutex
	buffer *util.FrequenciesBuffer
}

func (v *paddedVisualizer)Draw(screen *eb.Image){
	//v.ampMutex.Lock()
	//DrawFrequencies(screen,v.currentAmp)
	//v.ampMutex.Unlock()
	DrawFrequencies(screen,v.buffer.GetAverage())
}

//func GetPeaks(x []float64, sep int) []int {
//	pks := []int{}
//	for i := range x {
//		if isMax(i, i-sep, i+sep, x) {
//			pks = append(pks, i)
//		}
//	}
//	return pks
//}
//
//func isMax(i, min, max int, x []float64) bool {
//	if min < 0 {
//		min = 0
//	}
//	if max > len(x) {
//		max = len(x)
//	}
//	for j := min; j < i; j++ {
//		if x[j] >= x[i] {
//			return false
//		}
//	}
//	for j := i + 1; j < max; j++ {
//		if x[j] > x[i] {
//			return false
//		}
//	}
//	return true
//}
//
//func getWindow(i, sep int, x []float64) (min, max int) {
//	min, max = i-sep, i+sep
//	if min < 0 {
//		min = 0
//	}
//	if max > len(x) {
//		max = len(x)
//	}
//	return
//}
//
//func getEmptyPeaks(n int) []int {
//	epks := make([]int, n)
//	for i := range epks {
//		epks[i] = -1 // aka empty
//	}
//	return epks
//}
//
//func getSortedIndices(x []float64) []int {
//	idx := godsp.Range(len(x))
//	sort.SliceStable(idx, func(i, j int) bool { return x[i] > x[j] })
//	return idx
//}
//
//func markNeighbours(xi, sep int, pks []int) {
//	min := xi - sep
//	if min < 0 {
//		min = 0
//	}
//	max := xi + sep
//	if max > len(pks) {
//		max = len(pks)
//	}
//	for i := min; i < max; i++ {
//		pks[i] = xi
//	}
//}
//func avgFilter(start,end int,frequencies []float64)(avg float64){
//	low:=frequencies[start:end]
//	for _,l:=range low{
//		avg+=l
//	}
//	avg/=float64(len(low))
//	return
//}

