package util

func NewCircularAudioBuffer(len int)*CircularBuffer{
	return &CircularBuffer{
		len:len,
	}
}

type CircularBuffer struct {
	len int
	front,back int
	data []float64
}

func (c *CircularBuffer)GetPeaks(sep int){}