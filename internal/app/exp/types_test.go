package exp

import (
	"testing"

	"github.com/eriklupander/simd-tracer/internal/app/std"
)

var Sum float32

func BenchmarkFields(b *testing.B) {
	v1 := std.Vec3{X: 1.5, Y: 22.3, Z: 34.1}
	var sum float32
	for i := 0; i < b.N; i++ {
		sum = sum + v1.X*v1.X + v1.Y*v1.Y + v1.Z*v1.Z
	}
	Sum = sum
}

func BenchmarkArray(b *testing.B) {
	v1 := ArrayVec3{1.5, 22.3, 34.1}
	var sum float32
	for i := 0; i < b.N; i++ {
		sum = sum + v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2]
	}
	Sum = sum
}
