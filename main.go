package main

import (
	pa "github.com/gordonklaus/portaudio"
	"github.com/hajimehoshi/ebiten/v2"
	m "github.com/softtacos/go-visualizer/manager"
	v "github.com/softtacos/go-visualizer/manager/visualizers"
	s "github.com/softtacos/go-visualizer/stream"
	"log"
	"os"
	"os/signal"
)

const (
	wWidth,wHeight = 1280,960
)

func main() {
	var err error
	if err=pa.Initialize();err!=nil{
		log.Fatal("failed to initialize portaudio ",err)
		return
	}
	defer pa.Terminate()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	var buffer = 2
	streamOutput :=make(chan []float32,buffer)

	streamer:=s.NewStreamer(streamOutput,buffer,sig)
	streamer.Setup()
	streamer.Start2()

	ebiten.SetWindowSize(wWidth,wHeight)
	ebiten.SetWindowResizable(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	ebiten.SetWindowTitle("go visualizer")
	//ebiten.SetFullscreen(true)
	//ebiten.SetWindowSize(1791,1120)
	if err = ebiten.RunGame(m.NewVisManager(streamOutput,[]m.VisualizerConstructor{
		v.NewBasicVisualizer,
		//v.NewLazyPaddedVisualizer,
		//v.NewLazyCircleVisualizer,
	}));err!=nil{
		log.Fatal("error running game ",err)
	}
}

