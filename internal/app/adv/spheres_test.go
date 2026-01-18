package adv

import (
	"simd/archsimd"
	"testing"

	"github.com/eriklupander/simd-tracer/internal/app/std"
	"github.com/stretchr/testify/assert"
)

func TestScalarBitMasks(t *testing.T) {
	currentMask := 0b01000000

	t.Logf("%08b", currentMask)
	t.Logf("%v", currentMask|0b0000 == 0)
	t.Logf("%v", currentMask&0b1111 != 0)
	t.Logf("%v", currentMask^0b0000 == 0)

	currentMask = 0b00000010
	t.Logf("%08b", currentMask)
	t.Logf("%v", currentMask|0b0000 == 0)
	t.Logf("%v", currentMask&0b1111 != 0)
	t.Logf("%v", currentMask^0b0000 == 0)
}

func TestResolveCurrentIndexForLoop(t *testing.T) {
	lo := archsimd.LoadInt32x4(&[4]int32{0, 0, -1, 0})
	hi := archsimd.LoadInt32x4(&[4]int32{0, 0, 0, 0})
	index := resolveCurrentIndexForLoop(lo, hi, 32, -1)
	assert.Equal(t, 34, index)
}

func TestShadowRayNoIntersection(t *testing.T) {

	dir := &std.Vec4{0.45754904, 0.50628257, -0.7309766, 0}
	orig := &std.Vec4{-1.9835613, 0.9921495, 6.364105, 0}

	// translate origin a bit along the direction to the light
	origT := &std.Vec4{}
	origT[0] = orig[0] + dir[0]*0.01
	origT[1] = orig[1] + dir[1]*0.01
	origT[2] = orig[2] + dir[2]*0.01

	ray := &std.Ray{
		Orig: origT, // Y:-0.68457
		Dir:  dir,
	}

	spheres := AsStructOfArrays(std.SixteenSpheres())
	assert.False(t, IntersectSpheresSIMDShadowRay(ray, spheres))
}

func TestShadowRayWithIntersection(t *testing.T) {

	dir := &std.Vec4{0.19469759, 0.60765696, -0.76996493, 0}
	orig := &std.Vec4{-0.123568535, -1.2277203, 8.398015, 0} // Move down a bit

	ray := &std.Ray{
		Orig: orig,
		Dir:  dir,
	}

	spheres := AsStructOfArrays(std.SixteenSpheres())
	assert.True(t, IntersectSpheresSIMDShadowRay(ray, spheres))
}

func TestFindLowest(t *testing.T) {
	minT := archsimd.LoadFloat32x8Slice([]float32{3, 2, 1, 3, 3, 7, 12, 3})
	hi := minT.GetHi()
	i := hi.Min(minT.GetLo())

	i = i.Min(i.SelectFromPair(1, 2, 3, 0, i))
	i = i.Min(i.SelectFromPair(2, 3, 0, 1, i))
	i = i.Min(i.SelectFromPair(3, 0, 1, 2, i))

	assert.Equal(t, float32(1), i.GetElem(0))
}

func TestIntersectSpheres(t *testing.T) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}

	spheres := AsStructOfArrays(std.SixteenSpheres())
	t0, hitIndex, hit := IntersectSpheres(ray, spheres)
	assert.True(t, hit)
	assert.Equal(t, 5, hitIndex)
	assert.InEpsilon(t, 17, t0, 0.1)
}

// Hit sphere index 5
func TestIntersectSpheresSIMD(t *testing.T) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}

	spheres := AsStructOfArrays(std.SixteenSpheres())
	t0, hitIndex, hit := IntersectSpheresSIMD(ray, spheres)

	assert.True(t, hit)
	assert.Equal(t, 5, hitIndex)
	assert.InEpsilon(t, 17, t0, 0.1)
}

// Hit sphere index 1
func TestIntersectSpheresHitLowIndexSIMD(t *testing.T) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.3130755, -0.24104044, -0.91863114, 0},
	}

	spheres := AsStructOfArrays(std.SixteenSpheres())
	t0, hitIndex, hit := IntersectSpheresSIMD(ray, spheres)

	assert.True(t, hit)
	assert.Equal(t, 1, hitIndex)
	assert.InEpsilon(t, 17, t0, 0.1)
}

func TestIntersectSpheresSIMDMinIndexScalar(t *testing.T) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}

	spheres := AsStructOfArrays(std.SixteenSpheres())
	t0, hitIndex, hit := IntersectSpheresSIMDMinIndexScalar(ray, spheres)

	assert.True(t, hit)
	assert.Equal(t, 5, hitIndex)
	assert.InEpsilon(t, 17, t0, 0.1)
}

func BenchmarkIntersectSpheres(b *testing.B) {
	ray := &std.Ray{
		Orig: &std.Vec4{3.1, 4.1, 20.1, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}
	spheres := AsStructOfArrays(StandardSpheres())
	for b.Loop() {
		_, _, _ = IntersectSpheres(ray, spheres)
	}
}

func Benchmark16IntersectSpheres(b *testing.B) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}
	spheres := AsStructOfArrays(std.SixteenSpheres())
	for b.Loop() {
		_, _, _ = IntersectSpheres(ray, spheres)
	}
}

func Benchmark16IntersectSpheresSIMD(b *testing.B) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}
	spheres := AsStructOfArrays(std.SixteenSpheres())
	for b.Loop() {
		_, _, _ = IntersectSpheresSIMD(ray, spheres)
	}
}

func Benchmark16IntersectSpheresSIMDMinIndexScalar(b *testing.B) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}
	spheres := AsStructOfArrays(std.SixteenSpheres())
	for b.Loop() {
		_, _, _ = IntersectSpheresSIMDMinIndexScalar(ray, spheres)
	}
}

func BenchmarkSTDIntersectSpheres(b *testing.B) {
	spheres := std.SixteenSpheres()
	ray := &std.Ray{
		Orig: &std.Vec4{3.1, 4.1, 20.1, 0},
		Dir:  &std.Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}
	for b.Loop() {
		for _, s := range spheres {
			_, _ = std.IntersectSphere(ray, s.Center, s.RadiusSquared)
		}
	}
}

func BenchmarkSTDIntersectSpheresSIMD(b *testing.B) {
	spheres := std.SixteenSpheres()
	ray := &std.Ray{
		Orig: &std.Vec4{3.1, 4.1, 20.1, 0},
		Dir:  &std.Vec4{-0.47163668, 0.25912055, -0.8428614, 0},
	}
	for b.Loop() {
		for _, s := range spheres {
			_, _ = std.IntersectSphereSIMD(ray, s.Center, s.RadiusSquared)
		}
	}
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
