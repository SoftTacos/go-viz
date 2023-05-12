package util

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func NewFrequencyBeatDetector(windowSize, samplesPerSecond, maxBpm int, callback func()) *FrequencyBeatDetector {
	const buffer = 2.0
	maxBeatsPerSecond := float64(maxBpm) / 60.0
	minPeriod := time.Duration(1.0 / maxBeatsPerSecond * 1000000000) // base unit of time.Duration is a nanosecond, * 1000000000 turns into milliseconds
	samplesPerMinBeat := buffer * float64(maxBeatsPerSecond) * float64(samplesPerSecond)
	len := int((samplesPerMinBeat + 1.0) / float64(windowSize))
	recentLen := 5
	stepTime := 1.0 / float64(samplesPerSecond)
	fmt.Printf("%d, %d\n", len, windowSize)
	chunks := make([][]float64, len)
	for i := range chunks {
		chunks[i] = make([]float64, windowSize)
	}
	return &FrequencyBeatDetector{
		peaks:         NewCircularPeakSlice(10), // TODO: how many peaks do we want?
		samples:       make([]float64, len),
		recentSamples: make([]float64, recentLen),
		windowSize:    windowSize,
		len:           len,
		recentLen:     recentLen, //len / 10,
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

type FrequencyBeatDetector struct {
	sampleRate             float64
	minPeriod              time.Duration
	peaks                  *CircularPeakSlice
	windowSize             int
	lastPeak               peak
	stepTimeSeconds        float64
	callback               func()
	mutex                  *sync.Mutex
	len, recentLen         int
	front, recentFront     int // index of first/latest value
	samples, recentSamples []float64
	mean                   float64
	variance               float64
}

func (b *FrequencyBeatDetector) SetCallback(callback func()) {
	b.callback = callback
}

const coeff float64 = 1.5

var lastSampleTime time.Time

func (b *FrequencyBeatDetector) Push(frequencies []float64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	now := time.Now()
	//diff := now.Sub(lastSampleTime)
	//fmt.Println(1.0/float64(diff)*1000000000)
	//lastSampleTime=now

	//var freq = GroupFrequencies(30, frequencies)[0]
	var freq = frequencies[0]
	// get mean, max, std dev of args

	// if max > coeff*global-stddev, it's a beat?
	b.samples[b.front] = freq
	b.front = (b.front + 1) % b.len

	b.recentSamples[b.recentFront] = freq
	b.recentFront = (b.recentFront + 1) % b.recentLen

	timeSinceLastPeak := now.Sub(b.lastPeak.timestamp)
	if timeSinceLastPeak > b.minPeriod {
		var max float64
		var maxIndex,before,after int
		for i, s := range b.recentSamples {
			if s > max {
				max = s
				maxIndex = i
			}
		}
		//after = (maxIndex+1)%b.recentLen
		//before = maxIndex-1
		//if before < 0{
		//	before = b.recentLen-1
		//}
		_=maxIndex
		_=before
		_=after

		b.mean, b.variance = CalcMeanVariance(b.samples)
		sdev := math.Sqrt(b.variance)
		cutoff := b.mean + sdev * coeff
		if max > cutoff {
			fmt.Printf("BEAT %v m:%.2f sd:%.2f max:%.2f co:%.2f\n", timeSinceLastPeak, b.mean, sdev, max, cutoff)
			b.lastPeak.value = max
			b.lastPeak.timestamp = now
			b.callback()
		} else if timeSinceLastPeak > time.Second {
			fmt.Printf("m:%.2f sd:%.2f max:%.2f co:%.2f\n", b.mean, sdev, max, cutoff)
		}
	}

	return
}
