package util

import "math"

// breaks f into n groups/buckets, each group is the average
// value of all values of f in that group/bucket
//func GroupFrequencies(n int, freqs []float64)(groups []float64){
//	groups =make([]float64,n)
//	chunkSize := len(freqs)/n
//	for i:=0;i<n;i++{
//		for _,f:=range freqs[i*chunkSize:(i+1)*chunkSize]{
//			groups[i]+=f
//		}
//		groups[i]/=float64(chunkSize)
//	}
//
//	return
//}

func GroupFrequencies(n int, freqs []float64)(groups []float64){
	if len(freqs) == 0 || n > len(freqs){
		return
	}
	groups =make([]float64,n)
	chunkSize := float64(len(freqs))/float64(n)
	for i:=0;i<n;i++{
		start,end:=math.Round(float64(i)*chunkSize),math.Round(float64(i+1)*chunkSize)
		for _,f:=range freqs[int(start):int(end)]{
			groups[i]+=f
		}
		groups[i]/=end-start
	}

	return
}
