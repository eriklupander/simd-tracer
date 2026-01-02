package adv

import (
	"simd/archsimd"
	"slices"
	"testing"

	"github.com/eriklupander/simd-tracer/internal/app/std"
	"github.com/stretchr/testify/assert"
)

func TestDots(t *testing.T) {
	// 8 logical 4-element vectors represented column-wise
	x := archsimd.LoadFloat32x8Slice([]float32{1, 2, 11, 12, 22, 44, 66, 88})
	y := archsimd.LoadFloat32x8Slice([]float32{3, 4, 13, 14, 22, 44, 66, 88})
	z := archsimd.LoadFloat32x8Slice([]float32{5, 6, 15, 16, 22, 44, 66, 88})
	w := archsimd.LoadFloat32x8Slice([]float32{7, 8, 17, 18, 22, 44, 66, 88})

	// 8 other logical 4-element vectors represented column-wise
	x2 := archsimd.LoadFloat32x8Slice([]float32{11, 12, 111, 112, 122, 144, 166, 188})
	y2 := archsimd.LoadFloat32x8Slice([]float32{13, 14, 113, 114, 122, 144, 166, 188})
	z2 := archsimd.LoadFloat32x8Slice([]float32{15, 16, 115, 116, 122, 144, 166, 188})
	w2 := archsimd.LoadFloat32x8Slice([]float32{17, 18, 117, 118, 122, 144, 166, 188})

	dotProducts := x.Mul(x2).
		Add(y.Mul(y2)).
		Add(z.Mul(z2)).
		Add(w.Mul(w2))

	res := [8]float32{}
	dotProducts.Store(&res)
	assert.True(t, slices.Equal(res[:], []float32{244, 320, 6404, 6920, 10736, 25344, 43824, 66176}))
}

func BenchmarkDotPlainGoVertical(b *testing.B) {
	x, y, z, w, x2, y2, z2, w2 := testVectors()
	dotProducts := make([]float32, 8)
	tot := float32(0.0)
	for b.Loop() {
		std.DotProduct8(x, y, z, w, x2, y2, z2, w2, dotProducts)
		tot += dotProducts[0] + dotProducts[1] + dotProducts[2] + dotProducts[3] + dotProducts[4] + dotProducts[5] + dotProducts[6] + dotProducts[7]
	}
}

func BenchmarkDotPlainGoHorizontal(b *testing.B) {
	v1, v2, v3, v4, v5, v6, v7, v8 := testVectorsHorizontal()
	tot := float32(0.0)
	for b.Loop() {
		a, b := std.DotProduct2x4(v1, v2)
		c, d := std.DotProduct2x4(v3, v4)
		e, f := std.DotProduct2x4(v5, v6)
		g, h := std.DotProduct2x4(v7, v8)
		tot += a + b + c + d + e + f + g + h
	}
}

func BenchmarkDotSIMDHorizontal(b *testing.B) {

	v1, v2, v3, v4, v5, v6, v7, v8 := testVectorsHorizontal()

	v1Simd := archsimd.LoadFloat32x8Slice(v1)
	v2Simd := archsimd.LoadFloat32x8Slice(v2)
	v3Simd := archsimd.LoadFloat32x8Slice(v3)
	v4Simd := archsimd.LoadFloat32x8Slice(v4)
	v5Simd := archsimd.LoadFloat32x8Slice(v5)
	v6Simd := archsimd.LoadFloat32x8Slice(v6)
	v7Simd := archsimd.LoadFloat32x8Slice(v7)
	v8Simd := archsimd.LoadFloat32x8Slice(v8)

	tot := float32(0.0)
	for b.Loop() {
		a, b := std.DotProductSIMD2x4FromFloat32x8(v1Simd, v2Simd)
		c, d := std.DotProductSIMD2x4FromFloat32x8(v3Simd, v4Simd)
		e, f := std.DotProductSIMD2x4FromFloat32x8(v5Simd, v6Simd)
		g, h := std.DotProductSIMD2x4FromFloat32x8(v7Simd, v8Simd)
		tot += a + b + c + d + e + f + g + h
	}
}

func BenchmarkDotSIMDVertical(b *testing.B) {

	x1, y1, z1, w1, x2, y2, z2, w2 := testSIMDVectors()

	sums := [8]float32{}
	tot := float32(0)
	for b.Loop() {
		dots := DotProductSIMD8(x1, y1, z1, w1, x2, y2, z2, w2)
		dots.Store(&sums)
		tot += sums[0] + sums[1] + sums[2] + sums[3] + sums[4] + sums[5] + sums[6] + sums[7]
	}
}

func BenchmarkDotSIMDMulAddVertical(b *testing.B) {

	x1, y1, z1, w1, x2, y2, z2, w2 := testSIMDVectors()

	sums := [8]float32{}
	tot := float32(0)
	for b.Loop() {
		dots := DotProductSIMD8MulAdd(x1, y1, z1, w1, x2, y2, z2, w2)
		dots.Store(&sums)
		tot += sums[0] + sums[1] + sums[2] + sums[3] + sums[4] + sums[5] + sums[6] + sums[7]
	}
}

func TestDotProduct8MulAdd(t *testing.T) {

	x1, y1, z1, w1, x2, y2, z2, w2 := testSIMDVectors()
	sum := DotProductSIMD8MulAdd(x1, y1, z1, w1, x2, y2, z2, w2)
	assert.Equal(t, float32(244.0), sum.GetLo().GetElem(0))
}

func testVectors() ([]float32, []float32, []float32, []float32, []float32, []float32, []float32, []float32) {

	x := []float32{1, 2, 11, 12, 22, 44, 66, 88}          // X elements
	y := []float32{3, 4, 13, 14, 22, 44, 66, 88}          // Y elements
	z := []float32{5, 6, 15, 16, 22, 44, 66, 88}          // Z elements
	w := []float32{7, 8, 17, 18, 22, 44, 66, 88}          // W elements
	x2 := []float32{11, 12, 111, 112, 122, 144, 166, 188} // X elements
	y2 := []float32{13, 14, 113, 114, 122, 144, 166, 188} // Y elements
	z2 := []float32{15, 16, 115, 116, 122, 144, 166, 188} // Z elements
	w2 := []float32{17, 18, 117, 118, 122, 144, 166, 188} // W elements

	return x, y, z, w, x2, y2, z2, w2
}

func testVectorsHorizontal() ([]float32, []float32, []float32, []float32, []float32, []float32, []float32, []float32) {

	x, y, z, w, x2, y2, z2, w2 := testVectors()
	n := 0
	v1 := []float32{x[n], y[n], z[n], w[n], x2[n], y2[n], z2[n], w2[n]}
	n++
	v2 := []float32{x[n], y[n], z[n], w[n], x2[n], y2[n], z2[n], w2[n]}
	n++
	v3 := []float32{x[n], y[n], z[n], w[n], x2[n], y2[n], z2[n], w2[n]}
	n++
	v4 := []float32{x[n], y[n], z[n], w[n], x2[n], y2[n], z2[n], w2[n]}
	n++
	v5 := []float32{x[n], y[n], z[n], w[n], x2[n], y2[n], z2[n], w2[n]}
	n++
	v6 := []float32{x[n], y[n], z[n], w[n], x2[n], y2[n], z2[n], w2[n]}
	n++
	v7 := []float32{x[n], y[n], z[n], w[n], x2[n], y2[n], z2[n], w2[n]}
	n++
	v8 := []float32{x[n], y[n], z[n], w[n], x2[n], y2[n], z2[n], w2[n]}

	return v1, v2, v3, v4, v5, v6, v7, v8
}

func testSIMDVectors() (archsimd.Float32x8, archsimd.Float32x8, archsimd.Float32x8, archsimd.Float32x8, archsimd.Float32x8, archsimd.Float32x8, archsimd.Float32x8, archsimd.Float32x8) {
	simdX, simdY, simdZ, simdW, simdX2, simdY2, simdZ2, simdW2 := testVectors()

	x1 := archsimd.LoadFloat32x8Slice(simdX)
	y1 := archsimd.LoadFloat32x8Slice(simdY)
	z1 := archsimd.LoadFloat32x8Slice(simdZ)
	w1 := archsimd.LoadFloat32x8Slice(simdW)

	x2 := archsimd.LoadFloat32x8Slice(simdX2)
	y2 := archsimd.LoadFloat32x8Slice(simdY2)
	z2 := archsimd.LoadFloat32x8Slice(simdZ2)
	w2 := archsimd.LoadFloat32x8Slice(simdW2)

	return x1, y1, z1, w1, x2, y2, z2, w2
}
