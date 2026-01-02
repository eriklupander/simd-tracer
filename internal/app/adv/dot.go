package adv

import "simd/archsimd"

func DotProductSIMD8(x1, y1, z1, w1, x2, y2, z2, w2 archsimd.Float32x8) archsimd.Float32x8 {
	return x1.Mul(x2).
		Add(y1.Mul(y2)).
		Add(z1.Mul(z2)).
		Add(w1.Mul(w2))
}

func DotProductSIMD8MulAdd(x1, y1, z1, w1, x2, y2, z2, w2 archsimd.Float32x8) archsimd.Float32x8 {
	sum := archsimd.Float32x8{}
	sum = x1.MulAdd(x2, sum)
	sum = y1.MulAdd(y2, sum)
	sum = z1.MulAdd(z2, sum)
	return w1.MulAdd(w2, sum)
}
