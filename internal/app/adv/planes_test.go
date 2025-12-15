package adv

import (
	"math"
	"testing"

	"github.com/eriklupander/simd-tracer/internal/app/std"
	"github.com/stretchr/testify/assert"
)

func TestPlanes(t *testing.T) {
	dir := &std.Vec4{0.45754904, 0.50628257, -0.7309766, 0}
	orig := &std.Vec4{-1.9835613, 0.9921495, 6.364105, 0}

	planes := std.CornellBox()
	minT := float32(math.MaxFloat32)
	idx := -1
	for i, p := range planes {
		thisT, hit := IntersectPlane2(p.Normal, p.Point, orig, dir)

		if hit && thisT < minT {
			minT = thisT
			idx = i
		}
	}
	assert.InEpsilon(t, 8.9, minT, 0.2)
	assert.Equal(t, 3, idx)
}

func TestPlanesSIMD(t *testing.T) {

	dir := &std.Vec4{0.45754904, 0.50628257, -0.7309766, 0}
	orig := &std.Vec4{-1.9835613, 0.9921495, 6.364105, 0}
	ray := &std.Ray{
		Orig: orig,
		Dir:  dir,
	}

	planes := PlanesAsStructOfArrays(std.CornellBox())
	tt, idx, hit := IntersectPlanesSIMD(ray, planes)
	assert.InEpsilon(t, 8.9, tt, 0.2)
	assert.Equal(t, 3, idx)
	assert.True(t, hit)
}

func BenchmarkIntersectPlane(b *testing.B) {
	dir := &std.Vec4{0.45754904, 0.50628257, -0.7309766, 0}
	orig := &std.Vec4{-1.9835613, 0.9921495, 6.364105, 0}

	planes := std.CornellBox()
	minT := float32(math.MaxFloat32)

	for b.Loop() {
		for _, p := range planes {
			thisT, hit := IntersectPlane2(p.Normal, p.Point, orig, dir)

			if hit && thisT < minT {
				minT = thisT
			}
		}
	}

}

func BenchmarkIntersectPlaneSIMD(b *testing.B) {
	dir := &std.Vec4{0.45754904, 0.50628257, -0.7309766, 0}
	orig := &std.Vec4{-1.9835613, 0.9921495, 6.364105, 0}
	ray := &std.Ray{
		Orig: orig,
		Dir:  dir,
	}
	planes := PlanesAsStructOfArrays(std.CornellBox())
	for b.Loop() {
		_, _, _ = IntersectPlanesSIMD(ray, planes)
	}
}
