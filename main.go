package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log/slog"
	"os"
	"runtime/pprof"
	"time"

	"github.com/eriklupander/simd-tracer/internal/app/adv"
	"github.com/eriklupander/simd-tracer/internal/app/std"
)

// If changed to true, CPU profiling is enabled.
const enablePprof = false

func main() {
	var width, height, iterations int
	var outFile string
	flag.IntVar(&width, "width", 1024, "output image width in pixels")
	flag.IntVar(&height, "height", 768, "output image height in pixels")
	flag.IntVar(&iterations, "iterations", 1, "Number of times to render image, useful when running the profiler over extended periods of time")
	flag.StringVar(&outFile, "out", "out2.png", "output file name")

	var lightSamples int
	flag.IntVar(&lightSamples, "light-samples", 1, "number of random samples per light source for soft shadow computation")
	flag.Parse()

	if enablePprof {
		cleanFn := enableProfiling()
		defer cleanFn()
	}

	// Build scene objects
	spheres := std.SixteenSpheres()
	planes := std.CornellBox()

	// Transform plain slices to Struct of Arrays format.
	simdSpheres := adv.AsStructOfArrays(spheres)
	simdPlanes := adv.PlanesAsStructOfArrays(planes)
	//simdTriangles := adv.TrianglesSmallSideBySide()

	outBuf := new(bytes.Buffer)
	st := time.Now()
	lights := std.Lights()
	for range iterations {
		adv.Render(width, height, simdSpheres, simdPlanes, &adv.Triangles{}, lights, lightSamples, outBuf)
	}
	fmt.Printf("SIMD done in %v", time.Since(st))

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Pix = outBuf.Bytes()
	file, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		slog.Error("Error opening output file", slog.Any("error", err))
		os.Exit(1)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		slog.Error("Error encoding image buffer to PNG", slog.Any("error", err))
		os.Exit(1)
	}
}

func renderLegacy(width, height int, spheres []std.Sphere, planes []std.Plane) {

	outBuf := new(bytes.Buffer)

	st := time.Now()
	for range 1 {
		std.Render(width, height, spheres, planes, outBuf)
	}
	fmt.Printf("Legacy done in %v\n", time.Since(st))

	img := image.NewRGBA(image.Rect(0, 0, width, height))
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
	_ = pprof.StartCPUProfile(f2)
	return pprof.StopCPUProfile
}
