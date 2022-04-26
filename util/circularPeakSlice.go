package util

import "sync"


// TODO: can make this generic, nervous about the performance hit though
func NewCircularPeakSlice(len int)*CircularPeakSlice {
	values := make([]peak,len)
	return &CircularPeakSlice{
		len: len,
		Values: values,
		mutex: &sync.Mutex{},
	}
}

type CircularPeakSlice struct {
	len, front int
	Values     []peak
	mutex      *sync.Mutex
}

func (c *CircularPeakSlice) GetFront() int {
	return c.front
}

func (c *CircularPeakSlice) Push(args ...peak) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i := range args {
		c.Values[(c.front+i)%c.len] = args[i]
	}
	c.front = (c.front + len(args)) % c.len
}

func (c *CircularPeakSlice) Get(start,end int)[]peak{
	return nil // TODO
}