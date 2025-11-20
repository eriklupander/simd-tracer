package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersectSphere(t *testing.T) {

	ray := &Ray{
		orig: &Vec4{3, 4, 20, 0},
		dir:  &Vec4{-0.15241566, -0.00019870223, -0.9883165, 0},
	}

	dist, hit := intersectSphere(ray, &Vec4{0, 3, 0, 0}, 1.0)
	assert.True(t, hit)
	assert.InEpsilon(t, 20.0, float64(dist), 0.01)
}

func TestIntersectSphereFullSIMD(t *testing.T) {

	ray := &Ray{
		orig: &Vec4{3, 4, 20, 0},
		dir:  &Vec4{-0.15241566, -0.00019870223, -0.9883165, 0},
	}

	dist, hit := intersectSphereFullSIMD(ray, &Vec4{0, 3, 0, 0}, 1.0)
	assert.True(t, hit)
	assert.InEpsilon(t, 20.0, float64(dist), 0.01)
}

// 0 = {float32} -0.15241566
// 1 = {float32} -0.00019870223
// 2 = {float32} -0.9883165
// &[-0.14217138 -0.042353157 -0.98893553 0]
func TestIntersectSpheres(t *testing.T) {
	spheres := StandardSpheres()
	ray := &Ray{
		orig: &Vec4{3, 4, 20, 0},
		dir:  &Vec4{-0.29983652 - 0.15058117 - 0.9420315, 0},
	}
	for i, s := range spheres {
		dist, hit := intersectSphere(ray, s.Center, s.RadiusSquared)
		if hit {
			t.Logf("hit %d at %f", i, dist)
		}
	}
}

func BenchmarkIntersectSpheres(b *testing.B) {
	spheres := StandardSpheres()
	ray := &Ray{
		orig: &Vec4{3.1, 4.1, 20.1, 0},
		dir:  &Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}
	for b.Loop() {
		for _, s := range spheres {
			_, _ = intersectSphere(ray, s.Center, s.RadiusSquared)
		}
	}
}

func BenchmarkIntersectSpheresSIMD(b *testing.B) {
	spheres := StandardSpheres()
	ray := &Ray{
		orig: &Vec4{3.1, 4.1, 20.1, 0},
		dir:  &Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}
	for b.Loop() {
		for _, s := range spheres {
			_, _ = intersectSphereSIMD(ray, s.Center, s.RadiusSquared)
		}
	}
}
func BenchmarkIntersectSpheresFullSIMD(b *testing.B) {
	spheres := StandardSpheres()
	ray := &Ray{
		orig: &Vec4{3.1, 4.1, 20.1, 0},
		dir:  &Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}
	for b.Loop() {
		for _, s := range spheres {
			_, _ = intersectSphereFullSIMD(ray, s.Center, s.RadiusSquared)
		}
	}
}

//
//func BenchmarkIntersectSpheresIntegersOnly(b *testing.B) {
//	spheres := StandardSpheres()
//	ray := &Ray{
//		orig: &Vec4{3.1, 4.1, 20.1, 0},
//		dir:  &Vec4{-1, 0, -1, 0},
//	}
//	for b.Loop() {
//		for _, s := range spheres {
//			_, _ = intersectSphere(ray, s.Center, s.RadiusSquared)
//		}
//	}
//}
//
//func BenchmarkIntersectSpheresSIMDIntegersOnly(b *testing.B) {
//	spheres := StandardSpheres()
//	ray := &Ray{
//		orig: &Vec4{3.1, 4.1, 20.1, 0},
//		dir:  &Vec4{-1, 0, -1, 0},
//	}
//	for b.Loop() {
//		for _, s := range spheres {
//			_, _ = intersectSphereSIMD(ray, s.Center, s.RadiusSquared)
//		}
//	}
//}
