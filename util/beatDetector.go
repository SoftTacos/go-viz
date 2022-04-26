package util

import (
	"time"
)

func NewBeatDetector(len, samples int, minPeriod time.Duration) *BeatDetector {
	return &BeatDetector{
		lastBeat:  time.Now(),
		minPeriod: minPeriod,
		Buffer:    NewFrequencyBuffer(len,samples),
	}
}

type BeatDetector struct {
	*Buffer
	sampleRate float64
	lastBeat   time.Time
	minPeriod  time.Duration
}

func (b *BeatDetector)Push(args ...[]float64){
	b.Buffer.Push(args...)
	//
}