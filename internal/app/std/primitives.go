//go:build goexperiment.simd && amd64

package std

import (
	"math"
)

type Vec4 [4]float32

func NewVec4() Vec4 {
	out := Vec4{0, 0, 0, 0}
	return out
}

func (v1 *Vec4) normalize() {
	magnitude := float32(math.Sqrt(float64(v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2])))
	v1[0] = v1[0] / magnitude
	v1[1] = v1[1] / magnitude
	v1[2] = v1[2] / magnitude
}

func (v1 *Vec4) dotProduct(v2 *Vec4) float32 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

func (v1 *Vec4) crossProduct(t2 *Vec4, dst *Vec4) {

	dst[0] = v1[1]*t2[2] - v1[2]*t2[1]
	dst[1] = v1[2]*t2[0] - v1[0]*t2[2]
	dst[2] = v1[0]*t2[1] - v1[1]*t2[0]

	/*
		c[0] = a[1]*b[2] - a[2]*b[1]
		c[1] = a[2]*b[0] - a[0]*b[2]
		c[2] = a[0]*b[1] - a[1]*b[0]
	*/
}

func (v1 *Vec4) mul(f float32) {
	v1[0] = v1[0] * f
	v1[1] = v1[1] * f
	v1[2] = v1[2] * f
}

func (v1 *Vec4) negate() {
	v1[0] = -v1[0]
	v1[1] = -v1[1]
	v1[2] = -v1[2]
}

type Ray struct {
	Orig *Vec4
	Dir  *Vec4
}

//	func addMul(v1 Vec4, v2 Vec4, mulBy float32, dst Vec4) {
//		dst[0] = v1[0] + v2[0]*mulBy
//		dst[1] = v1[1] + v2[1]*mulBy
//		dst[2] = v1[2] + v2[2]*mulBy
//	}
func addMulVec4(v1 *Vec4, v2 *Vec4, mulBy float32, dst *Vec4) {

	dst[0] = v1[0] + v2[0]*mulBy
	dst[1] = v1[1] + v2[1]*mulBy
	dst[2] = v1[2] + v2[2]*mulBy
}

type Matrix44 [4][4]float32
type Matrix33 [3][3]float32
type Matrix22 [2][2]float32

func NewIdentityMatrix44() Matrix44 {
	return NewMatrix44f(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)
}

func NewMatrix44f(args ...float32) Matrix44 {
	x := Matrix44{}
	if len(args) == 0 {
		args = make([]float32, 16)
		for i := range 16 {
			args[i] = 0.0
		}
	}

	x[0][0] = args[0]
	x[0][1] = args[1]
	x[0][2] = args[2]
	x[0][3] = args[3]
	x[1][0] = args[4]
	x[1][1] = args[5]
	x[1][2] = args[6]
	x[1][3] = args[7]
	x[2][0] = args[8]
	x[2][1] = args[9]
	x[2][2] = args[10]
	x[2][3] = args[11]
	x[3][0] = args[12]
	x[3][1] = args[13]
	x[3][2] = args[14]
	x[3][3] = args[15]
	return x
}

func (x *Matrix44) indexed(i int) float32 {
	return x[i/4][i%4]
}

//
//func (x *Matrix44) mulVec(src Vec4f) (dst Vec4f) {
//
//	var a, b, c, w float32
//
//	a = src[0]*x[0][0] + src[1]*x[1][0] + src[2]*x[2][0] + x[3][0]
//	b = src[0]*x[0][1] + src[1]*x[1][1] + src[2]*x[2][1] + x[3][1]
//	c = src[0]*x[0][2] + src[1]*x[1][2] + src[2]*x[2][2] + x[3][2]
//	w = src[0]*x[0][3] + src[1]*x[1][3] + src[2]*x[2][3] + x[3][3]
//
//	dst = &Vec3f{a / w, b / w, c / w}
//
//	return dst
//}

func (x *Matrix44) mulVec2(src Vec4) Vec4 {

	var a, b, c, w float32

	a = src[0]*x[0][0] + src[1]*x[1][0] + src[2]*x[2][0] + x[3][0]
	b = src[0]*x[0][1] + src[1]*x[1][1] + src[2]*x[2][1] + x[3][1]
	c = src[0]*x[0][2] + src[1]*x[1][2] + src[2]*x[2][2] + x[3][2]
	w = src[0]*x[0][3] + src[1]*x[1][3] + src[2]*x[2][3] + x[3][3]

	dst := NewVec4()
	dst[0] = a / w
	dst[1] = b / w
	dst[2] = c / w
	//	dst = &Vec3f{a / w, b / w, c / w}

	return dst
}

//
//func (x *Matrix44) mulDir(src Vec4f) (dst Vec4f) {
//
//	var a, b, c float32
//
//	a = src[0]*x[0][0] + src[1]*x[1][0] + src[2]*x[2][0]
//	b = src[0]*x[0][1] + src[1]*x[1][1] + src[2]*x[2][1]
//	c = src[0]*x[0][2] + src[1]*x[1][2] + src[2]*x[2][2]
//
//	dst = &Vec3f{a, b, c}
//	return dst
//}

func (x *Matrix44) mulDir2(src Vec4) Vec4 {
	dst := NewVec4()
	dst[0] = src[0]*x[0][0] + src[1]*x[1][0] + src[2]*x[2][0]
	dst[1] = src[0]*x[0][1] + src[1]*x[1][1] + src[2]*x[2][1]
	dst[2] = src[0]*x[0][2] + src[1]*x[1][2] + src[2]*x[2][2]
	return dst
}

func (x *Matrix44) mulDirScalar(a float32, b float32, c float32, dst *Vec4) {

	dst[0] = a*x[0][0] + b*x[1][0] + c*x[2][0]
	dst[1] = a*x[0][1] + b*x[1][1] + c*x[2][1]
	dst[2] = a*x[0][2] + b*x[1][2] + c*x[2][2]

}

//	func sub(a, b Vec4f) Vec4f {
//		return &Vec3f{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
//	}
func sub2(a, b *Vec4, dst *Vec4) {
	dst[0] = a[0] - b[0]
	dst[1] = a[1] - b[1]
	dst[2] = a[2] - b[2]
	dst[3] = a[3] - b[3]
}

func multiplyMatricies(m1 Matrix44, m2 Matrix44) Matrix44 {
	m3 := NewMatrix44f()
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			m3[row][col] = multiply4x4(m1, m2, row, col)
		}
	}
	return m3
}

func newTranslationMatrix(x, y, z float32) Matrix44 {

	m1 := NewIdentityMatrix44() //NewMat4x4(make([]float64, 16))
	//copy(m1.Elems, IdentityMatrix.Elems)
	m1[0][3] = x
	m1[1][3] = y
	m1[2][3] = z
	return m1

}

func multiply4x4(m1 Matrix44, m2 Matrix44, row int, col int) float32 {
	// always row from m1, col from m2
	a0 := m1[row][0] * m2[0][col]
	a1 := m1[row][1] * m2[1][col]
	a2 := m1[row][2] * m2[2][col]
	a3 := m1[row][3] * m2[3][col]
	return a0 + a1 + a2 + a3
}

func Transpose(m1 Matrix44) Matrix44 {
	m3 := NewMatrix44f()
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			m3[row][col] = m1[col][row]
		}
	}
	return m3
}

// Determinant2x2 returns A-D minus B-C for a
// [A,B]
// [C,D]
// Matrix.
func Determinant2x2(m1 Matrix22) float32 {
	return m1[0][0]*m1[1][1] - m1[0][1]*m1[1][0]
}

// Determinant3x3 takes the first row of the passed matrix, summing the colvalue * Cofactor of the same col
func Determinant3x3(m1 Matrix33) float32 {
	det := float32(0.0)
	for col := 0; col < 3; col++ {
		det = det + m1[0][col]*Cofactor3x3(m1, 0, col)
	}
	return det
}

// Determinant4x4
func Determinant4x4(m1 Matrix44) float32 {
	det := float32(0.0)
	for col := 0; col < 4; col++ {
		det = det + m1[0][col]*Cofactor4x4(m1, 0, col)
	}
	return det
}

// Submatrix3x3 extracts the 2x2 submatrix after deleting row and col from the passed 3x3
func Submatrix3x3(m1 Matrix33, deleteRow, deleteCol int) Matrix22 {
	m3 := Matrix22{
		[2]float32{},
		[2]float32{},
	}
	idxRow := 0
	idxCol := 0
	for row := 0; row < 3; row++ {
		if row == deleteRow {
			continue
		}

		for col := 0; col < 3; col++ {
			if col == deleteCol {
				continue
			}

			m3[idxRow][idxCol] = m1[row][col] //.Get(row, col)
			idxCol++
		}
		idxRow++
		idxCol = 0
	}
	return m3
}

// Submatrix4x4 extracts the 3x3 submatrix after deleting row and col from the passed 4x4
func Submatrix4x4(m1 Matrix44, deleteRow, deleteCol int) Matrix33 {
	m3 := Matrix33{
		[3]float32{},
		[3]float32{},
		[3]float32{},
	}
	idxRow := 0
	idxCol := 0
	for row := 0; row < 4; row++ {
		if row == deleteRow {
			continue
		}
		for col := 0; col < 4; col++ {
			if col == deleteCol {
				continue
			}
			m3[idxRow][idxCol] = m1[row][col] //.Get(row, col)
			idxCol++
		}
		idxRow++
		idxCol = 0
	}
	return m3
}

// Minor3x3 computes the submatrix at row/col and returns the determinant of the computed matrix.
func Minor3x3(m1 Matrix33, row, col int) float32 {
	m2 := Submatrix3x3(m1, row, col)
	return Determinant2x2(m2)
}

// Minor4x4 computes the submatrix at row/col and returns the determinant of the computed matrix.
func Minor4x4(m1 Matrix44, row, col int) float32 {
	m2 := Submatrix4x4(m1, row, col)
	return Determinant3x3(m2)
}

// Cofactor3x3 may change the sign of the computed minor of the passed matrix
func Cofactor3x3(m1 Matrix33, row, col int) float32 {
	minor := Minor3x3(m1, row, col)
	if (row+col)%2 != 0 {
		return -minor
	}
	return minor
}

// Cofactor4x4 may change the sign of the computed minor of the passed matrix
func Cofactor4x4(m1 Matrix44, row, col int) float32 {
	minor := Minor4x4(m1, row, col)
	if (row+col)%2 != 0 {
		return -minor
	}
	return minor
}

func Inverse(m1 Matrix44) Matrix44 {
	m3 := NewMatrix44f()
	d4 := Determinant4x4(m1)
	for row := 0; row < 4; row++ {

		for col := 0; col < 4; col++ {
			c := Cofactor4x4(m1, row, col)
			// note that "col, row" here, instead of "row, col",
			// accomplishes the transpose operation!
			m3[col][row] = c / d4
		}
	}
	return m3
}
