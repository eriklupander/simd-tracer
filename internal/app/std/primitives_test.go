package std

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var Epsilon = 0.001

func TestInverse(t *testing.T) {
	m1 := NewMatrix44f(-5, 2, 6, -8, 1, -5, 1, 8, 7, 7, -6, -7, 1, -3, 7, 4)
	m3 := Inverse(m1)

	cf1 := Cofactor4x4(m1, 2, 3)
	cf2 := Cofactor4x4(m1, 3, 2)
	determinant := Determinant4x4(m1)
	assert.EqualValues(t, 532.0, determinant)
	assert.EqualValues(t, -160.0, cf1)
	assert.EqualValues(t, 105.0, cf2)

	expected := NewMatrix44f(0.21805, 0.45113, 0.24060, -0.04511,
		-0.80827, -1.45677, -0.44361, 0.52068,
		-0.07895, -0.22368, -0.05263, 0.19737,
		-0.52256, -0.81391, -0.30075, 0.30639)

	for i := 0; i < 16; i++ {
		assert.InEpsilon(t, expected.indexed(i), m3.indexed(i), Epsilon, fmt.Sprintf("index %d failed: values: %v %v", i, expected.indexed(i), m3.indexed(i)))
	}
}

func TestInverse2(t *testing.T) {
	m1 := NewMatrix44f(8, -5, 9, 2, 7, 5, 6, 1, -6, 0, 9, 6, -3, 0, -9, -4)
	m3 := Inverse(m1)
	expected := NewMatrix44f(-0.15385, -0.15385, -0.28205, -0.53846,
		-0.07692, 0.12308, 0.02564, 0.03077,
		0.35897, 0.35897, 0.43590, 0.92308,
		-0.69231, -0.69231, -0.76923, -1.92308)

	for i := 0; i < 16; i++ {
		assert.InEpsilon(t, expected.indexed(i), m3.indexed(i), Epsilon, fmt.Sprintf("index %d failed: values: %v %v", i, expected.indexed(i), m3.indexed(i)))
	}
}
func TestInverse3(t *testing.T) {
	m1 := NewMatrix44f(
		9, 3, 0, 9,
		-5, -2, -6, -3,
		-4, 9, 6, 4,
		-7, 6, 6, 2)
	m3 := Inverse(m1)

	expected := NewMatrix44f(-0.04074, -0.07778, 0.14444, -0.22222,
		-0.07778, 0.03333, 0.36667, -0.33333,
		-0.02901, -0.14630, -0.10926, 0.12963,
		0.17778, 0.06667, -0.26667, 0.33333)

	for i := 0; i < 16; i++ {
		assert.InEpsilon(t, expected.indexed(i), m3.indexed(i), Epsilon, fmt.Sprintf("index %d failed: values: %v %v", i, expected.indexed(i), m3.indexed(i)))
	}
}
