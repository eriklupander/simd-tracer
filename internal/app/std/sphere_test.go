package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersectSphere(t *testing.T) {

	ray := &Ray{
		Orig: &Vec4{3, 4, 20, 0},
		Dir:  &Vec4{-0.15241566, -0.00019870223, -0.9883165, 0},
	}

	dist, hit := IntersectSphere(ray, &Vec4{0, 3, 0, 0}, 1.0)
	assert.True(t, hit)
	assert.InEpsilon(t, 20.0, float64(dist), 0.01)
}

func TestIntersectSphereFullSIMD(t *testing.T) {

	ray := &Ray{
		Orig: &Vec4{3, 4, 20, 0},
		Dir:  &Vec4{-0.15241566, -0.00019870223, -0.9883165, 0},
	}

	dist, hit := IntersectSphereFullSIMD(ray, &Vec4{0, 3, 0, 0}, 1.0)
	assert.True(t, hit)
	assert.InEpsilon(t, 20.0, float64(dist), 0.01)
}

func TestIntersectSpheres(t *testing.T) {
	spheres := StandardSpheres()
	ray := &Ray{
		Orig: &Vec4{3, 4, 20, 0},
		Dir:  &Vec4{-0.29983652 - 0.15058117 - 0.9420315, 0},
	}
	for i, s := range spheres {
		dist, hit := IntersectSphere(ray, s.Center, s.RadiusSquared)
		if hit {
			t.Logf("hit %d at %f", i, dist)
		}
	}
}

func BenchmarkIntersectSpheres(b *testing.B) {
	spheres := StandardSpheres()
	ray := &Ray{
		Orig: &Vec4{3.1, 4.1, 20.1, 0},
		Dir:  &Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}
	for b.Loop() {
		for _, s := range spheres {
			_, _ = IntersectSphere(ray, s.Center, s.RadiusSquared)
		}
	}
}

func BenchmarkIntersectSpheresSIMD(b *testing.B) {
	spheres := StandardSpheres()
	ray := &Ray{
		Orig: &Vec4{3.1, 4.1, 20.1, 0},
		Dir:  &Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}
	for b.Loop() {
		for _, s := range spheres {
			_, _ = IntersectSphereSIMD(ray, s.Center, s.RadiusSquared)
		}
	}
}

func BenchmarkIntersectSingleSphere(b *testing.B) {
	s := Sphere{
		Center:        &Vec4{0, 3, 0, 0},
		RadiusSquared: 1,
		Color:         Vec4{0, 1, 1},
	}
	ray := &Ray{
		Orig: &Vec4{3.1, 4.1, 20.1, 0},
		Dir:  &Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}
	for b.Loop() {
		_, _ = IntersectSphere(ray, s.Center, s.RadiusSquared)
	}
}

func BenchmarkIntersectSingleSpheresSIMD(b *testing.B) {

	ray := &Ray{
		Orig: &Vec4{3.1, 4.1, 20.1, 0},
		Dir:  &Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}

	center := &Vec4{0, 3, 0, 0}
	b.Logf("%p", ray.Orig)
	b.Logf("%p", ray.Dir)
	b.Logf("%p", center)
	for b.Loop() {
		_, _ = IntersectSphereSIMD(ray, center, 1)
	}
}

func BenchmarkIntersectSpheresSIMDMethod(b *testing.B) {
	spheres := StandardSpheres()
	ray := &Ray{
		Orig: &Vec4{3.1, 4.1, 20.1, 0},
		Dir:  &Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}
	for b.Loop() {
		for _, s := range spheres {
			_, _ = IntersectSphereSIMDMethod(ray, s.Center, s.RadiusSquared)
		}
	}
}

func BenchmarkIntersectSpheresFullSIMD(b *testing.B) {
	spheres := StandardSpheres()
	ray := &Ray{
		Orig: &Vec4{3.1, 4.1, 20.1, 0},
		Dir:  &Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}
	for b.Loop() {
		for _, s := range spheres {
			_, _ = IntersectSphereFullSIMD(ray, s.Center, s.RadiusSquared)
		}
	}
}

func BenchmarkIntersectSpheresDotGoSIMD(b *testing.B) {
	spheres := StandardSpheres()
	ray := &Ray{
		Orig: &Vec4{3.1, 4.1, 20.1, 0},
		Dir:  &Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}
	for b.Loop() {
		for _, s := range spheres {
			_, _ = IntersectSphereDotGoSIMD(ray, s.Center, s.RadiusSquared)
		}
	}
}
