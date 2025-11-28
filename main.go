package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"runtime/pprof"
	"time"

	"github.com/eriklupander/simd-tracer/internal/app/std"
)

// Change to get a long-running main suitable for profiling.
const iterations = 1

// If changed to true, CPU profiling is enabled.
const enablePprof = false

const (
	screenWidth  = 3200
	screenHeight = 2400
)

func main() {
	if enablePprof {
		cleanFn := enableProfiling()
		defer cleanFn()
	}

	spheres := std.SixteenSpheres()
	planes := std.CornellBox()

	st := time.Now()
	outBuf := new(bytes.Buffer)
	for range 1 {
		std.Render(screenWidth, screenHeight, spheres, planes, outBuf)
	}
	fmt.Printf("Done in %v", time.Since(st))

	img := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	img.Pix = outBuf.Bytes()
	f, err := os.OpenFile("out1.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f, img); err != nil {
		panic(err.Error())
	}
	_ = f.Close()
}

func enableProfiling() func() {
	f2, err := os.Create("simd-tracer.pprof")
	if err != nil {
		panic(err.Error())
	}
	pprof.StartCPUProfile(f2)
	return pprof.StopCPUProfile
}
