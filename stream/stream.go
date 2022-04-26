package stream

import (
	"fmt"
	pa "github.com/gordonklaus/portaudio"
	"log"
	"os"
	"strings"
)

func NewStreamer(windowSize int, output chan []float64,cap int,stop chan os.Signal)*streamer{
	return &streamer{
		windowSize:windowSize,
		cap:cap,
		output:output,
		stop:stop,
	}
}

type streamer struct {
	cap,windowSize int
	stream *pa.Stream
	output chan []float64
	stop chan os.Signal
}

func (s *streamer) Start() {

	go func() {
		in := make([]int32, s.windowSize)
		stream,err:=pa.OpenDefaultStream(1,0,48000,len(in),in) // 44100
		if err!=nil{
			log.Panic("failed to open stream:",err)
		}
		defer stream.Close()

		if err=stream.Start();err!=nil{
			log.Panic("failed to start stream:",err)
		}
		for {
			//time.Sleep(time.Microsecond*500)
			if err = stream.Read(); err != nil {
				log.Panic("failed to read from stream ", err)
			}
			select {
			case <-s.stop:
				log.Println("shutting down listener")
				return
			default:
			}

			if len(s.output) == s.cap {
				<-s.output
			}
			o := make([]float64,len(in))
			for i:=range in{
				o[i] =float64(in[i])
			}
			s.output<-o

		}
	}()
}

func (s *streamer)Setup() {
	var device *pa.DeviceInfo
	devices,err:=pa.Devices()
	if err!=nil{
		panic(err)
	}
	for _,d:=range devices{
		//USB Audio Device
		if strings.Contains(d.Name,"USB") && d.MaxInputChannels > 0 {
			device = d
		 break
		}
		//if d.MaxInputChannels > 0 {
		//	device = d
		//	break
		//}
	}
	fmt.Printf("device:%+v\n",device)
	if device == nil{
		log.Panic("no USB device found!")
	}
	stream,err:=pa.OpenStream(pa.StreamParameters{
		Input: pa.StreamDeviceParameters{
			Device:device,
			Channels: 1,
		},
		SampleRate: device.DefaultSampleRate,
		FramesPerBuffer: s.windowSize,
	},s.readCallback)
	if err!=nil{
		log.Panic("failed to open stream:",err)
	}
	s.stream = stream
}

func (s *streamer)Start2(){
	if err:=s.stream.Start();err!=nil{
		log.Panic("failed to start stream",err)
	}
}
func (s *streamer) readCallback(in []float32){
	//fmt.Println(in)
	o := make([]float64,len(in))
	for i:=range in{
		o[i] =float64(in[i])
	}
	s.output<-o
}

func (s *streamer) Stop() {
	s.stop <- os.Interrupt
}