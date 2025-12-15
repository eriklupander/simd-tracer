package adv

import (
	"math"
	"simd"

	"github.com/eriklupander/simd-tracer/internal/app/std"
)

type Triangles struct {
	v0x []float32
	v0y []float32
	v0z []float32

	v1x []float32
	v1y []float32
	v1z []float32

	v2x []float32
	v2y []float32
	v2z []float32

	Count int
}

// IntersectTriangle intersects a triangle defined by the corners v0, v1 and v2 for the given ray.
// Code converted to Go from scratchapixel.com
func IntersectTriangle(ray *std.Ray, v0, v1, v2 *std.Vec4) (float32, int, bool) {

	// Compute the plane's normal
	v0v1 := std.Vec4{}
	std.Sub(v1, v0, &v0v1)

	v0v2 := std.Vec4{}
	std.Sub(v2, v0, &v0v2)

	// No need to normalize
	N := std.Vec4{}
	v0v1.CrossProduct(&v0v2, &N)

	// Step 1: Finding P

	// Check if the ray and plane are parallel
	NdotRayDirection := N.DotProduct(ray.Dir)
	if math.Abs(float64(NdotRayDirection)) < 0.01 {
		return 0, -1, false
	}

	// Compute d parameter using equation 2
	d := -N.DotProduct(v0)

	// Compute t (equation 3)
	t := -(N.DotProduct(ray.Orig) + d) / NdotRayDirection

	// Check if the triangle is behind the ray
	if t < 0 {
		return 0, -1, false // The triangle is behind
	}
	// Compute the intersection point using equation 1
	P := std.Vec4{}
	P[0] = ray.Orig[0] + ray.Dir[0]*t
	P[1] = ray.Orig[1] + ray.Dir[1]*t
	P[2] = ray.Orig[2] + ray.Dir[2]*t

	// Step 2: Inside-Outside Test
	Ne := std.Vec4{} // Vector perpendicular to triangle's plane

	// Test sidedness of P w.r.t. edge v0v1
	v0p := std.Vec4{}
	std.Sub(&P, v0, &v0p)
	v0v1.CrossProduct(&v0p, &Ne)
	if N.DotProduct(&Ne) < 0 {
		return 0, -1, false // P is on the right side
	}

	// Test sidedness of P w.r.t. edge v2v1
	v2v1 := std.Vec4{}
	v1p := std.Vec4{}
	std.Sub(v2, v1, &v2v1)
	std.Sub(&P, v1, &v1p)

	v2v1.CrossProduct(&v1p, &Ne)
	if N.DotProduct(&Ne) < 0 {
		return 0, -1, false // P is on the right side
	}

	// Test sidedness of P w.r.t. edge v2v0
	v2v0 := std.Vec4{}
	v2p := std.Vec4{}
	std.Sub(v0, v2, &v2v0)
	std.Sub(&P, v2, &v2p)

	v2v0.CrossProduct(&v2p, &Ne)
	if N.DotProduct(&Ne) < 0 {
		return 0, -1, false // P is on the right side
	}
	return t, 0, true // The ray hits the triangle
}

// IntersectTriangleMT implements Möller-Trumbore ray/triangle intersection (Converetd into Go from scratchapixel.com)
func IntersectTriangleMT(ray *std.Ray, v0, v1, v2 *std.Vec4, idx int) (float32, float32, float32, int, bool) {

	// Compute the plane's normal
	v0v1 := std.Vec4{}
	std.Sub(v1, v0, &v0v1)

	v0v2 := std.Vec4{}
	std.Sub(v2, v0, &v0v2)

	pvec := std.Vec4{}
	ray.Dir.CrossProduct(&v0v2, &pvec) //Vec3f pvec = dir.crossProduct(v0v2);

	det := v0v1.DotProduct(&pvec)

	// If det is close to 0, the ray and triangle are parallel.
	if math.Abs(float64(det)) < 0.01 {
		return 0, 0, 0, -1, false
	}

	invDet := 1 / det

	tvec := std.Vec4{}
	std.Sub(ray.Orig, v0, &tvec)

	u := tvec.DotProduct(&pvec) * invDet

	if u < 0 || u > 1 {
		return 0, 0, 0, -1, false
	}
	qvec := std.Vec4{}
	tvec.CrossProduct(&v0v1, &qvec)
	v := ray.Dir.DotProduct(&qvec) * invDet

	if v < 0 || u+v > 1 {
		return 0, 0, 0, -1, false
	}
	t := v0v2.DotProduct(&qvec) * invDet

	return t, u, v, 0, true
}

func IntersectTrianglesSIMD(r *std.Ray, triangles *Triangles) (float32, float32, float32, int, bool) {

	rayOriginX := simd.BroadcastFloat32x8(r.Orig[0])
	rayOriginY := simd.BroadcastFloat32x8(r.Orig[1])
	rayOriginZ := simd.BroadcastFloat32x8(r.Orig[2])
	rayDirectionX := simd.BroadcastFloat32x8(r.Dir[0])
	rayDirectionY := simd.BroadcastFloat32x8(r.Dir[1])
	rayDirectionZ := simd.BroadcastFloat32x8(r.Dir[2])

	one := simd.BroadcastFloat32x8(1.0)
	currentMin := simd.BroadcastFloat32x4(maxF32)
	currentU := simd.BroadcastFloat32x8(0)
	currentV := simd.BroadcastFloat32x8(0)

	var currentIndex = -1
	for i := 0; i < triangles.Count; i += 8 {
		v0x := simd.LoadFloat32x8Slice(triangles.v0x[i : i+8])
		v0y := simd.LoadFloat32x8Slice(triangles.v0y[i : i+8])
		v0z := simd.LoadFloat32x8Slice(triangles.v0z[i : i+8])

		v1x := simd.LoadFloat32x8Slice(triangles.v1x[i : i+8])
		v1y := simd.LoadFloat32x8Slice(triangles.v1y[i : i+8])
		v1z := simd.LoadFloat32x8Slice(triangles.v1z[i : i+8])

		v2x := simd.LoadFloat32x8Slice(triangles.v2x[i : i+8])
		v2y := simd.LoadFloat32x8Slice(triangles.v2y[i : i+8])
		v2z := simd.LoadFloat32x8Slice(triangles.v2z[i : i+8])

		// Compute the plane normal by creating two vectors from coord A to B and C, then cross product.
		v0v1x := v1x.Sub(v0x)
		v0v1y := v1y.Sub(v0y)
		v0v1z := v1z.Sub(v0z)

		v0v2x := v2x.Sub(v0x)
		v0v2y := v2y.Sub(v0y)
		v0v2z := v2z.Sub(v0z)

		// Cross product
		pVecX := rayDirectionY.Mul(v0v2z).Sub(rayDirectionZ.Mul(v0v2y))
		pVecY := rayDirectionZ.Mul(v0v2x).Sub(rayDirectionX.Mul(v0v2z))
		pVecZ := rayDirectionX.Mul(v0v2y).Sub(rayDirectionY.Mul(v0v2x))

		// DotProduct
		det := v0v1x.Mul(pVecX).Add(v0v1y.Mul(pVecY)).Add(v0v1z.Mul(pVecZ)) // det: 0.946111

		// TODO check for parallel using abs of det close to zero
		invDet := one.Div(det)

		tvecX := rayOriginX.Sub(v0x)
		tvecY := rayOriginY.Sub(v0y)
		tvecZ := rayOriginZ.Sub(v0z)

		u := (tvecX.Mul(pVecX).Add(tvecY.Mul(pVecY)).Add(tvecZ.Mul(pVecZ))).Mul(invDet)
		// 	u: 0.307912

		// Check if all u are < 0
		msk1 := u.Greater(zeroes)
		msk2 := u.Less(one)
		if msk1.AsInt32x8().IsZero() || msk2.AsInt32x8().IsZero() {
			continue
		}
		// We need to mask off any u that is NOT 0 to 1
		msk3 := msk1.And(msk2)
		u = u.Masked(msk3)

		qVecX := tvecY.Mul(v0v1z).Sub(tvecZ.Mul(v0v1y))
		qVecY := tvecZ.Mul(v0v1x).Sub(tvecX.Mul(v0v1z))
		qVecZ := tvecX.Mul(v0v1y).Sub(tvecY.Mul(v0v1x))

		v := (rayDirectionX.Mul(qVecX).Add(rayDirectionY.Mul(qVecY)).Add(rayDirectionZ.Mul(qVecZ))).Mul(invDet)
		v = v.Masked(msk3) // Mask off any triangles we know we cannot hit given mask from U
		//	v: 0.413464
		uv := u.Add(v)
		//fmt.Printf("u: %v\n", u)
		//fmt.Printf("v: %v\n", v)
		//fmt.Printf("uv: %v\n", uv)
		msk1x := v.Greater(zeroes)
		msk2x := uv.Less(one)
		//if v < 0 || u+v > 1 {
		if msk1x.AsInt32x8().IsZero() || msk2x.AsInt32x8().IsZero() {
			continue
		}

		// T now contains all intersection distances.
		t := (v0v2x.Mul(qVecX).Add(v0v2y.Mul(qVecY)).Add(v0v2z.Mul(qVecZ))).Mul(invDet)

		// Use mask to set all non-intersected t's to float32.Max
		t = t.Merge(maxT, msk1x.And(msk2x))
		if !findMinT(t, &currentMin) {
			continue
		}

		// We have a new closest intersection. We need to determine which index, which
		// is needed to figure out which u, v to store for this iteration.
		maskLo := currentMin.Equal(t.GetLo()).AsInt32x4()
		maskHi := currentMin.Equal(t.GetHi()).AsInt32x4()

		// Figure out _which_ index in the mask(s) that has value 1.
		currentIndex = resolveCurrentIndex(maskLo, maskHi, i, currentIndex)

		// Store this iteration's u and v result vectors. We'll extract the actual u and v's once the iterations are done.
		currentU = u
		currentV = v
	}
	if currentIndex == -1 {
		return 0, 0, 0, -1, false
	}
	// value is in lower or higher

	uvIndex := uint8(currentIndex % 8)
	all := currentU.SelectFromPairGrouped(uvIndex, uvIndex+4, uvIndex, uvIndex+4, currentV)
	if uvIndex < 4 {
		return currentMin.GetElem(0), all.GetLo().GetElem(0), all.GetLo().GetElem(1), currentIndex, true
	} else {
		return currentMin.GetElem(0), all.GetHi().GetElem(1), all.GetHi().GetElem(0), currentIndex, true
	}
	//// Shuffle so U ends up in element 0, V in element 1.
	//finalCoords := uVals.SelectFromPair(uvIndex, uvIndex+4, uvIndex, uvIndex+4, vVals)

}

func findMinT(minT simd.Float32x8, currentMin *simd.Float32x4) bool {
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
	msk := currentMin.Less(tX).AsInt32x4()
	if !msk.IsZero() {
		return false
	}

	// Current iteration has produced the lowest T (distance to intersection) yet. Assign tX to currentMin.
	*currentMin = tX
	return true
}

type Triangle struct {
	v0 *std.Vec4
	v1 *std.Vec4
	v2 *std.Vec4
}

func TrianglesSideBySide() *Triangles {
	tris := &Triangles{
		v0x: make([]float32, 8),
		v0y: make([]float32, 8),
		v0z: make([]float32, 8),
		v1x: make([]float32, 8),
		v1y: make([]float32, 8),
		v1z: make([]float32, 8),
		v2x: make([]float32, 8),
		v2y: make([]float32, 8),
		v2z: make([]float32, 8),
	}

	for i := range 8 {
		tris.v0x[i] = float32(i)
		tris.v0y[i] = 1
		tris.v0z[i] = 0

		tris.v1x[i] = -1 + float32(i)
		tris.v1y[i] = 0
		tris.v1z[i] = 0

		tris.v2x[i] = 1 + float32(i)
		tris.v2y[i] = 0
		tris.v2z[i] = 0
	}
	tris.Count = 8
	return tris
}

func TrianglesSmallSideBySide() *Triangles {
	tris := &Triangles{
		v0x: make([]float32, 8),
		v0y: make([]float32, 8),
		v0z: make([]float32, 8),
		v1x: make([]float32, 8),
		v1y: make([]float32, 8),
		v1z: make([]float32, 8),
		v2x: make([]float32, 8),
		v2y: make([]float32, 8),
		v2z: make([]float32, 8),
	}

	for i := range 8 {
		tris.v0x[i] = -3 + float32(i)*1
		tris.v0y[i] = 0.5
		tris.v0z[i] = 1

		tris.v1x[i] = -3.5 + float32(i)*1
		tris.v1y[i] = -0.5
		tris.v1z[i] = 1

		tris.v2x[i] = -2.5 + float32(i)*1
		tris.v2y[i] = -0.5
		tris.v2z[i] = 1
	}
	tris.Count = 8
	return tris
}

func TrianglesStacked() *Triangles {
	tris := &Triangles{
		v0x: make([]float32, 8),
		v0y: make([]float32, 8),
		v0z: make([]float32, 8),
		v1x: make([]float32, 8),
		v1y: make([]float32, 8),
		v1z: make([]float32, 8),
		v2x: make([]float32, 8),
		v2y: make([]float32, 8),
		v2z: make([]float32, 8),
	}

	for i := range 8 {
		tris.v0x[i] = 0
		tris.v0y[i] = 1
		tris.v0z[i] = -1 + float32(i)*0.5

		tris.v1x[i] = -1
		tris.v1y[i] = 0
		tris.v1z[i] = -1 + float32(i)*0.5

		tris.v2x[i] = 1
		tris.v2y[i] = 0
		tris.v2z[i] = -1 + float32(i)*0.5
	}
	tris.Count = 8
	return tris
}

func (t *Triangles) ToSlice() []Triangle {
	out := make([]Triangle, len(t.v0x))
	for i := range t.Count {
		tri := Triangle{
			v0: &std.Vec4{t.v0x[i], t.v0y[i], t.v0z[i]},
			v1: &std.Vec4{t.v1x[i], t.v1y[i], t.v1z[i]},
			v2: &std.Vec4{t.v2x[i], t.v2y[i], t.v2z[i]},
		}
		out[i] = tri
	}
	return out
}
