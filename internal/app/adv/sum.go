package adv

import "simd/archsimd"

func SumInt32SIMD(v *[8]int32) int32 {
	elems := archsimd.LoadInt32x8(v) // [1,2,3,4,5,6,7,8]
	lo := elems.GetLo()              // [1,2,3,4]
	hi := elems.GetHi()              // [5,6,7,8]

	lo = lo.AddPairs(lo) // [3,7,3,7]
	lo = lo.AddPairs(lo) // [10,10,10,10]

	hi = hi.AddPairs(hi) // [11,15,11,15]
	hi = hi.AddPairs(hi) // [26,26,26,26]

	return lo.GetElem(0) + hi.GetElem(0)
}

func SumMul(v *[8]int32) int32 {
	elems := archsimd.LoadInt32x8(v)                                         // [1,2,3,4,5,6,7,8]
	lo := elems.GetLo()                                                      // [1,2,3,4] == 24
	hi := elems.GetHi()                                                      // [5,6,7,8]
	tmp := lo.Mul(hi)                                                        // 1*5, 2*6, 3*7, 4*8
	return tmp.GetElem(0) * tmp.GetElem(1) * tmp.GetElem(2) * tmp.GetElem(3) // 6 * 12 * 21 * 32
}

func SumMulVertical(v *[8]int32) int32 {
	e1 := archsimd.LoadInt32x8Slice([]int32{v[0], 0, 0, 0, 0, 0, 0, 0})
	e2 := archsimd.LoadInt32x8Slice([]int32{v[1], 0, 0, 0, 0, 0, 0, 0})
	e3 := archsimd.LoadInt32x8Slice([]int32{v[2], 0, 0, 0, 0, 0, 0, 0})
	e4 := archsimd.LoadInt32x8Slice([]int32{v[3], 0, 0, 0, 0, 0, 0, 0})
	e5 := archsimd.LoadInt32x8Slice([]int32{v[4], 0, 0, 0, 0, 0, 0, 0})
	e6 := archsimd.LoadInt32x8Slice([]int32{v[5], 0, 0, 0, 0, 0, 0, 0})
	e7 := archsimd.LoadInt32x8Slice([]int32{v[6], 0, 0, 0, 0, 0, 0, 0})
	e8 := archsimd.LoadInt32x8Slice([]int32{v[7], 0, 0, 0, 0, 0, 0, 0})

	return e1.Mul(e2).Mul(e3).Mul(e4).
		Mul(e5).Mul(e6).Mul(e7).Mul(e8).
		GetLo().
		GetElem(0)
}

func SumMulVerticalSIMD(e1, e2, e3, e4, e5, e6, e7, e8 archsimd.Int32x8) int32 {
	return e1.Mul(e2).Mul(e3).Mul(e4).
		Mul(e5).Mul(e6).Mul(e7).Mul(e8).
		GetLo().
		GetElem(0)
}
func SumMulVerticalSIMDArrays(e1, e2, e3, e4, e5, e6, e7, e8 *[8]int32) int32 {
	v1 := archsimd.LoadInt32x8(e1)
	v2 := archsimd.LoadInt32x8(e2)
	v3 := archsimd.LoadInt32x8(e3)
	v4 := archsimd.LoadInt32x8(e4)
	v5 := archsimd.LoadInt32x8(e5)
	v6 := archsimd.LoadInt32x8(e6)
	v7 := archsimd.LoadInt32x8(e7)
	v8 := archsimd.LoadInt32x8(e8)
	return v1.Mul(v2).Mul(v3).Mul(v4).
		Mul(v5).Mul(v6).Mul(v7).Mul(v8).
		GetLo().
		GetElem(0)
}
func SumMulVerticalSIMDArraysFullResults(e1, e2, e3, e4, e5, e6, e7, e8 *[8]int32) archsimd.Int32x8 {
	v1 := archsimd.LoadInt32x8(e1)
	v2 := archsimd.LoadInt32x8(e2)
	v3 := archsimd.LoadInt32x8(e3)
	v4 := archsimd.LoadInt32x8(e4)
	v5 := archsimd.LoadInt32x8(e5)
	v6 := archsimd.LoadInt32x8(e6)
	v7 := archsimd.LoadInt32x8(e7)
	v8 := archsimd.LoadInt32x8(e8)
	return v1.Mul(v2).Mul(v3).Mul(v4).
		Mul(v5).Mul(v6).Mul(v7).Mul(v8)
}
