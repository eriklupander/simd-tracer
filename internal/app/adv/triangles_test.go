package adv

import (
	"fmt"
	"simd/archsimd"
	"testing"

	"github.com/eriklupander/simd-tracer/internal/app/std"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMasks(t *testing.T) {
	v1 := archsimd.LoadFloat32x4(&[4]float32{3, -1, 0, 4})
	zero := archsimd.BroadcastFloat32x4(0.0)
	msk := v1.Greater(zero)
	fmt.Println(msk.ToInt32x4())
}

func TestIntersectTriangleHit(t *testing.T) {
	ray := &std.Ray{Orig: &std.Vec4{0, 0.5, -2}, Dir: &std.Vec4{0, 0, 1}}
	tHit, _, hit := IntersectTriangle(ray, &std.Vec4{0, 1, 0}, &std.Vec4{-1, 0, 0}, &std.Vec4{1, 0, 0})
	assert.True(t, hit)
	assert.Equal(t, float32(2.0), tHit)
}
func TestIntersectTriangleMiss(t *testing.T) {
	ray := &std.Ray{Orig: &std.Vec4{0, -1, -2}, Dir: &std.Vec4{0, 0, 1}}
	tHit, _, hit := IntersectTriangle(ray, &std.Vec4{0, 1, 0}, &std.Vec4{-1, 0, 0}, &std.Vec4{1, 0, 0})
	assert.False(t, hit)
	assert.Equal(t, float32(0), tHit)

	ray = &std.Ray{Orig: &std.Vec4{1, 1, -2}, Dir: &std.Vec4{0, 0, 1}}
	tHit, _, hit = IntersectTriangle(ray, &std.Vec4{0, 1, 0}, &std.Vec4{-1, 0, 0}, &std.Vec4{1, 0, 0})
	assert.False(t, hit)
	assert.Equal(t, float32(0), tHit)

}
func TestIntersectTriangleMTHit(t *testing.T) {
	ray := &std.Ray{Orig: &std.Vec4{0, 0.5, -2}, Dir: &std.Vec4{0, 0, 1}}
	tHit, u, v, _, hit := IntersectTriangleMT(ray, &std.Vec4{0, 1, 0}, &std.Vec4{-1, 0, 0}, &std.Vec4{1, 0, 0}, 0)
	assert.True(t, hit)
	assert.Equal(t, float32(2.0), tHit)
	assert.Equal(t, float32(0.25), u)
	assert.Equal(t, float32(0.25), v)
}
func TestIntersectTriangleMTMiss(t *testing.T) {
	ray := &std.Ray{Orig: &std.Vec4{0, 1.5, -2}, Dir: &std.Vec4{0, 0, 1}} // Shoot above, origin is +1 y-axis
	tHit, _, _, _, hit := IntersectTriangleMT(ray, &std.Vec4{0, 1, 0}, &std.Vec4{-1, 0, 0}, &std.Vec4{1, 0, 0}, 0)
	assert.False(t, hit)
	assert.Equal(t, float32(0), tHit)
}

func TestIntersectTrianglesSIMD(t *testing.T) {
	ray := &std.Ray{Orig: &std.Vec4{0, 0.5, -2}, Dir: &std.Vec4{0, 0, 1}}
	tris := TrianglesSideBySide()

	t0, u, v, _, hit := IntersectTrianglesSIMD(ray, tris)
	assert.True(t, hit)
	assert.Equal(t, float32(2.0), t0)
	assert.Equal(t, float32(0.25), u)
	assert.Equal(t, float32(0.25), v)
}

func TestIntersectTriangleOne(t *testing.T) {

	ray := &std.Ray{Orig: &std.Vec4{3, 4, 20}, Dir: &std.Vec4{-0.2463486, -0.2102048, -0.9461111}}
	tris := TrianglesSmallSideBySide().ToSlice()
	tHit, u, v, _, hit := IntersectTriangleMT(ray, tris[1].v0, tris[1].v1, tris[1].v2, 1)
	assert.True(t, hit)
	assert.Equal(t, float32(20.082209), tHit)
	assert.Equal(t, float32(0.30791244), u)
	assert.Equal(t, float32(0.413464), v)
}
func TestIntersectTrianglesSIMDTriOne(t *testing.T) {

	ray := &std.Ray{Orig: &std.Vec4{3, 4, 20}, Dir: &std.Vec4{-0.2463486, -0.2102048, -0.9461111}}
	tris := TrianglesSmallSideBySide()

	t0, u, v, idx, hit := IntersectTrianglesSIMD(ray, tris)
	assert.True(t, hit)
	assert.Equal(t, 1, idx)
	assert.Equal(t, float32(20.082209), t0)
	assert.Equal(t, float32(0.30791244), u)
	assert.Equal(t, float32(0.413464), v)
}

func TestIntersectTriangleSeven(t *testing.T) {

	ray := &std.Ray{Orig: &std.Vec4{3, 4, 20}, Dir: &std.Vec4{0.0009813828, -0.21688892, -0.97619575}}
	tris := TrianglesSmallSideBySide().ToSlice()
	tHit, u, v, _, hit := IntersectTriangleMT(ray, tris[6].v0, tris[6].v1, tris[6].v2, 6)
	assert.True(t, hit)
	assert.Equal(t, float32(19.46331), tHit)
	assert.Equal(t, float32(0.34158704), u)
	assert.Equal(t, float32(0.3797891), v)
}
func TestIntersectTrianglesSIMDTriSeven(t *testing.T) {
	// 320x350 0.0009813828 -0.21688892 -0.97619575
	ray := &std.Ray{Orig: &std.Vec4{3, 4, 20}, Dir: &std.Vec4{0.0009813828, -0.21688892, -0.97619575}}
	tris := TrianglesSmallSideBySide()

	t0, u, v, idx, hit := IntersectTrianglesSIMD(ray, tris)
	assert.True(t, hit)
	assert.Equal(t, 6, idx)
	assert.Equal(t, float32(19.46331), t0)
	assert.Equal(t, float32(0.34158704), u)
	assert.Equal(t, float32(0.3797891), v)
}
func TestIntersectTriangleMissAbove(t *testing.T) {
	// 210x310 above the triangles
	ray := &std.Ray{Orig: &std.Vec4{3, 4, 20}, Dir: &std.Vec4{-0.2129862, -0.13712811, -0.9673845}}
	tris := TrianglesSmallSideBySide().ToSlice()
	for idx := range tris {
		_, _, _, _, hit := IntersectTriangleMT(ray, tris[idx].v0, tris[idx].v1, tris[idx].v2, 1)
		require.False(t, hit)
	}
}
func TestIntersectTrianglesSIMDMissAbove(t *testing.T) {
	// 210x310
	ray := &std.Ray{Orig: &std.Vec4{3, 4, 20}, Dir: &std.Vec4{-0.2129862, -0.13712811, -0.9673845}}
	tris := TrianglesSmallSideBySide()

	_, _, _, _, hit := IntersectTrianglesSIMD(ray, tris)
	assert.False(t, hit)
	//assert.Equal(t, float32(0.30791244), u)
	//assert.Equal(t, float32(0.413464), v)
}

func TestIntersectTrianglesStackedSIMD(t *testing.T) {
	ray := &std.Ray{Orig: &std.Vec4{0, 0.5, -2}, Dir: &std.Vec4{0, 0, 1}}
	tris := TrianglesStacked()

	t0, u, v, _, hit := IntersectTrianglesSIMD(ray, tris)
	assert.True(t, hit)
	assert.Equal(t, float32(1.0), t0)
	assert.Equal(t, float32(0.25), u)
	assert.Equal(t, float32(0.25), v)
}

func TestIntersectTriangleShadow(t *testing.T) {
	// Ray from wall at 300x270 towards point light, should not intersect
	ray := &std.Ray{Orig: &std.Vec4{2.136053, 2.6783404, -1.98}, Dir: &std.Vec4{-0.040390607, 0.80798984, 0.5878104}}
	tris := TrianglesStacked()

	hit := IntersectTrianglesShadowRaySIMD(ray, tris)
	assert.False(t, hit)

	// Ray close to the roof near light source that seems to cast gigantic shadow even though there is no triangle there
	ray = &std.Ray{Orig: &std.Vec4{2.7929351, 5.086907, -1.98}, Dir: &std.Vec4{-0.36782667, 0.14523755, 0.9184822}}
	hit = IntersectTrianglesShadowRaySIMD(ray, tris)
	assert.False(t, hit)

	hit = IntersectSpheresSIMDShadowRay(ray, AsStructOfArrays(std.SixteenSpheres()))
	assert.False(t, hit)
}

func BenchmarkIntersectTriangle(b *testing.B) {
	ray := &std.Ray{Orig: &std.Vec4{0, 0.5, -2}, Dir: &std.Vec4{0, 0, 1}}
	v0 := &std.Vec4{0, 1, 0}
	v1 := &std.Vec4{-1, 0, 0}
	v2 := &std.Vec4{1, 0, 0}
	for b.Loop() {
		_, _, _ = IntersectTriangle(ray, v0, v1, v2)
	}
}

func BenchmarkIntersectTriangleMT(b *testing.B) {
	ray := &std.Ray{Orig: &std.Vec4{0, 0.5, -2}, Dir: &std.Vec4{0, 0, 1}}
	v0 := &std.Vec4{0, 1, 0}
	v1 := &std.Vec4{-1, 0, 0}
	v2 := &std.Vec4{1, 0, 0}
	for b.Loop() {
		_, _, _, _, _ = IntersectTriangleMT(ray, v0, v1, v2, 0)
	}
}
func BenchmarkIntersectTriangleMiss(b *testing.B) {
	ray := &std.Ray{Orig: &std.Vec4{0, -1, -2}, Dir: &std.Vec4{0, 0, 1}}
	v0 := &std.Vec4{0, 1, 0}
	v1 := &std.Vec4{-1, 0, 0}
	v2 := &std.Vec4{1, 0, 0}
	for b.Loop() {
		_, _, _ = IntersectTriangle(ray, v0, v1, v2)
	}
}
func BenchmarkIntersectTriangleMissMT(b *testing.B) {
	ray := &std.Ray{Orig: &std.Vec4{0, -1, -2}, Dir: &std.Vec4{0, 0, 1}}
	v0 := &std.Vec4{0, 1, 0}
	v1 := &std.Vec4{-1, 0, 0}
	v2 := &std.Vec4{1, 0, 0}
	for b.Loop() {
		_, _, _, _, _ = IntersectTriangleMT(ray, v0, v1, v2, 0)
	}
}

func Benchmark8Triangles(b *testing.B) {
	ray := &std.Ray{Orig: &std.Vec4{0, 0.5, -2}, Dir: &std.Vec4{0, 0, 1}}
	triangles := TrianglesSideBySide().ToSlice()
	var minT float32
	for b.Loop() {
		minT = maxF32
		for _, tri := range triangles {
			t0, _, hit := IntersectTriangle(ray, tri.v0, tri.v1, tri.v2)
			if hit && t0 < minT {
				minT = t0
			}
		}
	}
}

func Benchmark8TrianglesMT(b *testing.B) {
	ray := &std.Ray{Orig: &std.Vec4{0, 0.5, -2}, Dir: &std.Vec4{0, 0, 1}}
	triangles := TrianglesSideBySide().ToSlice()
	var minT float32
	for b.Loop() {
		minT = maxF32
		for _, tri := range triangles {
			t0, _, _, _, hit := IntersectTriangleMT(ray, tri.v0, tri.v1, tri.v2, 0)
			if hit && t0 < minT {
				minT = t0
			}
		}
	}
}

func Benchmark8TrianglesSIMD(b *testing.B) {
	ray := &std.Ray{Orig: &std.Vec4{0, 0.5, -2}, Dir: &std.Vec4{0, 0, 1}}
	triangles := TrianglesSideBySide()

	for b.Loop() {
		_, _, _, _, _ = IntersectTrianglesSIMD(ray, triangles)
	}
}

func Benchmark8TrianglesStackedSIMD(b *testing.B) {
	ray := &std.Ray{Orig: &std.Vec4{0, 0.5, -2}, Dir: &std.Vec4{0, 0, 1}}
	triangles := TrianglesStacked()

	for b.Loop() {
		_, _, _, _, _ = IntersectTrianglesSIMD(ray, triangles)
	}
}
