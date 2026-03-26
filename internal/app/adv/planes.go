//go:build goexperiment.simd && amd64

package adv

import (
	"math"
	"simd/archsimd"

	"github.com/eriklupander/simd-tracer/internal/app/std"
)

func IntersectPlane2(normal *std.Vec4, pointOnPlane *std.Vec4, orig *std.Vec4, dir *std.Vec4) (float32, bool) {
	lower := dir.DotProduct(normal)
	if lower > -0.00001 {
		return 0.0, false
	}
	var tmp1 std.Vec4
	std.Sub(pointOnPlane, orig, &tmp1)

	upper := tmp1.DotProduct(normal)

	t := upper / lower

	return t, t > 0.0
}

var epsilons = archsimd.BroadcastFloat32x8(0.01)

func IntersectPlanesSIMD(r *std.Ray, planes *Planes) (float32, int, bool) {
	rayOriginX := archsimd.BroadcastFloat32x8(r.Orig[0])
	rayOriginY := archsimd.BroadcastFloat32x8(r.Orig[1])
	rayOriginZ := archsimd.BroadcastFloat32x8(r.Orig[2])
	rayDirectionX := archsimd.BroadcastFloat32x8(r.Dir[0])
	rayDirectionY := archsimd.BroadcastFloat32x8(r.Dir[1])
	rayDirectionZ := archsimd.BroadcastFloat32x8(r.Dir[2])

	currentMin := archsimd.BroadcastFloat32x4(math.MaxFloat32)

	currentIndex := -1
	for i := 0; i < planes.Count; i += 8 {
		normalX := archsimd.LoadFloat32x8Slice(planes.NormalX[i : i+8])
		normalY := archsimd.LoadFloat32x8Slice(planes.NormalY[i : i+8])
		normalZ := archsimd.LoadFloat32x8Slice(planes.NormalZ[i : i+8])

		lower := rayDirectionX.Mul(normalX).Add(rayDirectionY.Mul(normalY)).Add(rayDirectionZ.Mul(normalZ))

		// If none of the values in lower are negative, we have no intersections
		if lower.GreaterEqual(zeroes).ToInt32x8().IsZero() {
			continue
		}

		// All positives in lower needs to be set to a small epsilon, or we'll either get the wrong result, or if set
		// to 0.0, we'll divide by zeroes when calculating the final t.
		negativesMask := lower.LessEqual(zeroes)
		lower = lower.Merge(epsilons, negativesMask)

		pointX := archsimd.LoadFloat32x8Slice(planes.PointX[i : i+8])
		pointY := archsimd.LoadFloat32x8Slice(planes.PointY[i : i+8])
		pointZ := archsimd.LoadFloat32x8Slice(planes.PointZ[i : i+8])

		hitFromOriginX := pointX.Sub(rayOriginX)
		hitFromOriginY := pointY.Sub(rayOriginY)
		hitFromOriginZ := pointZ.Sub(rayOriginZ)

		upper := hitFromOriginX.Mul(normalX).Add(hitFromOriginY.Mul(normalY)).Add(hitFromOriginZ.Mul(normalZ))

		t := upper.Div(lower) // t := upper / lower

		// t's > 0 are hits. Find the closest t for this iteration. We can first check so at least one t is positive.
		// (In a cornell box, we'll always have an intersection with a plane)
		negativeMask := t.LessEqual(zeroes)
		if negativeMask.ToInt32x8().IsZero() { // All values negative
			continue
		}

		// Getting rid of zeroes by a masked merge, turning zeroes into float32 max.
		otherMask := t.Greater(zeroes)
		t = t.Merge(maxT, otherMask)

		// Now, we need to find the smallest non-zero element in minT. This code is identical to the spheres one, should
		// be possible to use a common function - but be aware of inlining.
		hi := t.GetHi()
		lo := t.GetLo()
		tX := hi.Min(lo) // min( [8,2,6,4], [5,1,8,3]) => [5,1,6,3]

		tX = tX.Min(tX.SelectFromPair(1, 2, 3, 0, tX)) // min( [5,1,6,3], [3,5,1,6]) => [3,1,1,3]
		tX = tX.Min(tX.SelectFromPair(2, 3, 0, 1, tX)) // min( [3,1,1,3], [3,3,1,1]) => [3,1,1,1]
		tX = tX.Min(tX.SelectFromPair(3, 0, 1, 2, tX)) // min( [3,1,1,1], [1,3,1,1]) => [1,1,1,1]

		// Check if this iteration's min is less than existing. If not we can skip
		// running the semi-expensive code figuring out which sphere index we've
		// intersected as closest.
		msk := currentMin.Less(tX).ToInt32x4()
		if !msk.IsZero() {
			continue
		}
		currentMin = tX

		// The final step is to figure out index of the intersected plane.
		maskHi := currentMin.Equal(hi).ToInt32x4()
		maskLo := currentMin.Equal(lo).ToInt32x4()

		for ii := range 4 {
			if maskLo.Xor(firstElemSetMask).IsZero() {
				currentIndex = i + ii
				break
			}
			if maskHi.Xor(firstElemSetMask).IsZero() {
				currentIndex = i + 4 + ii
				break
			}
			// Permute all elements one step to the left.
			maskLo = maskLo.PermuteScalars(1, 2, 3, 0)
			maskHi = maskHi.PermuteScalars(1, 2, 3, 0)
		}
	}
	return currentMin.GetElem(0), currentIndex, currentIndex != -1
}

type Planes struct {
	NormalX []float32
	NormalY []float32
	NormalZ []float32

	PointX []float32
	PointY []float32
	PointZ []float32
	Color  []std.Vec4
	Count  int
}

func PlanesAsStructOfArrays(planes []std.Plane) *Planes {
	padSize := 8 - (len(planes) % 8)
	totalSize := len(planes) + padSize

	out := &Planes{
		NormalX: make([]float32, totalSize),
		NormalY: make([]float32, totalSize),
		NormalZ: make([]float32, totalSize),
		PointX:  make([]float32, totalSize),
		PointY:  make([]float32, totalSize),
		PointZ:  make([]float32, totalSize),
		Color:   make([]std.Vec4, totalSize),
		Count:   totalSize,
	}
	for i, pl := range planes {
		out.NormalX[i] = pl.Normal[0]
		out.NormalY[i] = pl.Normal[1]
		out.NormalZ[i] = pl.Normal[2]
		out.PointX[i] = pl.Point[0]
		out.PointY[i] = pl.Point[1]
		out.PointZ[i] = pl.Point[2]
		out.Color[i] = pl.Material.Color
	}

	// Ugly, but we need to pad to lane-size for SIMD, e.g. multiples of 8 if using AVX2
	for i := range padSize {
		idx := len(planes) + i
		out.NormalX[idx] = 0
		out.NormalY[idx] = 1
		out.NormalZ[idx] = 0
		out.PointX[idx] = 0
		out.PointY[idx] = -1000
		out.PointZ[idx] = 0
		out.Color[idx] = std.Vec4{0, 0, 0}
	}

	return out
}
