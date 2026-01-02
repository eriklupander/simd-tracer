package std

import (
	"simd/archsimd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDot(t *testing.T) {
	v1 := []float32{2, 3, 4, 5, 2, 3, 4, 5}
	v2 := []float32{3, 4, 5, 6, 3, 4, 5, 7}

	sum1, sum2 := DotProduct2x4(v1, v2)
	assert.EqualValues(t, 68.0, sum1)
	assert.EqualValues(t, 73.0, sum2)

	sum1, sum2 = DotProductSIMDManualAdd(v1, v2)
	assert.EqualValues(t, 68.0, sum1)
	assert.EqualValues(t, 73.0, sum2)
}

func TestDotProductSIMDMethod(t *testing.T) {
	v1 := &Vec4{2, 3, 4, 5}
	v2 := &Vec4{3, 4, 5, 6}
	sum := v1.dotProductSIMD(v2)
	assert.EqualValues(t, 68.0, sum)
}

func TestDotProductSIMDSlice(t *testing.T) {

	v1 := []float32{2, 3, 4, 5}
	v2 := []float32{3, 4, 5, 6}
	v3 := []float32{2, 3, 4, 5}
	v4 := []float32{3, 4, 5, 7}
	sum1 := DotProductSIMDSlice(v1, v2)
	assert.EqualValues(t, 68.0, sum1)
	sum2 := DotProductSIMDSlice(v3, v4)
	assert.EqualValues(t, 73.0, sum2)
}

//func TestDotProductSIMD64(t *testing.T) {
//
//	v1 := &Vec4D{2, 3, 4, 5}
//	v2 := &Vec4D{3, 4, 5, 6}
//
//	sum1 := DotProductSIMDVec4D(v1, v2)
//	assert.EqualValues(t, 68.0, sum1)
//}

func TestDotProductSIMDArrayPointer8(t *testing.T) {
	v1 := [8]float32{2, 3, 4, 5, 3, 4, 5, 6}
	sum1 := DotProductSIMDArrayPointer8(v1)
	assert.EqualValues(t, 68.0, sum1)
}

//func TestDotProductSIMDSlice64(t *testing.T) {
//
//	v1 := []float64{2, 3, 4, 5}
//	v2 := []float64{3, 4, 5, 6}
//	v3 := []float64{2, 3, 4, 5}
//	v4 := []float64{3, 4, 5, 7}
//	sum1 := DotProductSIMDSlice64(v1, v2)
//	assert.EqualValues(t, 68.0, sum1)
//	sum2 := DotProductSIMDSlice64(v3, v4)
//	assert.EqualValues(t, 73.0, sum2)
//}

func TestDotProductSIMDArray(t *testing.T) {

	v1 := [4]float32{2, 3, 4, 5}
	v2 := [4]float32{3, 4, 5, 6}
	v3 := [4]float32{2, 3, 4, 5}
	v4 := [4]float32{3, 4, 5, 7}
	sum1 := DotProductSIMDArray(v1, v2)
	assert.EqualValues(t, 68.0, sum1)
	sum2 := DotProductSIMDArray(v3, v4)
	assert.EqualValues(t, 73.0, sum2)
}

func TestDotProductSIMDArrayPointer(t *testing.T) {

	v1 := [4]float32{2, 3, 4, 5}
	v2 := [4]float32{3, 4, 5, 6}
	v3 := [4]float32{2, 3, 4, 5}
	v4 := [4]float32{3, 4, 5, 7}
	sum1 := DotProductSIMDArrayPointer(&v1, &v2)
	assert.EqualValues(t, 68.0, sum1)
	sum2 := DotProductSIMDArrayPointer(&v3, &v4)
	assert.EqualValues(t, 73.0, sum2)
}

func TestDotProductSIMD2x4(t *testing.T) {
	v1 := []float32{2, 3, 4, 5, 2, 3, 4, 5}
	v2 := []float32{3, 4, 5, 6, 3, 4, 5, 7}

	sum1, sum2 := DotProductSIMD2x4(v1, v2)
	assert.EqualValues(t, 68.0, sum1)
	assert.EqualValues(t, 73.0, sum2)
}

func TestDotProductSIMD2x4HiLo(t *testing.T) {
	v1 := []float32{2, 3, 4, 5, 2, 3, 4, 5}
	v2 := []float32{3, 4, 5, 6, 3, 4, 5, 7}

	sum1, sum2 := DotProductSIMD2x4HiLo(v1, v2)
	assert.EqualValues(t, 68.0, sum1)
	assert.EqualValues(t, 73.0, sum2)
}

func BenchmarkDot(b *testing.B) {
	v1 := []float32{2, 3, 4, 5}
	v2 := []float32{3, 4, 5, 6}
	tot := float32(0.0)
	for b.Loop() {
		tot += DotProduct(v1, v2)
	}
}
func BenchmarkDotSIMDMethod(b *testing.B) {
	v1 := &Vec4{2, 3, 4, 5}
	v2 := &Vec4{3, 4, 5, 6}
	tot := float32(0.0)
	for b.Loop() {
		tot += v1.dotProductSIMD(v2)
	}
}
func BenchmarkDotSIMDArray(b *testing.B) {
	v1 := [4]float32{2, 3, 4, 5}
	v2 := [4]float32{3, 4, 5, 6}
	tot := float32(0.0)
	for b.Loop() {
		tot += DotProductSIMDArray(v1, v2)
	}
}
func BenchmarkDotSIMDArrayPtr(b *testing.B) {
	v1 := &[4]float32{2, 3, 4, 5}
	v2 := &[4]float32{3, 4, 5, 6}
	tot := float32(0.0)
	for b.Loop() {
		tot += DotProductSIMDArrayPointer(v1, v2)
	}
}
func BenchmarkDotSIMDSlice(b *testing.B) {
	v1 := []float32{2, 3, 4, 5}
	v2 := []float32{3, 4, 5, 6}
	tot := float32(0.0)
	for b.Loop() {
		tot += dot(v1, v2)
	}
}

var Result float32

func BenchmarkDotSIMDSlice2(b *testing.B) {
	v1 := []float32{2, 3, 4, 5}
	v2 := []float32{3, 4, 5, 6}

	for b.Loop() {
		a := dot(v1, v2)
		b := dot(v1, v2)
		c := dot(v1, v2)
		d := dot(v1, v2)
		Result += a + b + c + d
	}
}

func dot(v1, v2 []float32) float32 {

	r1 := archsimd.LoadFloat32x4((*[4]float32)(v1))
	r2 := archsimd.LoadFloat32x4((*[4]float32)(v2))
	sums := archsimd.Float32x4{}

	sums = r1.MulAdd(r2, sums)
	sums = sums.AddPairs(sums)
	sums = sums.AddPairs(sums)
	return sums.GetElem(0)

}

func BenchmarkDot2x(b *testing.B) {
	v1 := []float32{2, 3, 4, 5}
	v2 := []float32{3, 4, 5, 6}
	v3 := []float32{2, 3, 4, 5}
	v4 := []float32{3, 4, 5, 6}
	tot := float32(0.0)
	for b.Loop() {
		tot += DotProduct(v1, v2)
		tot += DotProduct(v3, v4)
	}
}

func BenchmarkDot2x4_2(b *testing.B) {
	v1 := []float32{2, 3, 4, 5, 2, 3, 4, 5}
	v2 := []float32{3, 4, 5, 6, 3, 4, 5, 6}
	tot := float32(0.0)
	for b.Loop() {
		a, b := DotProduct2x4(v1, v2)
		tot = tot + a + b
	}
}

//	func BenchmarkDot2x4(b *testing.B) {
//		v1 := []float32{2, 3, 4, 5, 2, 3, 4, 5}
//		v2 := []float32{3, 4, 5, 6, 3, 4, 5, 6}
//
//		for b.Loop() {
//			_, _ = DotProduct2x4(v1, v2)
//		}
//	}
//
//	func BenchmarkDotSIMDArray2x(b *testing.B) {
//		v1 := [4]float32{2, 3, 4, 5}
//		v2 := [4]float32{3, 4, 5, 6}
//		v3 := [4]float32{2, 3, 4, 5}
//		v4 := [4]float32{3, 4, 5, 7}
//		for b.Loop() {
//			_ = DotProductSIMDArray(v1, v2)
//			_ = DotProductSIMDArray(v3, v4)
//		}
//	}
//
//	func BenchmarkDotSIMDArrayPointer(b *testing.B) {
//		v1a := [4]float32{2, 3, 4, 5}
//		v2a := [4]float32{3, 4, 5, 6}
//		v3a := [4]float32{2, 3, 4, 5}
//		v4a := [4]float32{3, 4, 5, 7}
//		v1 := &v1a
//		v2 := &v2a
//		v3 := &v3a
//		v4 := &v4a
//		for b.Loop() {
//			Out = DotProductSIMDArrayPointer(v1, v2)
//			Out += DotProductSIMDArrayPointer(v3, v4)
//		}
//		b.Log(Out)
//	}
//
//	func BenchmarkDotSIMDArrayPointer8(b *testing.B) {
//		v1a := [8]float32{2, 3, 4, 5, 3, 4, 5, 6}
//		v3a := [8]float32{2, 3, 4, 5, 3, 4, 5, 7}
//		//v1 := &v1a
//		//v3 := &v3a
//
//		for b.Loop() {
//			Out = DotProductSIMDArrayPointer8(v1a)
//			Out += DotProductSIMDArrayPointer8(v3a)
//		}
//		b.Log(Out)
//	}
//
//	func BenchmarkDotSIMDSlice(b *testing.B) {
//		v1 := []float32{2, 3, 4, 5}
//		v2 := []float32{3, 4, 5, 6}
//		v3 := []float32{2, 3, 4, 5}
//		v4 := []float32{3, 4, 5, 7}
//		for b.Loop() {
//			_ = DotProductSIMDSlice(v1, v2)
//			_ = DotProductSIMDSlice(v3, v4)
//		}
//	}
func BenchmarkDotSIMD2x4(b *testing.B) {
	v1 := []float32{2, 3, 4, 5, 2, 3, 4, 5}
	v2 := []float32{3, 4, 5, 6, 3, 4, 5, 6}
	tot := float32(0.0)
	for b.Loop() {
		a, b := DotProductSIMD2x4(v1, v2)
		tot = tot + a + b
	}
}

//
//func BenchmarkDotFraction(b *testing.B) {
//	v1 := []float32{2.7, 3.4, 4.2, 5.9}
//	v2 := []float32{3.8, 4.1, 5.3, 6.7}
//	v3 := []float32{2.7, 3.4, 4.2, 5.9}
//	v4 := []float32{3.8, 4.1, 5.3, 6.7}
//
//	for b.Loop() {
//		_ = DotProduct(v1, v2)
//		_ = DotProduct(v3, v4)
//	}
//}
//func BenchmarkDotSIMDSliceFraction(b *testing.B) {
//	v1 := []float32{2.7, 3.4, 4.2, 5.9}
//	v2 := []float32{3.8, 4.1, 5.3, 6.7}
//	v3 := []float32{2.7, 3.4, 4.2, 5.9}
//	v4 := []float32{3.8, 4.1, 5.3, 6.7}
//	for b.Loop() {
//		_ = DotProductSIMDSlice(v1, v2)
//		_ = DotProductSIMDSlice(v3, v4)
//	}
//}
//func BenchmarkDotSIMD2x4Fraction(b *testing.B) {
//	v1 := []float32{2.7, 3.4, 4.2, 5.9, 3.8, 4.1, 5.3, 6.7}
//	v2 := []float32{2.7, 3.4, 4.2, 5.9, 3.8, 4.1, 5.3, 6.7}
//
//	for b.Loop() {
//		_, _ = DotProductSIMD2x4(v1, v2)
//	}
//}

func BenchmarkDotAVX256a(b *testing.B) {
	v1 := []float32{2, 3, 4, 5}
	v2 := []float32{2, 3, 4, 5}
	tot := float32(0.0)
	for b.Loop() {
		tot += DotAVX256a(v1, v2)
	}
}

func BenchmarkDotGoSIMD(b *testing.B) {
	v1 := []float32{2, 3, 4, 5}
	v2 := []float32{2, 3, 4, 5}
	tot := float32(0.0)
	for b.Loop() {
		tot += DotGoSIMD(v1, v2)
	}
}
