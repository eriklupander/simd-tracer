package adv

import (
	"math"
	"simd"

	"github.com/eriklupander/simd-tracer/internal/app/std"
)

// Spheres is a so-called "struct of arrays", essentially "transposing" from:
//
// xyzw           xxxxx
// xyzw           yyyyy
// xyzw    =>     zzzzz
// xyzw
// xyzw
//
// This is a good preparatory step for SIMD parallelization. It can also facilitate better memory layout
type Spheres struct {
	CenterX       []float32
	CenterY       []float32
	CenterZ       []float32
	RadiusSquared []float32
	Count         int
}

// IntersectSpheres performs a ray-spheres intersection given the "struct of arrays" Spheres type.
func IntersectSpheres(r *std.Ray, spheres *Spheres) (float32, int, bool) {
	minT := float32(3.4028235e38)
	currentIndex := -1
	for i := 0; i < spheres.Count; i++ {
		// Compute vector from sphere center to ray origin
		ocX := spheres.CenterX[i] - r.Orig[0]
		ocY := spheres.CenterY[i] - r.Orig[1]
		ocZ := spheres.CenterZ[i] - r.Orig[2]

		// Do dot product of center-origin vector and direction
		tca := ocX*r.Dir[0] + ocY*r.Dir[1] + ocZ*r.Dir[2] //+ v1[3]*v2[3]

		if tca < 0 {
			continue
		}

		// Dot product of L with itself
		lD := ocX*ocX + ocY*ocY + ocZ*ocZ
		d2 := lD - tca*tca

		if d2 > spheres.RadiusSquared[i] {
			continue
		}

		// Compute distance to hit point.
		beforeSqrt := spheres.RadiusSquared[i] - d2
		thc := float32(math.Sqrt(float64(beforeSqrt)))
		t0 := tca - thc
		t1 := tca + thc
		//fmt.Printf("%d: tca: %f ld: %f tcaSquared: %f d2: %f beforeSqrt: %f thc: %f t0: %f, t1: %f\n", i, tca, lD, tca*tca, d2, beforeSqrt, thc, t0, t1)

		// Swap if t0 is greater than t1
		if t0 > t1 {
			tmp := t0
			t0 = t1
			t1 = tmp
		}

		// If t0 is negative, let's use t1 instead.
		if t0 < 0 {
			t0 = t1
		}

		// Both t0 and t1 are negative. No intersection.
		if t0 < 0 {
			continue
		}

		// Check if t0 is smallest yet, in that case, store index and new lowest t
		if t0 < minT {
			minT = t0
			currentIndex = i
		}
	}
	if currentIndex != -1 {
		return minT, currentIndex, true
	}
	return 0.0, -1, false
}

var zeroes = simd.BroadcastFloat32x8(0.0001)
var maxT = simd.BroadcastFloat32x8(float32(3.4028235e38))
var firstElemSetMask = simd.LoadInt32x4(&[4]int32{-1, 0, 0, 0})

// IntersectSpheresSIMD performs a ray-spheres intersection given the "struct of arrays" Spheres type using SIMD.
func IntersectSpheresSIMD(r *std.Ray, spheres *Spheres) (float32, int, bool) {
	rayOriginX := simd.BroadcastFloat32x8(r.Orig[0])
	rayOriginY := simd.BroadcastFloat32x8(r.Orig[1])
	rayOriginZ := simd.BroadcastFloat32x8(r.Orig[2])
	rayDirectionX := simd.BroadcastFloat32x8(r.Dir[0])
	rayDirectionY := simd.BroadcastFloat32x8(r.Dir[1])
	rayDirectionZ := simd.BroadcastFloat32x8(r.Dir[2])

	currentMin := simd.BroadcastFloat32x4(math.MaxFloat32)

	currentIndex := -1
	it := -1
	for i := 0; i < spheres.Count; i += 8 {
		it++
		spheresCenterX := simd.LoadFloat32x8Slice(spheres.CenterX[i : i+8])
		spheresCenterY := simd.LoadFloat32x8Slice(spheres.CenterY[i : i+8])
		spheresCenterZ := simd.LoadFloat32x8Slice(spheres.CenterZ[i : i+8])
		spheresRadiusSquared := simd.LoadFloat32x8Slice(spheres.RadiusSquared[i : i+8])

		// Compute vectors from sphere center to ray origin
		ocX := spheresCenterX.Sub(rayOriginX)
		ocY := spheresCenterY.Sub(rayOriginY)
		ocZ := spheresCenterZ.Sub(rayOriginZ)

		// Dot products of center-origin vector and direction. tcas variable will hold 8 dot products
		tcas := ocX.Mul(rayDirectionX).Add(ocY.Mul(rayDirectionY)).Add(ocZ.Mul(rayDirectionZ))

		// If all tcas are less than zero, there cannot be any hits this iteration, continue
		if tcas.GreaterEqual(zeroes).AsInt32x8().IsZero() {
			continue
		}

		lD := ocX.Mul(ocX).Add(ocY.Mul(ocY)).Add(ocZ.Mul(ocZ)) //  lD holds 8 dot products
		tcaSquared := tcas.Mul(tcas)                           // tcaSquared holds the 8 dot products squared
		d2 := lD.Sub(tcaSquared)

		// Next mask, d2 must be greater than sphere radius squared. If no d2 fulfills that criteria, exit early.
		if d2.Less(spheresRadiusSquared).AsInt32x8().IsZero() {
			continue
		}

		// Compute distance to hit point.
		thc := spheresRadiusSquared.Sub(d2)
		thcSqrt := thc.Sqrt()
		t0 := tcas.Sub(thcSqrt)
		t1 := tcas.Add(thcSqrt)

		// t0 and t1 are the two intersection points on the sphere (typically enter and exit)
		// We need to find the lowest positive t0 / t1
		// Mask away negative values from both t0 and t1.
		// BUG HERE! When using t0.Greater, it worked when debugging in IDEA, but not when using run (or go test . on cmdline)
		t0PosMask := zeroes.Less(t0)
		t1PosMask := zeroes.Less(t1)

		t0Pos := t0.Masked(t0PosMask)
		t1Pos := t1.Masked(t1PosMask)

		minT := t0Pos.Min(t1Pos)
		if minT.Less(zeroes).AsInt32x8().IsZero() {
			continue
		}

		// Getting rid of zeroes by a masked merge, turning zeroes into float32 max.
		otherMask := minT.Greater(zeroes)
		minT = minT.Merge(maxT, otherMask)

		// Now, we need to find the smallest non-zero element in minT
		hi := minT.GetHi()
		lo := minT.GetLo()
		tX := hi.Min(lo) // min( [8,2,6,4], [5,1,8,3]) => [5,1,6,3]

		tX = tX.Min(tX.SelectFromPair(1, 2, 3, 0, tX)) // min( [5,1,6,3], [3,5,1,6]) => [3,1,1,3]
		tX = tX.Min(tX.SelectFromPair(2, 3, 0, 1, tX)) // min( [3,1,1,3], [3,3,1,1]) => [3,1,1,1]
		tX = tX.Min(tX.SelectFromPair(3, 0, 1, 2, tX)) // min( [3,1,1,1], [1,3,1,1]) => [1,1,1,1]

		// Check if this iteration's min is less than existing. If not we can skip
		// running the semi-expensive index
		msk := currentMin.Less(tX).AsInt32x4()
		if !msk.IsZero() {
			continue
		}
		currentMin = tX

		maskHi := currentMin.Equal(hi).AsInt32x4()
		maskLo := currentMin.Equal(lo).AsInt32x4()

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

	tMin := currentMin.GetElem(0)

	return tMin, currentIndex, currentIndex != -1
}

// IntersectSpheresSIMDExt performs a ray-spheres intersection given the "struct of arrays" Spheres type using SIMD.
func IntersectSpheresSIMDExt(r *std.Ray, spheres *Spheres, result []float32) {
	rayOriginX := simd.BroadcastFloat32x8(r.Orig[0])
	rayOriginY := simd.BroadcastFloat32x8(r.Orig[1])
	rayOriginZ := simd.BroadcastFloat32x8(r.Orig[2])
	rayDirectionX := simd.BroadcastFloat32x8(r.Dir[0])
	rayDirectionY := simd.BroadcastFloat32x8(r.Dir[1])
	rayDirectionZ := simd.BroadcastFloat32x8(r.Dir[2])

	for i := 0; i < spheres.Count; i += 8 {

		spheresCenterX := simd.LoadFloat32x8Slice(spheres.CenterX[i : i+8])
		spheresCenterY := simd.LoadFloat32x8Slice(spheres.CenterY[i : i+8])
		spheresCenterZ := simd.LoadFloat32x8Slice(spheres.CenterZ[i : i+8])
		spheresRadiusSquared := simd.LoadFloat32x8Slice(spheres.RadiusSquared[i : i+8])

		// Compute vectors from sphere center to ray origin
		ocX := spheresCenterX.Sub(rayOriginX)
		ocY := spheresCenterY.Sub(rayOriginY)
		ocZ := spheresCenterZ.Sub(rayOriginZ)

		// Dot products of center-origin vector and direction
		tcas := ocX.Mul(rayDirectionX).Add(ocY.Mul(rayDirectionY)).Add(ocZ.Mul(rayDirectionZ)) // tcas holds 8 dot products

		// Mask off all spheres whose tca is less than zero.
		tcaMask := tcas.GreaterEqual(zeroes)
		if tcaMask.AsInt32x8().IsZero() {
			continue
		}

		lD := ocX.Mul(ocX).Add(ocY.Mul(ocY)).Add(ocZ.Mul(ocZ)) //  lD holds 8 dot products
		tcaSquared := tcas.Mul(tcas)                           // tcaSquared holds the 8 dot products squared
		d2 := lD.Sub(tcaSquared)

		// d2 must be greater than sphere radius squared
		d2Mask := d2.Greater(spheresRadiusSquared)
		if d2Mask.AsInt32x8().IsZero() {
			continue
		}

		// Compute distance to hit point.
		thc := spheresRadiusSquared.Sub(d2)
		thcSqrt := thc.Sqrt()
		t0 := tcas.Sub(thcSqrt)
		t1 := tcas.Add(thcSqrt)

		// t0 and t1 are the two intersection points on the sphere (typically enter and exit)
		// We need to find the lowest positive t0 / t1
		// Mask away negative values from both t0 and t1.
		t0PosMask := zeroes.Less(t0)
		t1PosMask := zeroes.Less(t1)
		t0Pos := t0.Masked(t0PosMask)
		t1Pos := t1.Masked(t1PosMask)

		minT := t0Pos.Min(t1Pos)
		minT.StoreSlice(result[i : i+8])
	}

}
