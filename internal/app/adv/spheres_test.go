package adv

import (
	"simd"
	"slices"
	"testing"

	"github.com/eriklupander/simd-tracer/internal/app/std"
	"github.com/stretchr/testify/assert"
)

func TestRotate(t *testing.T) {
	vec1 := simd.LoadFloat32x8Slice([]float32{3, 2, 1, 0, 3, 2, 5, 0})
	vec2 := simd.LoadFloat32x8Slice([]float32{1, 1, 1, 1, 1, 1, 1, 1})
	res := vec1.Equal(vec2).AsInt32x8()
	s := make([]int32, 8)
	res.StoreSlice(s)
	t.Log(s)

}

func TestFindLowest(t *testing.T) {
	minT := simd.LoadFloat32x8Slice([]float32{3, 2, 1, 3, 3, 7, 12, 3})
	hi := minT.GetHi()
	i := hi.Min(minT.GetLo())

	i = i.Min(i.SelectFromPair(1, 2, 3, 0, i))
	i = i.Min(i.SelectFromPair(2, 3, 0, 1, i))
	i = i.Min(i.SelectFromPair(3, 0, 1, 2, i))

	t.Log(i)

}

func TestIntersectSpheres(t *testing.T) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}

	spheres := asSpheres(std.SixteenSpheres())
	t0, hitIndex, hit := IntersectSpheres(ray, spheres)
	assert.True(t, hit)
	assert.Equal(t, 5, hitIndex)
	assert.InEpsilon(t, 17, t0, 0.1)
}

func TestIntersectSpheresSIMD(t *testing.T) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}

	spheres := asSpheres(std.SixteenSpheres())
	t0, hitIndex, hit := IntersectSpheresSIMD(ray, spheres)

	assert.True(t, hit)
	assert.Equal(t, 5, hitIndex)
	assert.InEpsilon(t, 17, t0, 0.1)
}

func TestIntersectSpheresSIMDExt(t *testing.T) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}

	spheres := asSpheres(std.SixteenSpheres())
	results := make([]float32, spheres.Count)
	IntersectSpheresSIMDExt(ray, spheres, results)
	t.Log(results)
}

func BenchmarkIntersectSpheres(b *testing.B) {
	ray := &std.Ray{
		Orig: &std.Vec4{3.1, 4.1, 20.1, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}
	spheres := asSpheres(StandardSpheres())
	for b.Loop() {
		_, _, _ = IntersectSpheres(ray, spheres)
	}
}

func Benchmark16IntersectSpheres(b *testing.B) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}
	spheres := asSpheres(std.SixteenSpheres())
	for b.Loop() {
		_, _, _ = IntersectSpheres(ray, spheres)
	}
}

func Benchmark16IntersectSpheresSIMD(b *testing.B) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}
	spheres := asSpheres(std.SixteenSpheres())
	for b.Loop() {
		_, _, _ = IntersectSpheresSIMD(ray, spheres)
	}
}

func Benchmark16IntersectSpheresSIMDExt(b *testing.B) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}
	spheres := asSpheres(std.SixteenSpheres())
	results := make([]float32, spheres.Count)
	for b.Loop() {
		IntersectSpheresSIMDExt(ray, spheres, results)
		_ = slices.Min(results)
	}

}

func asSpheres(spheres []std.Sphere) *Spheres {
	sp := &Spheres{
		CenterX:       make([]float32, len(spheres)),
		CenterY:       make([]float32, len(spheres)),
		CenterZ:       make([]float32, len(spheres)),
		RadiusSquared: make([]float32, len(spheres)),
		Count:         len(spheres),
	}
	for i, s := range spheres {
		sp.CenterX[i] = s.Center[0]
		sp.CenterY[i] = s.Center[1]
		sp.CenterZ[i] = s.Center[2]
		sp.RadiusSquared[i] = s.RadiusSquared
	}
	return sp
}

func StandardSpheres() []std.Sphere {
	spheres := make([]std.Sphere, 0)

	// Origin
	spheres = append(spheres, std.Sphere{
		Center:        &std.Vec4{0, 0, 0, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         std.Vec4{1, 0, 0},
	})

	// 3 unit to the left (green)
	spheres = append(spheres, std.Sphere{
		Center:        &std.Vec4{-3, 0, 0, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         std.Vec4{0, 1, 0},
	})

	// 3 unit to the right (blue)
	spheres = append(spheres, std.Sphere{
		Center:        &std.Vec4{3, 0, 0, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         std.Vec4{0, 0, 1},
	})

	// 3 units up
	spheres = append(spheres, std.Sphere{
		Center:        &std.Vec4{0, 3, 0, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         std.Vec4{0, 1, 1},
	})

	// Closer to the camera
	spheres = append(spheres, std.Sphere{
		Center:        &std.Vec4{0, 0, 3, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         std.Vec4{1, 1, 0},
	})

	// Even closer
	spheres = append(spheres, std.Sphere{
		Center:        &std.Vec4{0, 0, 6, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         std.Vec4{0.8, 0.8, 0.8},
	})

	return spheres
}
