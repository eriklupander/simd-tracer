package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"runtime/pprof"
	"time"

	"github.com/eriklupander/simd-tracer/internal/app/adv"
	"github.com/eriklupander/simd-tracer/internal/app/std"
)

// Change to get a long-running main suitable for profiling.
const iterations = 1

// If changed to true, CPU profiling is enabled.
const enablePprof = false

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	if enablePprof {
		cleanFn := enableProfiling()
		defer cleanFn()
	}

	spheres := std.SixteenSpheres()
	planes := std.CornellBox()

	outBuf := new(bytes.Buffer)

	st := time.Now()
	for range 1 {
		std.Render(screenWidth, screenHeight, spheres, planes, outBuf)
	}
	fmt.Printf("Legacy done in %v\n", time.Since(st))

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

	simdSpheres := adv.AsStructOfArrays(spheres)
	simdPlanes := adv.PlanesAsStructOfArrays(planes)
	//simdTriangles := adv.TrianglesSmallSideBySide()
	outBuf2 := new(bytes.Buffer)
	st2 := time.Now()
	for range 1 {
		adv.Render(screenWidth, screenHeight, simdSpheres, simdPlanes, &adv.Triangles{}, outBuf2)
	}
	fmt.Printf("SIMD done in %v", time.Since(st2))

	img2 := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	img2.Pix = outBuf2.Bytes()
	f2, err := os.OpenFile("out2.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f2, img2); err != nil {
		panic(err.Error())
	}
	_ = f2.Close()
}

func enableProfiling() func() {
	f2, err := os.Create("simd-tracer.pprof")
	if err != nil {
		panic(err.Error())
	}
	pprof.StartCPUProfile(f2)
	return pprof.StopCPUProfile
}
