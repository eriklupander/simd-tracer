package adv

import (
	"math"
	"simd/archsimd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	v1 := archsimd.LoadInt32x8Slice([]int32{1, 2, 3, 4, 5, 6, 7, 8})
	v2 := archsimd.LoadInt32x8Slice([]int32{8, 7, 6, 5, 4, 3, 2, 1})
	minElems := v1.Min(v2)
	assert.Equal(t, int32(1), minElems.GetLo().GetElem(0))
	assert.Equal(t, int32(2), minElems.GetLo().GetElem(1))
}

func TestMinInVec(t *testing.T) {
	v := archsimd.LoadInt32x8Slice([]int32{8, 2, 6, 4, 5, 1, 8, 3})
	minVal := FindMinInt32x8(v)
	assert.Equal(t, int32(1), minVal)
}

func TestFindMinRegular(t *testing.T) {
	v := []int32{7, 2, 34, 44, 5, 16, 6, 3}
	m := min(v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7])
	assert.Equal(t, int32(2), m)
}

func TestFindMinWithForLoop(t *testing.T) {
	v := []int32{7, 2, 34, 44, 5, 16, 6, 3}
	m := int32(math.MaxInt32)
	for _, n := range v {
		if n < m {
			m = n
		}
	}
	assert.Equal(t, int32(2), m)
}

var M int32

func BenchmarkWithForLoop(b *testing.B) {
	v := []int32{7, 2, 34, 44, 5, 16, 6, 3}
	m := int32(math.MaxInt32)
	for b.Loop() {
		m = min(v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7])
	}
	M = m
}

func BenchmarkFindMinRegular(b *testing.B) {
	v := []int32{7, 2, 34, 44, 5, 16, 6, 3}
	m := int32(math.MaxInt32)
	for b.Loop() {
		for _, n := range v {
			if n < m {
				m = n
			}
		}
	}
	M = m
}
