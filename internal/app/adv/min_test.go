package adv

import (
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
