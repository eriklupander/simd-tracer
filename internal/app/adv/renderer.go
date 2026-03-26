//go:build goexperiment.simd && amd64

package adv

import (
	"bytes"
	"math"

	"github.com/eriklupander/simd-tracer/internal/app/std"
)

var obstructedColor = std.Vec4{0, 0, 0}

const (
	IntersectTypeNone = iota
	IntersectTypeSphere
	IntersectTypePlane
	IntersectTypeTriangle
)

// hard-coded make-believe sphere representing the point light.
var light = std.Sphere{
	Center:        &std.Vec4{2, 5.5, 0},
	Radius:        2,
	RadiusSquared: 4,
	Color:         std.Vec4{1, 1, 1},
}

// hitResult holds the outcome of closest-intersection testing for a single ray.
// idx is the raw index within the primitive's own slice — no offset encoding.
type hitResult struct {
	t    float32
	idx  int
	kind int     // IntersectType*
	u, v float32 // barycentric coords, used for triangles
}

func Render(width, height int, spheres *Spheres, planes *Planes, triangles *Triangles, lightSamples int, renderTo *bytes.Buffer) {
	hitToLightDir := &std.Vec4{2, 5.4, 0}

	cameraOrigin := &std.Vec4{3, 4, 20}
	lookAt := &std.Vec4{3, 4, 0}
	up := &std.Vec4{0, 1, 0}
	vt := std.ViewTransform(cameraOrigin, lookAt, up)
	cameraToWorld := std.Inverse(vt)
	cameraToWorld = std.Transpose(cameraToWorld)

	fov := 51.52
	scale := float32(math.Tan(deg2rad(fov * 0.5)))
	imageAspectRatio := float32(width) / float32(height)

	var orig = cameraToWorld.MulVec(std.Vec4{0, 0, 0})
	var ray = &std.Ray{Orig: &orig}
	var dir = &std.Vec4{}
	var shadowRay = &std.Ray{Dir: &std.Vec4{}, Orig: &std.Vec4{}}
	var hitPoint = &std.Vec4{}
	var hitNormal = &std.Vec4{}

	for j := 0; j < height; j++ {

		y := (1 - 2*(float32(j)+0.5)/float32(height)) * scale

		for i := 0; i < width; i++ {

			x := (2*(float32(i)+0.5)/float32(width) - 1) * imageAspectRatio * scale

			cameraToWorld.MulDirScalar(x, y, -1, dir)
			dir.Normalize()
			ray.Dir = dir

			hit := findClosestHit(ray, spheres, planes, triangles)
			if hit.kind == IntersectTypeNone {
				renderTo.Write([]byte{0x00, 0x00, 0x00, 0xFF})
				continue
			}

			std.AddMulVec4(ray.Orig, ray.Dir, hit.t, hitPoint)

			// Surface normal is needed before the shadow test for spheres (N·L self-shadow check).
			if hit.kind == IntersectTypeSphere {
				std.Sub(hitPoint, &std.Vec4{spheres.CenterX[hit.idx], spheres.CenterY[hit.idx], spheres.CenterZ[hit.idx], 0}, hitNormal)
				hitNormal.Normalize()
			}

			computeShadowRay(hitToLightDir, hitPoint, shadowRay)

			var color std.Vec4
			var fract float32
			if isObstructed(hit, shadowRay, hitNormal, spheres) {
				color = obstructedColor
			} else {
				color, fract = shade(hit, ray, hitNormal, spheres, planes)
			}

			finalFract := fract * 255
			renderTo.Write([]byte{byte(color[0] * finalFract), byte(color[1] * finalFract), byte(color[2] * finalFract), 0xFF})
		}
	}
}

// findClosestHit runs all three intersection tests and returns the nearest hit.
func findClosestHit(ray *std.Ray, spheres *Spheres, planes *Planes, triangles *Triangles) hitResult {
	result := hitResult{t: math.MaxFloat32, kind: IntersectTypeNone}

	if t, u, v, idx, hit := IntersectTrianglesSIMD(ray, triangles); hit {
		result = hitResult{t: t, idx: idx, kind: IntersectTypeTriangle, u: u, v: v}
	}
	if t, idx, hit := IntersectPlanesSIMD(ray, planes); hit && t < result.t {
		result = hitResult{t: t, idx: idx, kind: IntersectTypePlane}
	}
	if t, idx, hit := IntersectSpheresSIMD(ray, spheres); hit && t < result.t {
		result = hitResult{t: t, idx: idx, kind: IntersectTypeSphere}
	}

	return result
}

// isObstructed returns true if the path to the light is blocked.
// For sphere hits, a N·L ≤ 0 check short-circuits the more expensive SIMD shadow ray test.
func isObstructed(hit hitResult, shadowRay *std.Ray, hitNormal *std.Vec4, spheres *Spheres) bool {
	if hit.kind == IntersectTypeSphere && hitNormal.DotProduct(shadowRay.Dir) <= 0 {
		return true
	}
	return IntersectSpheresSIMDShadowRay(shadowRay, spheres)
}

// shade computes the surface color and brightness fraction for an unobstructed hit.
func shade(hit hitResult, ray *std.Ray, hitNormal *std.Vec4, spheres *Spheres, planes *Planes) (std.Vec4, float32) {
	switch hit.kind {
	case IntersectTypeSphere:
		fract := max(0.0, hitNormal[0]*-ray.Dir[0]+hitNormal[1]*-ray.Dir[1]+hitNormal[2]*-ray.Dir[2])
		return spheres.Color[hit.idx], fract

	case IntersectTypePlane:
		fract := max(0.0, planes.NormalX[hit.idx]*-ray.Dir[0]+planes.NormalY[hit.idx]*-ray.Dir[1]+planes.NormalZ[hit.idx]*-ray.Dir[2])
		return planes.Color[hit.idx], fract

	case IntersectTypeTriangle:
		return std.Vec4{hit.u, hit.v, hit.u * hit.v, 1.0}, 1.0
	default:
		return obstructedColor, 0
	}
}

func computeShadowRay(hitToLightDir *std.Vec4, hitPoint *std.Vec4, shadowRay *std.Ray) {
	std.Sub(hitToLightDir, hitPoint, shadowRay.Dir)

	shadowRay.Orig = &std.Vec4{}
	shadowRay.Orig[0] = hitPoint[0] + shadowRay.Dir[0]*0.01
	shadowRay.Orig[1] = hitPoint[1] + shadowRay.Dir[1]*0.01
	shadowRay.Orig[2] = hitPoint[2] + shadowRay.Dir[2]*0.01

	shadowRay.Dir.Normalize()
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
