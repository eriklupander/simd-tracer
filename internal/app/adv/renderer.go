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

// hitResult holds the outcome of closest-intersection testing for a single ray.
// idx is the raw index within the primitive's own slice — no offset encoding.
type hitResult struct {
	t    float32
	idx  int
	kind int     // IntersectType*
	u, v float32 // barycentric coords, used for triangles
}

func Render(width, height int, spheres *Spheres, planes *Planes, triangles *Triangles, lights []std.Sphere, lightSamples int, renderTo *bytes.Buffer) {
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
	var sphereCenter = &std.Vec4{}
	var lnormal = &std.Vec4{}
	var pointOnHemisphere = &std.Vec4{}

	invLightSamples := float32(1.0) / float32(lightSamples)

	for j := 0; j < height; j++ {

		y := (1 - 2*(float32(j)+0.5)/float32(height)) * scale

		for i := 0; i < width; i++ {

			x := (2*(float32(i)+0.5)/float32(width) - 1) * imageAspectRatio * scale

			cameraToWorld.MulDirScalar(x, y, -1, dir)
			dir.Normalize()
			ray.Dir = dir

			// findClosestHit — inlined
			hitT := float32(math.MaxFloat32)
			hitIdx := 0
			hitKind := IntersectTypeNone
			var hitU, hitV float32
			var color std.Vec4
			var litFract float32
			var finalFract float32
			if t, u, v, idx, ok := IntersectTrianglesSIMD(ray, triangles); ok {
				hitT, hitIdx, hitKind, hitU, hitV = t, idx, IntersectTypeTriangle, u, v
			}
			if t, idx, ok := IntersectPlanesSIMD(ray, planes); ok && t < hitT {
				hitT, hitIdx, hitKind = t, idx, IntersectTypePlane
			}
			if t, idx, ok := IntersectSpheresSIMD(ray, spheres); ok && t < hitT {
				hitT, hitIdx, hitKind = t, idx, IntersectTypeSphere
			}

			// Test primary ray against light spheres; emissive hit takes priority if closer.
			for li := range lights {
				lsphere := &lights[li]
				if t, ok := std.IntersectSphere(ray, lsphere.Center, lsphere.RadiusSquared); ok && t < hitT {
					intensity := float32(lsphere.Material.Intensity)
					renderTo.Write([]byte{
						byte(min(lsphere.Material.Color[0]*intensity*255, 255)),
						byte(min(lsphere.Material.Color[1]*intensity*255, 255)),
						byte(min(lsphere.Material.Color[2]*intensity*255, 255)),
						0xFF,
					})
					goto nextPixel
				}
			}

			if hitKind == IntersectTypeNone {
				renderTo.Write([]byte{0x00, 0x00, 0x00, 0xFF})
				continue
			}

			// Compute hit point, i.e. a point at distance hitT along the ray.
			std.AddMulVec4(ray.Orig, ray.Dir, hitT, hitPoint)

			// Compute surface normal for lighting.
			switch hitKind {
			case IntersectTypeSphere:
				sphereCenter[0] = spheres.CenterX[hitIdx]
				sphereCenter[1] = spheres.CenterY[hitIdx]
				sphereCenter[2] = spheres.CenterZ[hitIdx]
				sphereCenter[3] = 0
				std.Sub(hitPoint, sphereCenter, hitNormal)
				hitNormal.Normalize()
			case IntersectTypePlane:
				hitNormal[0] = planes.NormalX[hitIdx]
				hitNormal[1] = planes.NormalY[hitIdx]
				hitNormal[2] = planes.NormalZ[hitIdx]
				hitNormal[3] = 0
			}

			// Area light sampling: for each light, cast lightSamples shadow rays toward
			// random points on the light sphere and accumulate the N·L diffuse contribution.
			litFract = 0
			for li := range lights {
				lsphere := &lights[li]

				// Direction from the light center to the hit point. After normalising this
				// becomes the outward normal of the light sphere at the point facing us, which
				// defines which hemisphere we pointOnHemisphere — we only want points on the side of the
				// light that is visible from the surface, not the far side.
				lnormal[0] = hitPoint[0] - lsphere.Center[0]
				lnormal[1] = hitPoint[1] - lsphere.Center[1]
				lnormal[2] = hitPoint[2] - lsphere.Center[2]
				lnormal[3] = 0
				lnormal.Normalize()

				for range lightSamples {
					// Pick a uniformly distributed random point on the hemisphere of the
					// light sphere that faces the hit point, writing the world-space position
					// into pointOnHemisphere.
					RandomPointOnHemisphere(*lsphere, *lnormal, pointOnHemisphere)

					// Build a shadow ray from the hit point toward the sampled light position.
					// The direction is the un-normalised vector from hit point to pointOnHemisphere; we
					// offset the origin slightly along that direction to avoid self-intersection
					// with the surface we just hit.
					shadowRay.Dir[0] = pointOnHemisphere[0] - hitPoint[0]
					shadowRay.Dir[1] = pointOnHemisphere[1] - hitPoint[1]
					shadowRay.Dir[2] = pointOnHemisphere[2] - hitPoint[2]
					shadowRay.Dir[3] = 0
					shadowRay.Orig[0] = hitPoint[0] + shadowRay.Dir[0]*0.01
					shadowRay.Orig[1] = hitPoint[1] + shadowRay.Dir[1]*0.01
					shadowRay.Orig[2] = hitPoint[2] + shadowRay.Dir[2]*0.01
					shadowRay.Dir.Normalize()

					// N·L: cosine of the angle between the surface normal and the direction
					// toward this light's pointOnHemisphere. Negative or zero means the light is behind or
					// tangent to the surface — no contribution.
					nl := hitNormal.DotProduct(shadowRay.Dir)
					if nl <= 0 {
						continue
					}
					// Only add the contribution if nothing blocks the path to the light pointOnHemisphere.
					// The N·L term gives a physically correct Lambertian falloff; multiplying by
					// intensity scales the brightness of the light source.
					if !IntersectSpheresSIMDShadowRay(shadowRay, spheres) {
						litFract += nl * float32(lsphere.Material.Intensity)
					}
				}
			}
			// Average over samples so that increasing lightSamples improves quality without
			// changing the overall brightness of the scene.
			litFract *= invLightSamples

			color = shade(hitResult{t: hitT, idx: hitIdx, kind: hitKind, u: hitU, v: hitV}, spheres, planes)

			finalFract = litFract * 255
			renderTo.Write([]byte{byte(color[0] * finalFract), byte(color[1] * finalFract), byte(color[2] * finalFract), 0xFF})
		nextPixel:
		}
	}
}

// shade returns the surface color for a hit. Lighting (N·L) is computed separately in the area light loop.
func shade(hit hitResult, spheres *Spheres, planes *Planes) std.Vec4 {
	switch hit.kind {
	case IntersectTypeSphere:
		return spheres.Color[hit.idx]
	case IntersectTypePlane:
		return planes.Color[hit.idx]
	case IntersectTypeTriangle:
		return std.Vec4{hit.u, hit.v, hit.u * hit.v, 1.0}
	default:
		return obstructedColor
	}
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
