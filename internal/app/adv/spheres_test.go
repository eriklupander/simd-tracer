package adv

import (
	"math/bits"
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

// TestShadowRaySixtyFourNoIntersection checks that a shadow ray from a floor hit
// (where no sphere should be in the path to the light) returns false.
// Floor is at y=-4, light at {2,5.4,0}. Shadow ray from (0,-4,5) toward light.
func TestShadowRaySixtyFourNoIntersection(t *testing.T) {
	// Shadow ray direction (unnormalized) = lightPos - hitPoint = (2, 9.4, -5)
	dir := &std.Vec4{0.1847, 0.8687, -0.4617, 0} // normalized
	orig := &std.Vec4{0.02, -3.906, 4.95, 0}     // hitPoint + epsilon*dir

	ray := &std.Ray{Orig: orig, Dir: dir}
	spheres := AsStructOfArrays(std.SixtyFourSpheres())
	assert.False(t, IntersectSpheresSIMDShadowRay(ray, spheres))
}

func TestShadowRayWithIntersection(t *testing.T) {
	// Ray origin is outside all spheres, pointing directly at sphere centered at {-0.5, -1, 3}.
	// t0 = dist - radius = 2 - 1 = 1 > 0, so this is a genuine (non-self) intersection.
	dir := &std.Vec4{0, 0, -1, 0}
	orig := &std.Vec4{-0.5, -1, 5, 0}

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

func Benchmark512IntersectSpheres(b *testing.B) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}
	spheres := AsStructOfArrays(std.SixtyFourSpheres())
	for b.Loop() {
		_, _, _ = IntersectSpheres(ray, spheres)
	}
}

func Benchmark512IntersectSpheresSIMD(b *testing.B) {
	ray := &std.Ray{
		Orig: &std.Vec4{3, 4, 20, 0},
		Dir:  &std.Vec4{-0.16645914, -0.24594483, -0.95488346, 0},
	}
	spheres := AsStructOfArrays(std.SixtyFourSpheres())
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
		Material:      std.Material{Color: std.Vec4{1, 0, 0}},
	})

	// 3 unit to the left (green)
	spheres = append(spheres, std.Sphere{
		Center:        &std.Vec4{-3, 0, 0, 0},
		Radius:        1,
		RadiusSquared: 1,
		Material:      std.Material{Color: std.Vec4{0, 1, 0}},
	})

	// 3 unit to the right (blue)
	spheres = append(spheres, std.Sphere{
		Center:        &std.Vec4{3, 0, 0, 0},
		Radius:        1,
		RadiusSquared: 1,
		Material:      std.Material{Color: std.Vec4{0, 0, 1}},
	})

	// 3 units up
	spheres = append(spheres, std.Sphere{
		Center:        &std.Vec4{0, 3, 0, 0},
		Radius:        1,
		RadiusSquared: 1,
		Material:      std.Material{Color: std.Vec4{0, 1, 1}},
	})

	// Closer to the camera
	spheres = append(spheres, std.Sphere{
		Center:        &std.Vec4{0, 0, 3, 0},
		Radius:        1,
		RadiusSquared: 1,
		Material:      std.Material{Color: std.Vec4{1, 1, 0}},
	})

	// Even closer
	spheres = append(spheres, std.Sphere{
		Center:        &std.Vec4{0, 0, 6, 0},
		Radius:        1,
		RadiusSquared: 1,
		Material:      std.Material{Color: std.Vec4{0.8, 0.8, 0.8}},
	})

	return spheres
}

func TestDivide(t *testing.T) {
	assert.Equal(t, 4, findBitIndex(0b00001000))
	assert.Equal(t, 6, findBitIndex(0b00000010))
	assert.Equal(t, 0, findBitIndex(0b10000000))
	assert.Equal(t, 2, findBitIndex(0b00100000))
}

var R int

func BenchmarkFindBitIndex(b *testing.B) {
	for b.Loop() {
		R = findBitIndex(0b00001000)
		R = findBitIndex(0b00000010)
		R = findBitIndex(0b10000000)
		R = findBitIndex(0b00100000)
		R = findBitIndex(0b00001000)
		R = findBitIndex(0b00000010)
		R = findBitIndex(0b10000000)
		R = findBitIndex(0b00100000)
		R = findBitIndex(0b00001000)
		R = findBitIndex(0b00000010)
		R = findBitIndex(0b10000000)
		R = findBitIndex(0b00100000)
	}
}

func BenchmarkFindBitIndex2(b *testing.B) {
	for b.Loop() {
		R = findCurrentIndex(0b00001000, 0)
		R = findCurrentIndex(0b00000010, 0)
		R = findCurrentIndex(0b10000000, 0)
		R = findCurrentIndex(0b00100000, 0)
		R = findCurrentIndex(0b00001000, 0)
		R = findCurrentIndex(0b00000010, 0)
		R = findCurrentIndex(0b10000000, 0)
		R = findCurrentIndex(0b00100000, 0)
		R = findCurrentIndex(0b00001000, 0)
		R = findCurrentIndex(0b00000010, 0)
		R = findCurrentIndex(0b10000000, 0)
		R = findCurrentIndex(0b00100000, 0)
	}
}

func findBitIndex(b1 uint8) int {
	if b1&0b1111 != 0 {
		// Low nibble
		if b1&0011 != 0 {
			// Low half
			if b1>>2&01 != 0 {
				return 5
			} else {
				return 4
			}
		} else {
			// Low half
			if b1&01 != 0 {
				return 7
			} else {
				return 6
			}
		}
	} else {
		// High nibble
		if b1>>4&0011 != 0 {
			// Low half
			if b1>>6&01 != 0 {
				return 1
			} else {
				return 0
			}
		} else {
			// Low half
			if b1>>4&01 != 0 {
				return 3
			} else {
				return 2
			}
		}
	}
}

func TestTrailing(t *testing.T) {
	assert.Equal(t, 6, bits.TrailingZeros8(0b01000000))
}
