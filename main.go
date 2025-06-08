package main

import (
    "github.com/gordonklaus/portaudio"
    "log"
)

func main() {
    portaudio.Initialize()
    defer portaudio.Terminate()

    const sampleRate = 44100
    const framesPerBuffer = 64

    in := make([]float32, framesPerBuffer)
	// out := make([]float32, framesPerBuffer)

    stream, err := portaudio.OpenDefaultStream(1, 1, sampleRate, len(in), func(inBuf, outBuf []float32) {
        copy(outBuf, inBuf)
    })
    if err != nil {
        log.Fatal(err)
    }
    defer stream.Close()

    if err := stream.Start(); err != nil {
        log.Fatal(err)
    }

    log.Println("Streaming mic input to output. Press Ctrl+C to stop.")
    select {} // block forever
}
