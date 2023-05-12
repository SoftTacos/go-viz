package util

import (
	"fmt"
	"sync"
	"time"
)

func NewBeatDetector2(windowSize, samplesPerSecond, maxBpm int, callback func()) *BeatDetector2 {
	maxBeatsPerSecond := float64(maxBpm) / 60.0
	minPeriod := time.Duration(1.0 / maxBeatsPerSecond * 1000000000) // base unit of time.Duration is a nanosecond, * 1000000000 turns into milliseconds
	len := 90                                                        //int((samplesPerMinBeat + 1.0) / float64(windowSize))
	stepTime := 1.0 / float64(samplesPerSecond)
	fmt.Println("LEN", len)
	beats := 10
	return &BeatDetector2{
		peaks:      NewCircularPeakSlice(10), // TODO: how many peaks do we want?
		windowSize: windowSize,
		lastPeak: peak{
			timestamp: time.Now(),
			value:     0,
		},
		stepTimeSeconds: stepTime,
		minPeriod:       minPeriod,
		sampleRate:      float64(samplesPerSecond),
		mutex:           &sync.Mutex{},
		callback:        callback,
		beats:           make([]peak, beats),
		beatsLen:        beats,
	}
}

type BeatDetector2 struct {
	sampleRate      float64
	minPeriod       time.Duration
	peaks           *CircularPeakSlice
	windowSize      int
	lastPeak        peak
	stepTimeSeconds float64
	callback        func()
	mutex           *sync.Mutex

	beats     []peak
	beatsFront int
	beatsLen  int
}

func (b *BeatDetector2) SetCallback(callback func()) {
	b.callback = callback
}

func (b *BeatDetector2) Listen(amps chan []float64) {
	for {
		b.Push(<-amps)
	}
}

var allTimeMax float64

func (b *BeatDetector2) Push(args []float64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var (
		//slopes      = make([]float64, len(args))
		//previousVal = b.buffer[b.front]
		//argMean  float64
		maxIndex = -1
		max      float64
	)

	for i, arg := range args {
		//arg = math.Abs(arg)
		//slopes[i] = (arg - previousVal) / b.stepTimeSeconds

		if arg > max {
			max = arg
			maxIndex = i
			if max>allTimeMax{
				allTimeMax=max
			}
		}
	}

	now := time.Now()
	timeSinceLastPeak := now.Sub(b.lastPeak.timestamp)
	if timeSinceLastPeak > b.minPeriod && maxIndex > -1 {
		//var decayScale float64 = decayScaleBase - (float64(timeSinceLastPeak)-decayScaleCoeff*float64(b.minPeriod))/1000000000
		threshold:=allTimeMax*.9
		if max > threshold {
			fmt.Printf("BEAT DETECTED: %+v last:%.2f t:%.2f\n", timeSinceLastPeak, b.lastPeak.value, threshold)
			b.lastPeak = peak{
				timestamp: now,
				value:     max,
			}
			b.beats[b.beatsFront] = b.lastPeak
			b.beatsFront = (b.beatsFront + 1) % b.beatsLen
			b.callback()
		}
	}
}
