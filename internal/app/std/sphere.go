package std

import (
	"math"
	"simd"
)

// IntersectSphere performs a ray-sphere intersection using "plain go" DotProduct.
func IntersectSphere(ray *Ray, center *Vec4, radiusSquared float32) (float32, bool) {

	// Vector subtraction to get new vector from sphere center to ray origin
	var L Vec4
	Sub(center, ray.Orig, &L)

	// Compute dot product of L and ray direction. If tca is negative, no intersection.
	tca := L.DotProduct(ray.Dir)
	if tca < 0 {
		return 0.0, false
	}

	// Dot product of L with itself
	lD := L.DotProduct(&L)
	d2 := lD - tca*tca

	// If larger than radius squared, no intersection
	if d2 > radiusSquared {
		return 0.0, false
	}

	// Compute distance to hit point.
	thc := float32(math.Sqrt(float64(radiusSquared - d2)))
	t0 := tca - thc
	t1 := tca + thc

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
		return 0.0, false
	}

	t := t0

	return t, true
}

// IntersectSphereSIMD performs a ray-sphere intersection using DotProductSIMDVec4.
func IntersectSphereSIMD(ray *Ray, center *Vec4, radiusSquared float32) (float32, bool) {

	var L Vec4
	Sub(center, ray.Orig, &L)

	tca := DotProductSIMDVec4(&L, ray.Dir)
	if tca < 0 {
		return 0.0, false
	}

	d2 := DotProductSIMDVec4(&L, &L) - tca*tca
	if d2 > radiusSquared {
		return 0.0, false
	}

	thc := float32(math.Sqrt(float64(radiusSquared - d2)))
	t0 := tca - thc
	t1 := tca + thc

	// Swap
	if t0 > t1 {
		tmp := t0
		t0 = t1
		t1 = tmp
	}

	if t0 < 0 {
		t0 = t1 // If t0 is negative, let's use t1 instead.
	}
	if t0 < 0 {
		return 0.0, false // Both t0 and t1 are negative.
	}

	return t0, true
}

// IntersectSphereSIMDMethod is identical to IntersectSphereSIMD, except that it calls dotProductSIMD as a method on
// Vec4. It typically performs slightly worse than the func-based IntersectSphereSIMD approach.
func IntersectSphereSIMDMethod(ray *Ray, center *Vec4, radiusSquared float32) (float32, bool) {

	var L = Vec4{0, 0, 0, 0}
	Sub(center, ray.Orig, &L)

	tca := L.dotProductSIMD(ray.Dir)
	if tca < 0 {
		return 0.0, false
	}

	d2 := L.dotProductSIMD(&L) - tca*tca
	if d2 > radiusSquared {
		return 0.0, false
	}

	thc := float32(math.Sqrt(float64(radiusSquared - d2)))
	t0 := tca - thc
	t1 := tca + thc

	// Swap
	if t0 > t1 {
		tmp := t0
		t0 = t1
		t1 = tmp
	}

	if t0 < 0 {
		t0 = t1 // If t0 is negative, let's use t1 instead.
	}
	if t0 < 0 {
		return 0.0, false // Both t0 and t1 are negative.
	}

	return t0, true
}

var zeroes = simd.BroadcastFloat32x4(0.0)

// IntersectSphereFullSIMD is a WiP rewrite of ray-sphere intersection, trying to use SIMD more broadly.
func IntersectSphereFullSIMD(ray *Ray, center *Vec4, radiusSquared float32) (float32, bool) {
	var t0, t1 float32

	r1 := simd.LoadFloat32x4((*[4]float32)(center))
	r2 := simd.LoadFloat32x4((*[4]float32)(ray.Orig))
	rDir := simd.LoadFloat32x4((*[4]float32)(ray.Dir))

	// Vector from ray origin to center of sphere
	L := r1.Sub(r2)

	// Dot product of L and ray direction
	sumdst := simd.Float32x4{}
	sumdst = L.MulAdd(rDir, sumdst)  // => [6,12,20,30]
	sumdst = sumdst.AddPairs(sumdst) // => [18,50,18,50]
	sumdst = sumdst.AddPairs(sumdst) // => [68,68,68,68]

	//mask := sumdst.Less(zeroes)
	tca := sumdst.GetElem(0)
	// If dot product is negative, the ray cannot intersect.
	if tca < 0 {
		return 0.0, false
	}

	// Dot product of L and L
	sumdst2 := simd.Float32x4{}
	sumdst2 = L.MulAdd(L, sumdst2)
	sumdst2 = sumdst2.AddPairs(sumdst2) // => [18,50,18,50]
	sumdst2 = sumdst2.AddPairs(sumdst2) // => [68,68,68,68]
	tcaSquared := simd.BroadcastFloat32x4(tca)
	tcaSquared = tcaSquared.Mul(tcaSquared)
	d2 := sumdst2.Sub(tcaSquared).GetElem(0)

	// If d2 is larger than the squared radius of sphere we have no intersection.
	if d2 > radiusSquared {
		return 0.0, false
	}

	// irritating float32 <-> float64 type conversions
	thc := float32(math.Sqrt(float64(radiusSquared - d2)))
	t0 = tca - thc
	t1 = tca + thc

	if t0 > t1 {
		tmp := t0
		t0 = t1
		t1 = tmp
	}

	if t0 < 0 {
		t0 = t1 // If t0 is negative, let's use t1 instead.
	}
	if t0 < 0 {
		return 0.0, false // Both t0 and t1 are negative.
	}

	t := t0

	return t, true
}

// IntersectSphereDotGoSIMD uses https://go.dev/play/p/NY5rJYPoJcl for its Dot Product computations.
func IntersectSphereDotGoSIMD(ray *Ray, center *Vec4, radiusSquared float32) (float32, bool) {
	var t0, t1 float32 // solutions for t if the ray intersects

	// Geometric solution
	var L Vec4
	Sub(center, ray.Orig, &L)

	tca := DotGoSIMD(L[:], ray.Dir[:])
	if tca < 0 {
		return 0.0, false
	}
	d2 := DotGoSIMD(L[:], L[:]) - tca*tca
	if d2 > radiusSquared {
		return 0.0, false
	}
	thc := float32(math.Sqrt(float64(radiusSquared - d2)))
	t0 = tca - thc
	t1 = tca + thc

	if t0 > t1 {
		tmp := t0
		t0 = t1
		t1 = tmp
	}

	if t0 < 0 {
		t0 = t1
	}
	if t0 < 0 {
		return 0.0, false // Both t0 and t1 are negative.
	}

	t := t0

	return t, true
}

type Sphere struct {
	Center        *Vec4
	Radius        float32
	RadiusSquared float32
	Color         Vec4
}
