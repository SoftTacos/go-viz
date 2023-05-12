package util

import (
	"fmt"
	"sync"
	"time"
)

func NewBeatDetector(windowSize, samplesPerSecond, maxBpm int, callback func()) *BeatDetector {
	const buffer = 1.0
	maxBeatsPerSecond := float64(maxBpm) / 60.0
	minPeriod := time.Duration(1.0 / maxBeatsPerSecond * 1000000000) // base unit of time.Duration is a nanosecond, * 1000000000 turns into milliseconds
	samplesPerMinBeat := buffer * float64(maxBeatsPerSecond) * float64(samplesPerSecond)
	len := 90 //int((samplesPerMinBeat + 1.0) / float64(windowSize))
	stepTime := 1.0 / float64(samplesPerSecond)
	fmt.Println("LEN", len)
	chunks := make([][]float64, len)
	for i := range chunks {
		chunks[i] = make([]float64, windowSize)
	}
	//var variances = make([]float64, len)
	//for i := range chunks {
	//	variances[i] = .1
	//}
	fmt.Println(maxBeatsPerSecond, minPeriod, samplesPerMinBeat, len)
	return &BeatDetector{
		peaks:      NewCircularPeakSlice(10), // TODO: how many peaks do we want?
		chunks:     chunks,
		means:      make([]float64, len),
		variances:  make([]float64, len),
		windowSize: windowSize,
		len:        len,
		lastPeak: peak{
			timestamp: time.Now(),
			value:     0,
		},
		stepTimeSeconds: stepTime,
		minPeriod:       minPeriod,
		sampleRate:      float64(samplesPerSecond),
		mutex:           &sync.Mutex{},
		callback:        callback,
	}
}

type BeatDetector struct {
	sampleRate      float64
	minPeriod       time.Duration
	peaks           *CircularPeakSlice
	windowSize      int
	lastPeak        peak
	stepTimeSeconds float64
	callback        func()
	mutex           *sync.Mutex
	front           int // index of first/latest value
	len             int
	chunks          [][]float64
	means           []float64
	variances       []float64

	//mean   float64
	//stdDev float64
}

type peak struct {
	value     float64
	timestamp time.Time
}

func (b *BeatDetector) SetCallback(callback func()) {
	b.callback = callback
}

func (b *BeatDetector) Listen(amps chan []float64) {
	for {
		b.Push(<-amps)
	}
}

//func (b *BeatDetector) Push(args []float64) {
//	b.mutex.Lock()
//	defer b.mutex.Unlock()
//
//	// get mean, max, std dev of args
//	argMean, argVariance, argMax := CalcArgStats(args)
//	b.chunks[b.front] = args
//	b.means[b.front] = argMean
//	b.variances[b.front] = argVariance
//	b.front = (b.front + 1) % b.len
//	// if max > coeff*global-stddev, it's a beat?
//
//	now := time.Now()
//	timeSinceLastPeak := now.Sub(b.lastPeak.timestamp)
//	if timeSinceLastPeak > b.minPeriod {
//
//		mean, variance := CalcMultiSetStats(float64(b.windowSize), b.means, b.variances)
//		sdev := math.Sqrt(variance)
//		cutoff := mean + (sdev)*4
//
//		if argMax > cutoff {
//			fmt.Printf("BEAT %v m:%.2f sd:%.2f max:%.2f co:%.2f r:%.4f\n", timeSinceLastPeak, mean, sdev, argMax, cutoff,argMax-cutoff)
//			b.lastPeak.value = argMax
//			b.lastPeak.timestamp = now
//			b.callback()
//		} else if timeSinceLastPeak > time.Second {
//			fmt.Printf("m:%.2f sd:%.2f max:%.2f co:%.2f r:%.4f\n",  mean, sdev, argMax, cutoff,argMax-cutoff)
//		}
//	}
//
//	return
//}

func CalcMeanVariance(set []float64) (mean, variance float64) {
	var (
		dev  float64
		flen = float64(len(set))
	)
	for _, arg := range set {
		mean += arg
	}
	mean /= float64(len(set))

	for _, arg := range set {
		dev = arg - mean
		variance += dev * dev
	}
	variance /= flen
	return
}

func CalcArgStats(set []float64) (mean, variance, max float64) {
	var (
		dev  float64
		flen = float64(len(set))
	)
	for _, arg := range set {
		mean += arg
		if arg > max {
			max = arg
		}
	}
	mean /= float64(len(set))

	for _, arg := range set {
		dev = arg - mean
		variance += dev * dev
	}
	variance /= flen
	return
}

// assumes means and sdevs have the same lengths and that the lengths are greater than 0
// assumes all sets have length of n
func CalcMultiSetStats(n float64, means, variances []float64) (mean, variance float64) {
	num := n
	mean = means[0]
	variance = variances[0] //calc2SetVariance(n,means[0],means[1],variances[0],variances[1])
	for i := 0; i < len(means)-1; i++ {
		variance = calc2SetVariance(num, n, mean, means[i+1], variance, variances[i+1])
		mean = calc2SetMean(num, n, mean, means[i+1])
		num += n
	}

	return
}

func calc2SetMean(n1, n2, m1, m2 float64) float64 {
	// [n1 /(n1+n2)]*Xbar1 + [n2 /(n1+n2)]*Xbar2
	return (n1/(n1+n2))*m1 + (n2/(n1+n2))*m2
}

func calc2SetVariance(n1, n2, m1, m2, v1, v2 float64) float64 {
	// http://www.talkstats.com/threads/standard-deviation-of-multiple-sample-sets.7130/
	// [ n1^2*Var1 + n2^2*Var2 – n1*Var1 – n1*Var2 – n2*Var1 - n2*Var2 + n1*n2*Var1 + n1*n2*Var2 +n1*n2*(Xbar1 – Xbar2)^2 ] / [ (n1+n2-1)*(n1+n2) ]
	//n*n*variances[0] + n*n*variances[1] - n*variances[0] - n*variances[1] - n*variances[0] - n*variances[1] + n*n*variances[0] + n*n*variances[1] + n*n*(means[0]-means[1])*(means[0]-means[1])) / ((2*n - 1) * (2 * n))
	return (n1*n1*v1 + n2*n2*v2 - n1*v1 - n1*v2 - n2*v1 - n2*v2 + n1*n2*v1 + n1*n2*v2 + n1*n2*(m1-m2)*(m1-m2)) / ((n1+n2 - 1) * (n1 + n2))
}

const (
	decayScaleBase float64 = .9
	decayScaleCoeff float64 = 1.5
)

func (b *BeatDetector) Push(args []float64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var (
		//slopes      = make([]float64, len(args))
		//previousVal = b.buffer[b.front]
		//argMean  float64
		max      float64
	)

	// find max amplitude of the incoming sample
	for _, arg := range args {
		//arg = math.Abs(arg)
		//slopes[i] = (arg - previousVal) / b.stepTimeSeconds
		if arg*arg > max {
			max = arg * arg
		}
		//argMean += arg
	}
	//argMean /= float64(len(args))
	//b.front = (b.front + len(args)) % b.len

	now := time.Now()
	timeSinceLastPeak := now.Sub(b.lastPeak.timestamp)

	if timeSinceLastPeak > b.minPeriod  {
		var decayScale float64 = decayScaleBase - (float64(timeSinceLastPeak)-decayScaleCoeff*float64(b.minPeriod))/1000000000
		if max > b.lastPeak.value*decayScale {
			fmt.Printf("BEAT DETECTED: %+v\tlast:%.2f\tt:%.2f\n",timeSinceLastPeak, b.lastPeak.value,b.lastPeak.value*decayScale)
			b.lastPeak = peak{
				timestamp: now,
				value:     max,
			}
			b.callback()
			//b.peaks.Push(b.lastPeak)
		}
	}

}
