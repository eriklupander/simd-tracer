//go:build goexperiment.simd && amd64

package std

import (
	"simd"
)

func DotProduct(v1, v2 []float32) float32 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2] + v1[3]*v2[3]
}

func DotProduct2x4(v1, v2 []float32) (float32, float32) {
	sum := v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2] + v1[3]*v2[3]
	sum2 := v1[4]*v2[4] + v1[5]*v2[5] + v1[6]*v2[6] + v1[7]*v2[7]
	return sum, sum2
}

func DotProductSIMDManualAdd(v1, v2 []float32) (float32, float32) {
	sumdst := simd.Float32x8{}

	r1 := simd.LoadFloat32x8((*[8]float32)(v1))
	r2 := simd.LoadFloat32x8((*[8]float32)(v2))

	sumdst = r1.MulAdd(r2, sumdst)

	out := [8]float32{}
	sumdst.Store(&out)

	return out[0] + out[1] + out[2] + out[3], out[4] + out[5] + out[6] + out[7]
}

const zero = 0

// DotProductSIMDSlice performs a single dot product for the passed slices.
func DotProductSIMDSlice(v1, v2 []float32) float32 {
	r1 := simd.LoadFloat32x4((*[4]float32)(v1))
	r2 := simd.LoadFloat32x4((*[4]float32)(v2))
	sumdst := simd.Float32x4{}

	sumdst = r1.MulAdd(r2, sumdst)       // => [6,12,20,30]
	sumdst2 := sumdst.AddPairs(sumdst)   // => [18,50,18,50]
	sumdst3 := sumdst2.AddPairs(sumdst2) // => [68,68,68,68]

	return sumdst3.GetElem(zero)
}

// DotProductSIMDArray performs a single dot product for the passed arrayas
func DotProductSIMDArray(v1, v2 [4]float32) float32 {
	sumdst := simd.Float32x4{}

	r1 := simd.LoadFloat32x4(&v1)
	r2 := simd.LoadFloat32x4(&v2)

	sumdst = r1.MulAdd(r2, sumdst)   // => [6,12,20,30]
	sumdst = sumdst.AddPairs(sumdst) // => [18,50,18,50]
	sumdst = sumdst.AddPairs(sumdst) // => [68,68,68,68]
	return sumdst.GetElem(zero)
}

// DotProductSIMDArrayPointer performs a single dot product for the passed arrayas
func DotProductSIMDArrayPointer(v1, v2 *[4]float32) float32 {
	sumdst := simd.Float32x4{}

	r1 := simd.LoadFloat32x4(v1)
	r2 := simd.LoadFloat32x4(v2)

	sumdst = r1.MulAdd(r2, sumdst)   // => [6,12,20,30]
	sumdst = sumdst.AddPairs(sumdst) // => [18,50,18,50]
	sumdst = sumdst.AddPairs(sumdst) // => [68,68,68,68]
	return sumdst.GetElem(zero)
}

// DotProductSIMDArrayPointer8 performs a single dot product for the passed arrayas
func DotProductSIMDArrayPointer8(v1 [8]float32) float32 {
	sumdst := simd.Float32x4{}

	r1 := simd.LoadFloat32x4((*[4]float32)(v1[0:4]))
	r2 := simd.LoadFloat32x4((*[4]float32)(v1[4:8]))

	sumdst = r1.MulAdd(r2, sumdst)   // => [6,12,20,30]
	sumdst = sumdst.AddPairs(sumdst) // => [18,50,18,50]
	sumdst = sumdst.AddPairs(sumdst) // => [68,68,68,68]
	return sumdst.GetElem(zero)
}

func DotProductSIMDVec4(v1, v2 *Vec4) float32 {
	sumdst := simd.Float32x4{}

	r1 := simd.LoadFloat32x4((*[4]float32)(v1))
	r2 := simd.LoadFloat32x4((*[4]float32)(v2))

	sumdst = r1.MulAdd(r2, sumdst)   // => [6,12,20,30]
	sumdst = sumdst.AddPairs(sumdst) // => [18,50,18,50]
	sumdst = sumdst.AddPairs(sumdst) // => [68,68,68,68]
	return sumdst.GetElem(zero)
}

// DotProductSIMD2x4 computes full dot products for two 4-element vectors. Elements 0-3 form the first dot product,
// elements 4-7 form the second dot product.
func DotProductSIMD2x4(v1, v2 []float32) (float32, float32) {

	r1 := simd.LoadFloat32x8((*[8]float32)(v1)) // [2,3,4,5,2,3,4,5]
	r2 := simd.LoadFloat32x8((*[8]float32)(v2)) // [3,4,5,6,3,4,5,7]

	sumdst := simd.Float32x8{}
	sumdst = r1.MulAdd(r2, sumdst)   // => [6,12,20,30,6,12,20,35]
	sumdst = sumdst.AddPairs(sumdst) // => [18,50,18,50,18,55,18,55]
	sumdst = sumdst.AddPairs(sumdst) // => [68,68,68,68,73,73,73,73]
	out := [8]float32{}
	sumdst.Store(&out)

	return out[0], out[4]
}

// DotProductSIMD2x4HiLo is identical to DotProductSIMD2x4, except that it uses GetHi/GetLo with GetElem to extract
// the final dot products instead of performing a store.
func DotProductSIMD2x4HiLo(v1, v2 []float32) (float32, float32) {

	r1 := simd.LoadFloat32x8((*[8]float32)(v1))
	r2 := simd.LoadFloat32x8((*[8]float32)(v2))

	sumdst := simd.Float32x8{}
	sumdst = r1.MulAdd(r2, sumdst)   // => [6,12,20,30,6,12,20,35]
	sumdst = sumdst.AddPairs(sumdst) // => [18,50,18,50,18,55,18,55]
	sumdst = sumdst.AddPairs(sumdst) // => [68,68,68,68,73,73,73,73]
	lo := sumdst.GetLo()
	hi := sumdst.GetHi()

	return lo.GetElem(0), hi.GetElem(0)
}

// From https://go.dev/play/p/ZaruM_PzP1X
func DotAVX256a(x []float32, y []float32) float32 {
	var a simd.Float32x8
	if len(y) < len(x) {
		panic("slice y is shorter than slice x")
	}
	i := 0
	for ; i < len(x)-8; i += 8 { // this idiom is friendly to bounds check elimination
		xv := simd.LoadFloat32x8Slice(x[i : i+8])
		yv := simd.LoadFloat32x8Slice(y[i : i+8])
		a = yv.MulAdd(xv, a)
	}
	xv := simd.LoadFloat32x8SlicePart(x[i:])
	yv := simd.LoadFloat32x8SlicePart(y[i:])
	a = yv.MulAdd(xv, a)
	a = a.AddPairs(a) // 01234567                AP 01234567                -> 0+1 2+3 _ _ 4+5 6+7 _ _
	a = a.AddPairs(a) // 0+1 2+3 _ _ 4+5 6+7 _ _ AP 0+1 2+3 _ _ 4+5 6+7 _ _ -> 0+1+2+3 _ _ _ 4+5+6+7 _ _ _
	b := a.GetLo().Add(a.GetHi())
	return b.GetElem(0)
}

// From https://go.dev/play/p/NY5rJYPoJcl
func DotGoSIMD(x, y []float32) float32 {
	var (
		s0, s1, s2, s3 simd.Float32x8
	)

	// Writing anything slice indexing related in constant can reduce the bound checks.
	// Our bound-check elimination pass is clever at reasoning constants, but struggles
	// at reasoning expressions with variables.
	for len(x) >= 32 && len(y) >= 32 {
		x3 := simd.LoadFloat32x8Slice(x[24:])
		x2 := simd.LoadFloat32x8Slice(x[16:])
		x1 := simd.LoadFloat32x8Slice(x[8:])
		x0 := simd.LoadFloat32x8Slice(x[:])
		x = x[32:]
		y3 := simd.LoadFloat32x8Slice(y[24:])
		y2 := simd.LoadFloat32x8Slice(y[16:])
		y1 := simd.LoadFloat32x8Slice(y[8:])
		y0 := simd.LoadFloat32x8Slice(y[:])
		y = y[32:]

		s0 = x0.MulAdd(y0, s0)
		s1 = x1.MulAdd(y1, s1)
		s2 = x2.MulAdd(y2, s2)
		s3 = x3.MulAdd(y3, s3)
	}

	// Reduce to one value
	s0 = s0.Add(s1).Add(s2.Add(s3))
	low, high := s0.GetLo(), s0.GetHi()
	sum4 := low.Add(high)
	sum2 := sum4.AddPairs(sum4)
	sum1 := sum2.AddPairs(sum2)
	sum1Slice := make([]float32, 4)
	sum1.StoreSlice(sum1Slice)
	sum := sum1Slice[0]

	// Handle the tail.
	if len(x) == len(y) {
		// Again remove on unnecessary bound check.
		// Our bound-check elimination pass is also clever at reasoning ==.
		for i := range len(x) {
			sum += x[i] * y[i]
		}
	}

	return sum
}
