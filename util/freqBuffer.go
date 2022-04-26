package util

import (
	"sync"
	"time"
)

func NewFrequencyBuffer(len, samples int) *Buffer {
	frequencies:= make([][]float64, len)
	for i:=range frequencies{
		frequencies[i] = make([]float64, samples)
	}
	return &Buffer{
		len:         len,
		samples:     samples,
		mutex:       &sync.Mutex{},
		frequencies: frequencies,
		//frequencies: make([][]float64, len),
		//derivative: make([]float64,len),
		lastSample: time.Now(),
	}
}

type Buffer struct {
	front        int
	len, samples int
	frequencies  [][]float64 // [i] is the timeslice, [j] is the frequency band
	lastSample   time.Time
	derivative   []float64
	mutex        *sync.Mutex
}

func (q *Buffer) Push(args ...[]float64) {
	q.mutex.Lock()
	for i := range args {
		q.frequencies[(q.front+i)%q.len] = args[i]
	}
	q.front = (q.front + len(args)) % q.len

	f := q.front - 1
	if f < 0 {
		f = q.len - 1
	}
	now := time.Now()
	q.derivative = make([]float64, len(args[0]))
	for i := range q.frequencies[0] {
		q.derivative[i] = (q.frequencies[q.front][i] - q.frequencies[f][i]) / float64(q.lastSample.Sub(now))
	}
	q.lastSample = now
	q.mutex.Unlock()
}

func (q *Buffer) Get() [][]float64 {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.frequencies
}

func (q *Buffer) GetDerivative() []float64 {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.derivative
}

func (q *Buffer) GetAverage() (avg []float64) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	l := float64(len(q.frequencies))
	avg = make([]float64, len(q.frequencies[0]))
	//for i:=range q.frequencies{
	//	for j:=range q.frequencies[i]{
	//		avg[j]+=q.frequencies[i][j]
	//	}
	//}
	for j := range q.frequencies[0] {
		for i := range q.frequencies {
			avg[j] += q.frequencies[i][j]
		}
		avg[j] /= l
	}
	return
}
