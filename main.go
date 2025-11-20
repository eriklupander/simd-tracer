package main

import (
	"fmt"
	"image"
	"image/png"
	rnd "math/rand/v2"
	"os"
	"runtime/pprof"
	"time"

	"github.com/eriklupander/simd-tracer/internal/app/std"
	//
	//"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 3200
	screenHeight = 2400
)

//
//type rand struct {
//	x, y, z, w uint32
//}
//
//func (r *rand) next() uint32 {
//	// math/rand is too slow to keep 60 FPS on web browsers.
//	// Use Xorshift instead: http://en.wikipedia.org/wiki/Xorshift
//	t := r.x ^ (r.x << 11)
//	r.x, r.y, r.z = r.y, r.z, r.w
//	r.w = (r.w ^ (r.w >> 19)) ^ (t ^ (t >> 8))
//	return r.w
//}
//
//var theRand = &rand{12345678, 4185243, 776511, 45411}

type Game struct {
	noiseImage *image.RGBA
}

//
//func (g *Game) Update() error {
//	return nil
//}
//
//func (g *Game) Draw(screen *ebiten.Image) {
//	screen.WritePixels(g.noiseImage.Pix)
//	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
//}

//func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
//	return screenWidth, screenHeight
//}

func main() {
	f2, err := os.Create("simd-legacy-5.pprof")
	if err != nil {
		panic(err.Error())
	}
	pprof.StartCPUProfile(f2)
	defer pprof.StopCPUProfile()

	//ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	//ebiten.SetWindowTitle("Noise (Ebitengine Demo)")
	//g := &Game{
	//	noiseImage: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	//}

	// generate a scene made of random spheres
	//spheres := generateSpheres()

	//sphereData, err := os.ReadFile("scene1.json")
	//if err != nil {
	//	panic(err.Error())
	//}
	//spheres := make([]std.Sphere, 0)
	//if err := json.Unmarshal(sphereData, &spheres); err != nil {
	//	panic(err.Error())
	//}

	//	go func() {
	// Generate the noise with random RGB values.
	img := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))

	spheres := std.StandardSpheres()
	planes := std.CornellBox()

	st := time.Now()
	for range 10 {
		std.Render(screenWidth, screenHeight, spheres, planes, img.Pix)
	}
	fmt.Printf("Done in %v", time.Since(st))

	f, err := os.OpenFile("out1.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f, img); err != nil {
		panic(err.Error())
	}
	_ = f.Close()
	//	}()

	//if err := ebiten.RunGame(g); err != nil {
	//	log.Fatal(err)
	//}
}

func generateSpheres() []std.Sphere {
	spheres := make([]std.Sphere, 0)
	for range 32 {
		randPos := &std.Vec4{(0.5 - rnd.Float32()) * 10, (0.5 - rnd.Float32()) * 10, (0.5 + rnd.Float32()) * 10}
		//Vec4f randPos((, (, ());
		randRadius := (0.5 + rnd.Float32()) * 0.5
		spheres = append(spheres, std.Sphere{
			Center: randPos,
			Radius: randRadius,
			Color:  std.Vec4{rnd.Float32(), rnd.Float32(), rnd.Float32()},
		})
	}
	return spheres
}
