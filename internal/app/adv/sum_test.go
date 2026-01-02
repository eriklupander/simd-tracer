package adv

import (
	"simd/archsimd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSumSIMD(t *testing.T) {
	vec := [8]int32{1, 2, 3, 4, 5, 6, 7, 8}
	assert.Equal(t, int32(36), SumInt32SIMD(&vec))
}

func TestSumMul(t *testing.T) {
	vec := [8]int32{1, 2, 3, 4, 5, 6, 7, 8}
	assert.Equal(t, int32(40320), SumMul(&vec))
}

func TestSumMulVertical(t *testing.T) {
	vec := [8]int32{1, 2, 3, 4, 5, 6, 7, 8}
	assert.Equal(t, int32(40320), SumMulVertical(&vec))
}

var S int32

func BenchmarkSumMul(b *testing.B) {
	v := &[8]int32{1, 2, 3, 4, 5, 6, 7, 8}
	for b.Loop() {
		S = SumMul(v)
	}
}
func BenchmarkSumMulVertical(b *testing.B) {
	v := &[8]int32{1, 2, 3, 4, 5, 6, 7, 8}
	e1 := archsimd.LoadInt32x8Slice([]int32{v[0], 0, 0, 0, 0, 0, 0, 0})
	e2 := archsimd.LoadInt32x8Slice([]int32{v[1], 0, 0, 0, 0, 0, 0, 0})
	e3 := archsimd.LoadInt32x8Slice([]int32{v[2], 0, 0, 0, 0, 0, 0, 0})
	e4 := archsimd.LoadInt32x8Slice([]int32{v[3], 0, 0, 0, 0, 0, 0, 0})
	e5 := archsimd.LoadInt32x8Slice([]int32{v[4], 0, 0, 0, 0, 0, 0, 0})
	e6 := archsimd.LoadInt32x8Slice([]int32{v[5], 0, 0, 0, 0, 0, 0, 0})
	e7 := archsimd.LoadInt32x8Slice([]int32{v[6], 0, 0, 0, 0, 0, 0, 0})
	e8 := archsimd.LoadInt32x8Slice([]int32{v[7], 0, 0, 0, 0, 0, 0, 0})
	for b.Loop() {
		S = SumMulVerticalSIMD(e1, e2, e3, e4, e5, e6, e7, e8)
	}
}

func BenchmarkSumMulVerticalArrays(b *testing.B) {

	e1 := &[8]int32{1, 0, 0, 0, 0, 0, 0, 0}
	e2 := &[8]int32{2, 0, 0, 0, 0, 0, 0, 0}
	e3 := &[8]int32{3, 0, 0, 0, 0, 0, 0, 0}
	e4 := &[8]int32{4, 0, 0, 0, 0, 0, 0, 0}
	e5 := &[8]int32{5, 0, 0, 0, 0, 0, 0, 0}
	e6 := &[8]int32{6, 0, 0, 0, 0, 0, 0, 0}
	e7 := &[8]int32{7, 0, 0, 0, 0, 0, 0, 0}
	e8 := &[8]int32{8, 0, 0, 0, 0, 0, 0, 0}
	for b.Loop() {
		S = SumMulVerticalSIMDArrays(e1, e2, e3, e4, e5, e6, e7, e8)
	}
}
func BenchmarkSumMul8(b *testing.B) {
	v := &[8]int32{1, 2, 3, 4, 5, 6, 7, 8}
	for b.Loop() {
		S = SumMul(v)
		S = SumMul(v)
		S = SumMul(v)
		S = SumMul(v)
		S = SumMul(v)
		S = SumMul(v)
		S = SumMul(v)
		S = SumMul(v)
	}
}
func BenchmarkSumMulVerticalArrays8(b *testing.B) {

	e1 := &[8]int32{1, 1, 1, 1, 1, 1, 1, 1}
	e2 := &[8]int32{2, 2, 2, 2, 2, 2, 2, 2}
	e3 := &[8]int32{3, 3, 3, 3, 3, 3, 3, 3}
	e4 := &[8]int32{4, 4, 4, 4, 4, 4, 4, 4}
	e5 := &[8]int32{5, 5, 5, 5, 5, 5, 5, 5}
	e6 := &[8]int32{6, 6, 6, 6, 6, 6, 6, 6}
	e7 := &[8]int32{7, 7, 7, 7, 7, 7, 7, 7}
	e8 := &[8]int32{8, 8, 8, 8, 8, 8, 8, 8}
	for b.Loop() {
		S = SumMulVerticalSIMDArrays(e1, e2, e3, e4, e5, e6, e7, e8)
	}
}
func BenchmarkSumMulVerticalArrays8FullResults(b *testing.B) {

	e1 := &[8]int32{1, 1, 1, 1, 1, 1, 1, 1}
	e2 := &[8]int32{2, 2, 2, 2, 2, 2, 2, 2}
	e3 := &[8]int32{3, 3, 3, 3, 3, 3, 3, 3}
	e4 := &[8]int32{4, 4, 4, 4, 4, 4, 4, 4}
	e5 := &[8]int32{5, 5, 5, 5, 5, 5, 5, 5}
	e6 := &[8]int32{6, 6, 6, 6, 6, 6, 6, 6}
	e7 := &[8]int32{7, 7, 7, 7, 7, 7, 7, 7}
	e8 := &[8]int32{8, 8, 8, 8, 8, 8, 8, 8}
	results := &[8]int32{}
	for b.Loop() {
		res := SumMulVerticalSIMDArraysFullResults(e1, e2, e3, e4, e5, e6, e7, e8)
		res.Store(results)
	}
}
