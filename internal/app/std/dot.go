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
func DotProductSIMDSlice64(v1, v2 []float64) float64 {

	r1 := simd.LoadFloat64x4((*[4]float64)(v1))
	r2 := simd.LoadFloat64x4((*[4]float64)(v2))
	sumdst := simd.Float64x4{}

	sumdst = r1.MulAdd(r2, sumdst)       // => [6,12,20,30]
	sumdst2 := sumdst.AddPairs(sumdst)   // => [18,50,18,50]
	sumdst3 := sumdst2.AddPairs(sumdst2) // => [68,68,68,68]
	out := [4]float64{}
	sumdst3.Store(&out)
	return out[zero]
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
	sumdst := simd.Float32x8{}

	r1 := simd.LoadFloat32x8((*[8]float32)(v1))
	r2 := simd.LoadFloat32x8((*[8]float32)(v2))

	sumdst = r1.MulAdd(r2, sumdst)   // => [6,12,20,30,6,12,20,35]
	sumdst = sumdst.AddPairs(sumdst) // => [18,50,18,50,18,55,18,55]
	sumdst = sumdst.AddPairs(sumdst) // => [68,68,68,68,73,73,73,73]
	out := [8]float32{}
	sumdst.Store(&out)

	return out[0], out[4]
}
