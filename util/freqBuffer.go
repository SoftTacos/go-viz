package util

import "sync"

func NewFrequencyBuffer(len int)*FrequenciesBuffer {
	return &FrequenciesBuffer{
		len:len,
		mutex: &sync.Mutex{},
		stuff:make([][]float64,len),
	}
}

type FrequenciesBuffer struct {
	front int
	len int
	stuff [][]float64
	mutex *sync.Mutex
}

func (q *FrequenciesBuffer)Push(args ...[]float64){
	//for i:=range args{
	//	q.stuff[(q.front+i)%q.len] = args[i]
	//}
	//q.front = (q.front+len(args))%q.len
	q.mutex.Lock()
	for i :=range args{
		q.stuff[(q.front+i)%q.len] = args[i]
	}
	q.front = (q.front+len(args))%q.len
	q.mutex.Unlock()
}

func (q *FrequenciesBuffer)Get()[][]float64{
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.stuff
}

func (q *FrequenciesBuffer)GetAverage()(avg []float64){
	q.mutex.Lock()
	defer q.mutex.Unlock()
	l:=float64(len(q.stuff))
	avg = make([]float64,len(q.stuff[0]))
	//for i:=range q.stuff{
	//	for j:=range q.stuff[i]{
	//		avg[j]+=q.stuff[i][j]
	//	}
	//}
	for j:=range q.stuff[0]{
		for i:=range q.stuff{
			avg[j]+=q.stuff[i][j]
		}
		avg[j]/=l
	}
	return
}
