package main

import (
	"fmt"
	"time"

	"github.com/eriklupander/simd-tracer/internal/app/std"
)

// main can be used to run one of the ray-sphere intersection code-paths in a loop in order to either collect profiling
// data or possibly attach to with an external profiler such as Intruments.
func main() {

	//f2, err := os.Create("simd.pprof")
	//if err != nil {
	//	panic(err.Error())
	//}
	//pprof.StartCPUProfile(f2)
	//defer pprof.StopCPUProfile()

	spheres := std.StandardSpheres()
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.15241566, -0.00019870223, -0.9883165, 0},
	}
	hits := 0
	st := time.Now()
	var hit = false
	for range 100_000_000 {
		for _, s := range spheres {
			_, hit = std.IntersectSphereSIMD(ray, s.Center, s.RadiusSquared)
			if hit {
				hits++
			}
		}
	}
	fmt.Printf("hit: %d times in %v\n", hits, time.Since(st))

	//for range 100_000_000 {
	//	for _, s := range spheres {
	//		_, hit := std.IntersectSphere(ray, s.Center, s.RadiusSquared)
	//		if hit {
	//			hits++
	//		}
	//	}
	//}
	//fmt.Printf("hit: %d times in %v\n", hits, time.Since(st))
	//
	//hits = 0
	//st = time.Now()

	//hits = 0
	//st = time.Now()
	//for range 100_000_000 {
	//	for _, s := range spheres {
	//		_, hit := std.IntersectSphereFullSIMD(ray, s.Center, s.RadiusSquared)
	//		if hit {
	//			hits++
	//		}
	//	}
	//}
	//fmt.Printf("hit: %d times in %v\n", hits, time.Since(st))
}
