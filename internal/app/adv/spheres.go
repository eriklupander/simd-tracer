package adv

import (
	"math"
	"simd/archsimd"

	"github.com/eriklupander/simd-tracer/internal/app/std"
)

const maxF32 = float32(3.4028235e38)

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
	Color         []std.Vec4
	Count         int
}

// IntersectSpheres performs a ray-spheres intersection given the "struct of arrays" Spheres type.
func IntersectSpheres(r *std.Ray, spheres *Spheres) (float32, int, bool) {
	minT := maxF32
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

var true0 = archsimd.BroadcastFloat32x8(0)
var zeroes = archsimd.BroadcastFloat32x8(0.0001)
var maxT = archsimd.BroadcastFloat32x8(float32(3.4028235e38))
var negMaxT = archsimd.BroadcastFloat32x8(float32(-3.4028235e38))
var firstElemSetMask = archsimd.LoadInt32x4(&[4]int32{-1, 0, 0, 0})

// IntersectSpheresSIMD performs a ray-spheres intersection given the "struct of arrays" Spheres type using archsimd.
// Note the deliberate design where the code inside the for-statement is pretty much exclusively using branchless-style
// SIMD-only code, only using if-statements where we can exit early.
func IntersectSpheresSIMD(r *std.Ray, spheres *Spheres) (float32, int, bool) {
	rayOriginX := archsimd.BroadcastFloat32x8(r.Orig[0])
	rayOriginY := archsimd.BroadcastFloat32x8(r.Orig[1])
	rayOriginZ := archsimd.BroadcastFloat32x8(r.Orig[2])
	rayDirectionX := archsimd.BroadcastFloat32x8(r.Dir[0])
	rayDirectionY := archsimd.BroadcastFloat32x8(r.Dir[1])
	rayDirectionZ := archsimd.BroadcastFloat32x8(r.Dir[2])

	currentMin := archsimd.BroadcastFloat32x4(math.MaxFloat32)

	currentIndex := -1
	for i := 0; i < spheres.Count; i += 8 {
		spheresCenterX := archsimd.LoadFloat32x8Slice(spheres.CenterX[i : i+8])
		spheresCenterY := archsimd.LoadFloat32x8Slice(spheres.CenterY[i : i+8])
		spheresCenterZ := archsimd.LoadFloat32x8Slice(spheres.CenterZ[i : i+8])
		spheresRadiusSquared := archsimd.LoadFloat32x8Slice(spheres.RadiusSquared[i : i+8])

		// Compute vectors from sphere center to ray origin
		ocX := spheresCenterX.Sub(rayOriginX)
		ocY := spheresCenterY.Sub(rayOriginY)
		ocZ := spheresCenterZ.Sub(rayOriginZ)

		// Dot products of center-origin vector and direction. tcas variable will hold 8 dot products
		tcas := ocX.Mul(rayDirectionX).Add(ocY.Mul(rayDirectionY)).Add(ocZ.Mul(rayDirectionZ))

		// If all tcas are less than zero, there cannot be any hits this iteration, continue
		if tcas.GreaterEqual(zeroes).ToInt32x8().IsZero() {
			continue
		}

		// Compute the squared magnitude of the element-wise sphere-center -> origin vectors. (This is essentially dot product on itself)
		// lD then holds the squared magnitude per lane.
		lD := ocX.Mul(ocX).Add(ocY.Mul(ocY)).Add(ocZ.Mul(ocZ))

		// Then, also square tcas, tcaSquared holds the 8 dot products squared
		tcaSquared := tcas.Mul(tcas)

		// Element-wise subtracting of lD and tcaSquared produces
		d2 := lD.Sub(tcaSquared)

		// Next mask, d2 must be less than sphere radius squared to be a hit. If no d2 fulfills that criteria, exit early.
		if d2.Less(spheresRadiusSquared).ToInt32x8().IsZero() {
			continue
		}

		// Compute distance to hit point.
		thc := spheresRadiusSquared.Sub(d2)
		thcSqrt := thc.Sqrt()
		t0 := tcas.Sub(thcSqrt)
		t1 := tcas.Add(thcSqrt)

		// t0 and t1 are the two intersection points on the sphere (typically enter and exit, but there are edge cases)
		// We need to find the lowest positive t0 / t1
		// Mask away negative values from both t0 and t1.
		t0PosMask := zeroes.Less(t0)
		t1PosMask := zeroes.Less(t1)

		t0Pos := t0.Masked(t0PosMask)
		t1Pos := t1.Masked(t1PosMask)

		// Find min-values element-wise.
		minT := t0Pos.Min(t1Pos)

		// If all elements are negative, we can skip since there were no intersection.
		if minT.Less(zeroes).ToInt32x8().IsZero() {
			continue
		}

		// Get rid of zeroes in minT. First, create a mask for all elements in minT being greater than zero (the ones we want to keep)
		// Then, do a masked merge, turning elements where mask is false to float32 max value.
		otherMask := minT.Greater(zeroes) // [0,-2,2,0,3,0,-1,0] => greater mask => [0,0,1,0,1,0,0,0]
		minT = minT.Merge(maxT, otherMask)
		// [0,-2,2,0,3,0,-1,0] => merge([1e38,1e38,1e38,1e38], [0,0,1,0,1,0,0,0])
		// => [1e38,1e38,2,1e38,3,1e38,1e38,1e38]

		// Now, we need to find the smallest non-zero element in minT. Split our 8 elements into two four-element halves and
		// call element-wise min. This results in the 4 lowest of the 8 elements being retained.
		hi := minT.GetHi()
		lo := minT.GetLo()
		tX := hi.Min(lo) // min( [8,2,6,4], [5,1,8,3]) => [5,1,6,3]

		// Next, we need to combine shuffling with more Min-calls, i.e. "rotating" the elements sideways, performing a new
		// min call per iteration so we after three "rotations" have checked all elements against each other and the resulting
		// min-value is present in all 4 elements.
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

		// Current iteration has produced the lowest T (distance to intersection) yet. Assign tX to currentMin.
		currentMin = tX

		// Figure out which sphere (index) that produced the closest intersection. This is quite tricky, but we can take
		// advantage of already having the currentMin vector that has the lowest T in all elements. By performing an
		// element-wise equal operation against the lo and hi vectors containing the element-wise t values, we can produce
		// a mask indicating which element from 0-7 that matches the computed min T.
		//
		// Example:
		// [7,3,2,5] Equal [2,2,2,2] => mask [0,0,1,0]  (do for both hi and lo)
		maskLo := currentMin.Equal(lo).ToInt32x4()
		maskHi := currentMin.Equal(hi).ToInt32x4()

		// Finally, figure out _which_ index in the mask(s) that has value 1.
		currentIndex = resolveCurrentIndex(maskLo, maskHi, i, currentIndex)
	}

	// Extract scalar value for the tMin.
	tMin := currentMin.GetElem(0)

	return tMin, currentIndex, currentIndex != -1
}

//
//var zeroElemSet = [8]int32{-1, 0, 0, 0, 0, 0, 0, 0}
//var firstElemSet = [8]int32{0, -1, 0, 0, 0, 0, 0, 0}
//var secondElemSet = [8]int32{0, 0, -1, 0, 0, 0, 0, 0}
//var thirdElemSet = [8]int32{0, 0, 0, -1, 0, 0, 0, 0}
//var fourthElemSet = [8]int32{0, 0, 0, -1, 0, 0, 0, 0}
//var fifthElemSet = [8]int32{0, 0, 0, -1, 0, 0, 0, 0}
//var sixthElemSet = [8]int32{0, 0, 0, -1, 0, 0, 0, 0}
//var seventhElemSet = [8]int32{0, 0, 0, -1, 0, 0, 0, 0}
//var bitMap = map[[8]int32]int{
//	zeroElemSet: 0,
//	firstElemSet: 1,
//	thirdElemSet: 2,
//	zeroElemSet: 3,
//	fifthElemSet: 4,
//	zeroElemSet: 5,
//	sixthElemSet: 6,
//	seventhElemSet: 7,
//}

// resolveCurrentIndex: Given two 4x masks, indexed 0-3, 4-7, where at most one mask element is non-zero, this function
// returns the index of the set element, with offset added. If no match could be obtained,
// the supplied currentIndex is returned.
//
// This is probably an idiotic way to figure out which element (by index) that is set, but it works.
//
// The code can also be written using a for i := range 4 loop, but manually unrolling the loop improves performance. Note that
// the best performance is obtained if the code of this function is copied directly into the intersection code (e.g. manually
// inlined)
func resolveCurrentIndex(maskLo archsimd.Int32x4, maskHi archsimd.Int32x4, offset, currentIndex int) int {

	// Unroll #1
	if maskLo.Xor(firstElemSetMask).IsZero() {
		return offset
	}
	if maskHi.Xor(firstElemSetMask).IsZero() {
		return offset + 4
	}
	maskLo = maskLo.PermuteScalars(1, 2, 3, 0) // shift every element one step to the right, e.g. [0,0,1,0] => [0,0,0,1], next time [0,0,0,1] becomes [1,0,0,0] and so on.
	maskHi = maskHi.PermuteScalars(1, 2, 3, 0)

	// Unroll #2
	if maskLo.Xor(firstElemSetMask).IsZero() {
		return offset + 1
	}
	if maskHi.Xor(firstElemSetMask).IsZero() {
		return offset + 5
	}
	maskLo = maskLo.PermuteScalars(1, 2, 3, 0)
	maskHi = maskHi.PermuteScalars(1, 2, 3, 0)

	// Unroll #3
	if maskLo.Xor(firstElemSetMask).IsZero() {
		return offset + 2
	}
	if maskHi.Xor(firstElemSetMask).IsZero() {
		return offset + 6
	}
	maskLo = maskLo.PermuteScalars(1, 2, 3, 0)
	maskHi = maskHi.PermuteScalars(1, 2, 3, 0)

	// Unroll #4
	if maskLo.Xor(firstElemSetMask).IsZero() {
		return offset + 3
	}
	if maskHi.Xor(firstElemSetMask).IsZero() {
		return offset + 7
	}
	return currentIndex
}

// IntersectSpheresSIMDShadowRay performs a ray-spheres intersection given the "struct of arrays" Spheres type using archsimd.
func IntersectSpheresSIMDShadowRay(r *std.Ray, spheres *Spheres) bool {
	rayOriginX := archsimd.BroadcastFloat32x8(r.Orig[0])
	rayOriginY := archsimd.BroadcastFloat32x8(r.Orig[1])
	rayOriginZ := archsimd.BroadcastFloat32x8(r.Orig[2])
	rayDirectionX := archsimd.BroadcastFloat32x8(r.Dir[0])
	rayDirectionY := archsimd.BroadcastFloat32x8(r.Dir[1])
	rayDirectionZ := archsimd.BroadcastFloat32x8(r.Dir[2])

	for i := 0; i < spheres.Count; i += 8 {

		spheresCenterX := archsimd.LoadFloat32x8Slice(spheres.CenterX[i : i+8])
		spheresCenterY := archsimd.LoadFloat32x8Slice(spheres.CenterY[i : i+8])
		spheresCenterZ := archsimd.LoadFloat32x8Slice(spheres.CenterZ[i : i+8])
		spheresRadiusSquared := archsimd.LoadFloat32x8Slice(spheres.RadiusSquared[i : i+8])

		// Compute vectors from sphere center to ray origin
		ocX := spheresCenterX.Sub(rayOriginX)
		ocY := spheresCenterY.Sub(rayOriginY)
		ocZ := spheresCenterZ.Sub(rayOriginZ)

		// Dot products of center-origin vector and direction. tcas variable will hold 8 dot products
		tcas := ocX.Mul(rayDirectionX).Add(ocY.Mul(rayDirectionY)).Add(ocZ.Mul(rayDirectionZ))

		// If no values in tcas are greater than zero, there cannot be any hits this iteration, continue
		if tcas.GreaterEqual(zeroes).ToInt32x8().IsZero() {
			continue
		}

		// Compute element-wise dot products. (Each element in lD contains that column's dot product.
		lD := ocX.Mul(ocX).Add(ocY.Mul(ocY)).Add(ocZ.Mul(ocZ))
		tcaSquared := tcas.Mul(tcas)
		d2 := lD.Sub(tcaSquared)

		// Next mask, d2 must be greater than sphere radius squared. If no d2 fulfills that criteria, exit early.
		d2Mask := spheresRadiusSquared.Greater(d2)

		//fmt.Printf("D2Mask: %v\n", d2Mask)
		if d2Mask.ToInt32x8().IsZero() {
			continue
		}

		// Compute distance to hit point.
		thc := spheresRadiusSquared.Sub(d2)
		thcSqrt := thc.Sqrt()
		t0 := tcas.Sub(thcSqrt)
		t1 := tcas.Add(thcSqrt)

		zero := archsimd.BroadcastFloat32x8(0.0)
		m1 := t0.Greater(zero)
		zero = archsimd.BroadcastFloat32x8(0.0)
		m2 := t1.Greater(zero)

		if m1.ToInt32x8().IsZero() && m2.ToInt32x8().IsZero() {
			// No intersection
			continue
		}

		// If we've gotten here something (doesn't really matter what) intersected the shadow ray.
		return true
	}
	return false
}

func AsStructOfArrays(spheres []std.Sphere) *Spheres {
	sp := &Spheres{
		CenterX:       make([]float32, len(spheres)),
		CenterY:       make([]float32, len(spheres)),
		CenterZ:       make([]float32, len(spheres)),
		RadiusSquared: make([]float32, len(spheres)),
		Count:         len(spheres),
		Color:         make([]std.Vec4, len(spheres)),
	}
	for i, s := range spheres {
		sp.CenterX[i] = s.Center[0]
		sp.CenterY[i] = s.Center[1]
		sp.CenterZ[i] = s.Center[2]
		sp.RadiusSquared[i] = s.RadiusSquared
		sp.Color[i] = s.Color
	}
	return sp
}
