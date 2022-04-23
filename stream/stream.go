package stream

import (
	pa "github.com/gordonklaus/portaudio"
	"log"
	"os"
	"time"
)

func NewStreamer(output chan []float64,cap int,stop chan os.Signal)*streamer{
	return &streamer{
		cap:cap,
		output:output,
		stop:stop,
	}
}

type streamer struct {
	cap int
	output chan []float64
	stop chan os.Signal
}

func (s *streamer) Start() {
	const numSamples = 128
	go func() {
		in := make([]int32, numSamples)
		stream,err:=pa.OpenDefaultStream(1,0,44100,len(in),in)
		if err!=nil{
			log.Panic("failed to open stream:",err)
		}
		defer stream.Close()

		if err=stream.Start();err!=nil{
			log.Panic("failed to start stream:",err)
		}
		for {
			time.Sleep(time.Microsecond*500)
			if err = stream.Read(); err != nil {
				log.Panic("failed to read from stream ", err)
			}
			select {
			case <-s.stop:
				return
			default:
			}
			var out =make([]float64,numSamples)
			for i:=range in{
				out[i] = float64(in[i])
			}

			if len(s.output) ==s.cap {
				<-s.output
			}
			s.output <- out
		}
	}()
}

func (s *streamer) Stop() {
	s.stop <- os.Interrupt
}