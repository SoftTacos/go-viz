package analyzer

import (
	"github.com/mjibson/go-dsp/fft"
	"math"
)

type ProcessFunc func(in []float64)(out []float64)

type pipeline struct {
	input <-chan []float64
	operations []ProcessFunc
	outputs []chan<- []float64
}

func (v *pipeline)listen(){
	for {
		input:=<-v.input
		for _,operation:=range v.operations {
			input = operation(input)
		}
		for _,output:=range v.outputs{
			output<-input
		}
	}
}

func ExtractFrequencies(input []float64)(out []float64){
	cf:=fft.FFTReal(input)
	out = make([]float64,len(cf))
	for i:=range cf {
		out[i] = math.Abs(real(cf[i]))
	}
	return
}