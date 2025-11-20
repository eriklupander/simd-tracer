package std

import (
	"math"
	"simd"
)

func IntersectSphere(ray *Ray, center *Vec4, radiusSquared float32) (float32, bool) {
	var t0, t1 float32 // solutions for t if the ray intersects
	//#if 0
	// Geometric solution
	var L Vec4
	sub2(center, ray.Orig, &L) // -3,-1,-20 ray Orig is 3,4,20 and center is 0,3,0
	//, &L) // Vec3f
	tca := L.dotProduct(ray.Dir) // 20.22 ray Dir is   -0.15241566 -0.00019870223 -0.9883165 // DotProductSIMDSlice(L, ray.Dir)
	if tca < 0 {
		return 0.0, false
	}
	lD := L.dotProduct(&L) // 410
	d2 := lD - tca*tca     // 0,998 L.dotProduct(&L) - tca*tca
	if d2 > radiusSquared {
		return 0.0, false
	}
	thc := float32(math.Sqrt(float64(radiusSquared - d2))) // irritating float32 <-> float64 type conversions
	t0 = tca - thc
	t1 = tca + thc
	//#else
	// Analytic solution
	//L := sub(ray.Orig, center)
	//a := ray.Dir.dotProduct(ray.Dir)
	//b := 2.0 * ray.Dir.dotProduct(L)
	//c := L.dotProduct(L) - radius*radius
	//t0, t1, ok := solveQuadratic(a, b, c)
	//if !ok {
	//	return 0.0, false
	//}
	//#endif
	if t0 > t1 {
		tmp := t0
		t0 = t1
		t1 = tmp
	} //std::swap(t0, t1);

	if t0 < 0 {
		t0 = t1 // If t0 is negative, let's use t1 instead.
	}
	if t0 < 0 {
		return 0.0, false // Both t0 and t1 are negative.
	}

	t := t0

	return t, true
}

func IntersectSphereSIMD(ray *Ray, center *Vec4, radiusSquared float32) (float32, bool) {
	var t0, t1 float32 // solutions for t if the ray intersects
	//#if 0
	// Geometric solution
	var L Vec4
	sub2(center, ray.Orig, &L)
	//, &L) // Vec3f
	tca := DotProductSIMDVec4(&L, ray.Dir) //L.DotProductSIMDSlice(ray.Dir)
	if tca < 0 {
		return 0.0, false
	}
	d2 := DotProductSIMDVec4(&L, &L) - tca*tca // L.dotProduct(&L) - tca*tca
	if d2 > radiusSquared {
		return 0.0, false
	}
	thc := float32(math.Sqrt(float64(radiusSquared - d2))) // irritating float32 <-> float64 type conversions
	t0 = tca - thc
	t1 = tca + thc
	//#else
	// Analytic solution
	//L := sub(ray.Orig, center)
	//a := ray.Dir.dotProduct(ray.Dir)
	//b := 2.0 * ray.Dir.dotProduct(L)
	//c := L.dotProduct(L) - radius*radius
	//t0, t1, ok := solveQuadratic(a, b, c)
	//if !ok {
	//	return 0.0, false
	//}
	//#endif
	if t0 > t1 {
		tmp := t0
		t0 = t1
		t1 = tmp
	} //std::swap(t0, t1);

	if t0 < 0 {
		t0 = t1 // If t0 is negative, let's use t1 instead.
	}
	if t0 < 0 {
		return 0.0, false // Both t0 and t1 are negative.
	}

	t := t0

	return t, true
}

func IntersectSphereFullSIMD(ray *Ray, center *Vec4, radiusSquared float32) (float32, bool) {
	var t0, t1 float32 // solutions for t if the ray intersects
	// Geometric solution
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
	var d2 = float32(0.0)
	d2 = sumdst2.GetElem(0) - tca*tca

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

type Sphere struct {
	Center        *Vec4
	Radius        float32
	RadiusSquared float32
	Color         Vec4
}
