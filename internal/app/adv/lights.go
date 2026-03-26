//go:build goexperiment.simd && amd64

package adv

import (
	"math"
	"math/rand/v2"

	"github.com/eriklupander/simd-tracer/internal/app/std"
)

// RandomPointOnHemisphere returns a uniformly distributed random point on the
// surface of sphere s, restricted to the hemisphere that faces toward normal.
// normal should be the unit vector pointing from the sphere center toward the
// shading point (i.e. the hemisphere of the sphere visible from that point).
// RandomPointOnHemisphere writes a uniformly distributed random point on the
// surface of sphere s into out, restricted to the hemisphere that faces toward normal.
// normal should be the unit vector pointing from the sphere center toward the
// shading point (i.e. the hemisphere of the sphere visible from that point).
func RandomPointOnHemisphere(s std.Sphere, normal std.Vec4, out *std.Vec4) {
	// Uniform hemisphere sampling: u1 = cos(theta), u2 drives phi.
	u1 := rand.Float32()
	u2 := rand.Float32()

	sinTheta := float32(math.Sqrt(float64(1 - u1*u1)))
	phi := 2 * math.Pi * float64(u2)

	// Local-space sample (hemisphere aligned with +z).
	lx := sinTheta * float32(math.Cos(phi))
	ly := sinTheta * float32(math.Sin(phi))
	lz := u1 // cos(theta) ≥ 0, so always in upper hemisphere

	// Build an orthonormal basis (tangent, bitangent, normal).
	var up std.Vec4
	if math.Abs(float64(normal[0])) < 0.9 {
		up = std.Vec4{1, 0, 0, 0}
	} else {
		up = std.Vec4{0, 1, 0, 0}
	}

	var tangent, bitangent std.Vec4
	up.CrossProduct(&normal, &tangent)
	tangent.Normalize()
	normal.CrossProduct(&tangent, &bitangent)

	// Rotate local sample into world space and scale to sphere surface.
	out[0] = s.Center[0] + (tangent[0]*lx+bitangent[0]*ly+normal[0]*lz)*s.Radius
	out[1] = s.Center[1] + (tangent[1]*lx+bitangent[1]*ly+normal[1]*lz)*s.Radius
	out[2] = s.Center[2] + (tangent[2]*lx+bitangent[2]*ly+normal[2]*lz)*s.Radius
	out[3] = 0
}
