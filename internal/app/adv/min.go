package adv

import "simd/archsimd"

func FindMinInt32x8(v archsimd.Int32x8) int32 {
	lo := v.GetLo()
	hi := v.GetHi()
	tmp := hi.Min(lo) // min( [8,2,6,4], [5,1,8,3]) => [5,1,6,3]

	// Next, we need to combine shuffling with more Min-calls, i.e. "rotating" the elements sideways, performing a new
	// min call per iteration so we after three "rotations" have checked all elements against each other and the resulting
	// min-value is present in all 4 elements.
	tmp = tmp.Min(tmp.SelectFromPair(1, 2, 3, 0, tmp)) // min( [5,1,6,3], [3,5,1,6]) => [3,1,1,3]
	tmp = tmp.Min(tmp.SelectFromPair(2, 3, 0, 1, tmp)) // min( [3,1,1,3], [3,3,1,1]) => [3,1,1,1]
	tmp = tmp.Min(tmp.SelectFromPair(3, 0, 1, 2, tmp)) // min( [3,1,1,1], [1,3,1,1]) => [1,1,1,1]

	return tmp.GetElem(0)
}
