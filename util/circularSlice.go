package util

import "sync"

// semi-encapsulated for reusability

func NewCircularSlice(len,front int)*CircularSlice{
	values := make([]float64,len)
	return &CircularSlice{
		len: len,
		Values: values,
		mutex: &sync.Mutex{},
	}
}

type CircularSlice struct {
	len, front int
	Values     []float64
	mutex      *sync.Mutex
}

func (c *CircularSlice) GetFront() int {
	return c.front
}

func (c *CircularSlice) Push(args ...float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i := range args {
		c.Values[(c.front+i)%c.len] = args[i]
	}
	c.front = (c.front + len(args)) % c.len
}
