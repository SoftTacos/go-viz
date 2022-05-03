package util

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func NewBeatDetector(windowSize, samplesPerSecond, maxBpm int, callback func()) *BeatDetector {
	const buffer = 2.0
	maxBeatsPerSecond := float64(maxBpm) / 60.0
	minPeriod := time.Duration(1.0 / maxBeatsPerSecond * 1000000000) // base unit of time.Duration is a nanosecond, * 1000000000 turns into milliseconds
	samplesPerMinBeat := buffer * float64(maxBeatsPerSecond) * float64(samplesPerSecond)
	len := int(samplesPerMinBeat) + 1
	stepTime := 1.0 / float64(samplesPerSecond)
	return &BeatDetector{
		peaks:        NewCircularPeakSlice(10), // TODO: how many peaks do we want?
		buffer:       make([]float64, len),
		len:          len,
		lastBeatTime: time.Now(),
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
	lastBeatTime    time.Time
	minPeriod       time.Duration
	peaks           *CircularPeakSlice
	windowSize      int
	lastPeak        peak
	stepTimeSeconds float64
	callback        func()
	mutex           *sync.Mutex
	front           int // index of first/latest value
	len             int
	buffer          []float64

	mean   float64
	stdDev float64
}

type peak struct {
	value     float64
	timestamp time.Time
}

func (b *BeatDetector) SetCallback(callback func()) {
	b.callback = callback
}

func (b *BeatDetector) Push(args []float64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var (
		argMean  float64
		maxIndex = -1
		max      float64
	)

	for i, arg := range args {
		arg = math.Abs(arg)

		if arg*arg > max {
			max = arg * arg
			maxIndex = i
		}
		b.buffer[(b.front+i)%b.len] = arg
		argMean += arg
		// TODO: take the rolling mean

		// TODO: wtf do with std dev
	}
	argMean /= float64(len(args))
	b.front = (b.front + len(args)) % b.len

	now := time.Now()
	timeSinceLastPeak := now.Sub(b.lastPeak.timestamp)
	if timeSinceLastPeak > b.minPeriod && maxIndex > -1 {
		var decayScale float64 = .9 - (float64(timeSinceLastPeak)-1.5*float64(b.minPeriod))/1000000000
		if max > b.lastPeak.value*decayScale {
			fmt.Println("BEAT DETECTED:", b.lastPeak.value, timeSinceLastPeak, decayScale)
			b.lastPeak = peak{
				timestamp: now,
				value:     max,
			}
			b.callback()
			//b.peaks.Push(b.lastPeak)
		}
	}

	// track the std dev
	rollingMean := b.mean * (float64(b.len-len(args)) / float64(b.len))
	b.mean = rollingMean + argMean*(float64(len(args))/float64(b.len))

	b.stdDev = 0
	b.mean /= float64(b.len)
	for j := 0; j < b.len; j++ {
		b.stdDev += math.Pow(b.buffer[j]-b.mean, 2)
	}
	b.stdDev = math.Sqrt(b.stdDev / float64(b.len))

}


//func (b *BeatDetector) Push(args ...float64) {
//	b.mutex.Lock()
//	defer b.mutex.Unlock()
//
//	var (
//		//slopes      = make([]float64, len(args))
//		//previousVal = b.buffer[b.front]
//		argMean  float64
//		maxIndex = -1
//		max      float64
//	)
//
//	for i, arg := range args {
//		arg = math.Abs(arg)
//		//slopes[i] = (arg - previousVal) / b.stepTimeSeconds
//
//		if arg*arg > max {
//			max = arg * arg
//			maxIndex = i
//		}
//		b.buffer[(b.front+i)%b.len] = arg
//		argMean += arg
//		// TODO: take the rolling mean
//
//		// TODO: wtf do with std dev
//	}
//	argMean /= float64(len(args))
//	b.front = (b.front + len(args)) % b.len
//
//	now := time.Now()
//	timeSinceLastPeak := now.Sub(b.lastPeak.timestamp)
//	if timeSinceLastPeak > b.minPeriod && maxIndex > -1 {
//		var decayScale float64 = .9 - (float64(timeSinceLastPeak)-1.5*float64(b.minPeriod))/1000000000
//		if max > b.lastPeak.value*decayScale {
//			fmt.Println("BEAT DETECTED:", b.lastPeak.value, timeSinceLastPeak, decayScale)
//			b.lastPeak = peak{
//				timestamp: now,
//				value:     max,
//			}
//			b.callback()
//			//b.peaks.Push(b.lastPeak)
//		}
//	}
//
//	// track the std dev
//	rollingMean := b.mean * (float64(b.len-len(args)) / float64(b.len))
//	b.mean = rollingMean + argMean*(float64(len(args))/float64(b.len))
//
//	b.stdDev = 0
//	b.mean /= float64(b.len)
//	for j := 0; j < b.len; j++ {
//		b.stdDev += math.Pow(b.buffer[j]-b.mean, 2)
//	}
//	b.stdDev = math.Sqrt(b.stdDev / float64(b.len))
//
//}
