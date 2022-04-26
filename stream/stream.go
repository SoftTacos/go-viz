package stream

import (
	"fmt"
	pa "github.com/gordonklaus/portaudio"
	"log"
	"strings"
)

func NewStreamer(windowSize int, output chan []float64, cap int) *Streamer {
	s := &Streamer{
		windowSize: windowSize,
		cap:        cap,
		output:     output,
	}
	s.setup()
	return s
}

type Streamer struct {
	cap, windowSize int
	device          *pa.DeviceInfo
	stream          *pa.Stream
	output          chan []float64
}

func (s *Streamer) setup() {
	devices, err := pa.Devices()
	if err != nil {
		panic(err)
	}
	for _, d := range devices {
		//USB Audio Device
		if strings.Contains(d.Name, "USB") && d.MaxInputChannels > 0 {
			s.device = d
			break
		}
		//if d.MaxInputChannels > 0 {
		//	device = d
		//	break
		//}
	}
	fmt.Printf("device:%+v\n", s.device)
	if s.device == nil {
		log.Panic("no USB device found!")
	}
	stream, err := pa.OpenStream(pa.StreamParameters{
		Input: pa.StreamDeviceParameters{
			Device:   s.device,
			Channels: 1,
		},
		SampleRate:      s.device.DefaultSampleRate,
		FramesPerBuffer: s.windowSize,
	}, s.readCallback)
	if err != nil {
		log.Panic("failed to open stream:", err)
	}
	s.stream = stream
	fmt.Println("Streamer setup complete")
}

func (s *Streamer) Start() {
	if err := s.stream.Start(); err != nil {
		log.Panic("failed to start stream", err)
	}
	fmt.Println("Streamer started")
}
func (s *Streamer) readCallback(in []float32) {
	// TODO: add overflow <-s.output if len is full
	o := make([]float64, len(in))
	for i := range in {
		o[i] = float64(in[i])
	}
	s.output <- o
}

//func (s *Streamer) Start_old() {
//
//	go func() {
//		in := make([]int32, s.windowSize)
//		stream, err := pa.OpenDefaultStream(1, 0, 48000, len(in), in) // 44100
//		if err != nil {
//			log.Panic("failed to open stream:", err)
//		}
//		defer stream.Close()
//
//		if err = stream.Start(); err != nil {
//			log.Panic("failed to start stream:", err)
//		}
//		for {
//			//time.Sleep(time.Microsecond*500)
//			if err = stream.Read(); err != nil {
//				log.Panic("failed to read from stream ", err)
//			}
//			select {
//			case <-s.stop:
//				log.Println("shutting down listener")
//				return
//			default:
//			}
//
//			if len(s.output) == s.cap {
//				<-s.output
//			}
//			o := make([]float64, len(in))
//			for i := range in {
//				o[i] = float64(in[i])
//			}
//			s.output <- o
//
//		}
//	}()
//}
