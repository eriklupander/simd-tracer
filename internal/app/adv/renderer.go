package adv

import (
	"bytes"
	"math"

	"github.com/eriklupander/simd-tracer/internal/app/std"
)

var obstructedColor = std.Vec4{0, 0, 0}
var red = std.Vec4{1, 0, 0}

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

func Render(width, height int, spheres *Spheres, planes *Planes, triangles *Triangles, renderTo *bytes.Buffer) {
	hitToLightDir := &std.Vec4{2, 5.4, 0}

	cameraOrigin := &std.Vec4{3, 4, 20}
	lookAt := &std.Vec4{3, 4, 0}
	up := &std.Vec4{0, 1, 0}
	vt := std.ViewTransform(cameraOrigin, lookAt, up)
	cameraToWorld := std.Inverse(vt)             // the inverse view-transform matrix is used in actual rendering.
	cameraToWorld = std.Transpose(cameraToWorld) // convert to column-major

	fov := 51.52
	scale := float32(math.Tan(deg2rad(fov * 0.5)))
	imageAspectRatio := float32(width) / float32(height)

	// Transform the ray origin to world-space
	var orig = cameraToWorld.MulVec(std.Vec4{0, 0, 0})

	var ray = &std.Ray{Orig: &orig}
	var dir = &std.Vec4{}
	var shadowRay = &std.Ray{Dir: &std.Vec4{}, Orig: &std.Vec4{}}
	var hitPoint = &std.Vec4{}
	var hitNormal = &std.Vec4{}
	var intersectedType = IntersectTypeNone

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {

			// Useful to troubleshoot / debug a single pixel.

			// Compute x and y components of ray direction given pixel coords.
			x := (2*(float32(i)+0.5)/float32(width) - 1) * imageAspectRatio * scale
			y := (1 - 2*(float32(j)+0.5)/float32(height)) * scale // Optimize: This can be moved to before the inner for-loop

			// Transform the ray direction using the camera-to-world matrix.
			cameraToWorld.MulDirScalar(x, y, -1, dir)
			dir.Normalize()
			ray.Dir = dir

			//if j != 370 || i != 150 {
			//	continue
			//}
			//fmt.Printf("%v\n", ray.Dir)
			minT := float32(3.4028235e38)
			intersectedIdx := -1

			// Perform ray-triangle intersection testing.
			//tris := triangles.ToSlice()
			//for idx, tri := range tris {
			//	t, _, hit := IntersectTriangle(ray, tri.v0, tri.v1, tri.v2)
			//	if hit && t < minT {
			//		intersectedIdx = planes.Count + spheres.Count + idx
			//		minT = t
			//	}
			//}

			t, u, v, idx, hit := IntersectTrianglesSIMD(ray, triangles)
			if hit {
				intersectedIdx = planes.Count + spheres.Count + idx
				minT = t
				intersectedType = IntersectTypeTriangle
			}

			// Perform ray-plane intersection testing.
			t, idx, hit = IntersectPlanesSIMD(ray, planes)
			if hit && t < minT {
				intersectedIdx = spheres.Count + idx
				minT = t
				intersectedType = IntersectTypePlane
			}

			// Perform ray-sphere intersection testing.
			t, idx, hit = IntersectSpheresSIMD(ray, spheres)
			if hit && t < minT {
				intersectedIdx = idx
				minT = t
				intersectedType = IntersectTypeSphere
			}

			// The code below is absolutely horrific, it needs a severe clean-up and a more uniform handling of ray/sphere/light and
			// shadow rays.
			if intersectedIdx > -1 {
				// Compute point of hit
				std.AddMulVec4(ray.Orig, ray.Dir, minT, hitPoint)

				var fract float32
				var color std.Vec4

				// cast a shadow ray against point light (represented by white sphere)
				computeShadowRay(hitToLightDir, hitPoint, shadowRay)

				// For shadows, first compute hit normal and then check if it is facing away from the light source
				var obstructed = false
				if intersectedType == IntersectTypeSphere {
					std.Sub(hitPoint, &std.Vec4{spheres.CenterX[intersectedIdx], spheres.CenterY[intersectedIdx], spheres.CenterZ[intersectedIdx], 0}, hitNormal)
					hitNormal.Normalize()
					if hitNormal.DotProduct(shadowRay.Dir) <= 0 {
						obstructed = true
					}
				}

				// In our cornell box, the planes cannot cast shadows anyway, so just try to intersect spheres occluding
				// the path to the point light.

				// Check first if self-shadowed, no need to perform expensive shadow ray casting if intersected point is
				// pointing away from direction of point light source.
				if !obstructed {
					obstructed = IntersectSpheresSIMDShadowRay(shadowRay, spheres)
				}

				// Compute final pixel color. Uses a hideous trick to discern spheres from planes and triangles by offsetting
				// the intersection index (intersectedIdx) with the number of spheres and/or triangles.
				if obstructed {
					fract = 0.0
					color = obstructedColor
				} else {
					switch intersectedType {
					case IntersectTypeSphere:
						// Note: Hit normal already computed.

						normalToReverseCamera := hitNormal[0]*-ray.Dir[0] + hitNormal[1]*-ray.Dir[1] + hitNormal[2]*-ray.Dir[2]
						fract = max(0.0, normalToReverseCamera)
						color = spheres.Color[intersectedIdx]
					case IntersectTypePlane:
						// Plane
						planeIdx := intersectedIdx - spheres.Count
						planeNormal := std.Vec4{planes.NormalX[planeIdx], planes.NormalY[planeIdx], planes.NormalZ[planeIdx], 0}
						normalToReverseCamera := planeNormal[0]*-ray.Dir[0] + planeNormal[1]*-ray.Dir[1] + planeNormal[2]*-ray.Dir[2]
						fract = max(0.0, normalToReverseCamera)
						color = planes.Color[planeIdx]
					case IntersectTypeTriangle:
						// Triangle. Note that this is just placeholder/boilerplate.
						color = std.Vec4{1.0 * u, 1.0 * v, 1.0 * u * v, 1.0}
						fract = 1.0
					default:
						// No intersection
					}
				}

				// Render RGBA for current pixel in one go. fract is multiplied by 255 so it can be stuffed directly into
				// the output buffer as uint8.
				finalFract := fract * 255
				renderTo.Write([]byte{byte(color[0] * finalFract), byte(color[1] * finalFract), byte(color[2] * finalFract), 0xFF})
			} else {
				// RGBA, no transparency.
				renderTo.Write([]byte{0x00, 0x00, 0x00, 0xFF})
			}
		}
	}
}

func computeShadowRay(hitToLightDir *std.Vec4, hitPoint *std.Vec4, shadowRay *std.Ray) {
	std.Sub(hitToLightDir, hitPoint, shadowRay.Dir)

	shadowRay.Orig = &std.Vec4{}
	shadowRay.Orig[0] = hitPoint[0] + shadowRay.Dir[0]*0.01 // Translate by epsilon along shadow Dir
	shadowRay.Orig[1] = hitPoint[1] + shadowRay.Dir[1]*0.01
	shadowRay.Orig[2] = hitPoint[2] + shadowRay.Dir[2]*0.01

	shadowRay.Dir.Normalize()
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
