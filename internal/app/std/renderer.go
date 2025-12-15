package std

import (
	"bytes"
	"fmt"
	"math"
	"time"
)

// hard-coded make-believe sphere representing the point light.
var light = Sphere{
	Center:        &Vec4{2, 5.5, 0},
	Radius:        2,
	RadiusSquared: 4,
	Color:         Vec4{1, 1, 1},
}

func Render(width, height int, spheres []Sphere, planes []Plane, renderTo *bytes.Buffer) {
	hitToLightDir := &Vec4{2, 5.4, 0}

	cameraOrigin := &Vec4{3, 4, 20}
	lookAt := &Vec4{3, 4, 0}
	up := &Vec4{0, 1, 0}
	vt := ViewTransform(cameraOrigin, lookAt, up)
	cameraToWorld := Inverse(vt)             // the inverse view-transform matrix is used in actual rendering.
	cameraToWorld = Transpose(cameraToWorld) // convert to column-major

	fov := 51.52
	scale := float32(math.Tan(deg2rad(fov * 0.5)))
	imageAspectRatio := float32(width) / float32(height)

	// Transform the ray origin to world-space
	var orig = cameraToWorld.MulVec(Vec4{0, 0, 0})

	var ray = &Ray{Orig: &orig}
	var dir = &Vec4{}
	var shadowRay = &Ray{}
	var shadowDir = &Vec4{}
	var hitPoint = &Vec4{}
	var hitNormal = &Vec4{}

	st := time.Now()
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {

			// Useful to troubleshoot / debug a single pixel.
			//if j != 1840 || i != 1166 {
			//	continue
			//}

			// Compute x and y components of ray direction given pixel coords.
			x := (2*(float32(i)+0.5)/float32(width) - 1) * imageAspectRatio * scale
			y := (1 - 2*(float32(j)+0.5)/float32(height)) * scale // Optimize: This can be moved to before the inner for-loop

			// Transform the ray direction using the camera-to-world matrix.
			cameraToWorld.MulDirScalar(x, y, -1, dir)
			dir.Normalize()
			ray.Dir = dir

			minT := float32(3.4028235e38)
			intersectedIdx := -1

			// The code below is absolutely horrific, it needs a severe clean-up and a more uniform handling of ray/sphere/light and
			// shadow rays.
			for n := range spheres {

				t, hit := IntersectSphereSIMD(ray, spheres[n].Center, spheres[n].RadiusSquared)
				if hit && t < minT {
					intersectedIdx = n
					minT = t
				}
			}
			for n := range planes {
				planeT, hit := IntersectPlane2(planes[n].Normal, planes[n].Point, ray.Orig, ray.Dir)
				if hit && planeT < minT {
					intersectedIdx = len(spheres) + n
					minT = planeT
				}
			}

			t, hit := IntersectSphereSIMD(ray, light.Center, light.RadiusSquared)
			if hit && t < minT {
				intersectedIdx = 1000
				minT = t
			}

			if intersectedIdx > -1 {
				// Compute point of hit
				AddMulVec4(ray.Orig, ray.Dir, minT, hitPoint)

				var fract float32
				var color Vec4

				// cast a shadow ray against point light (represented by white sphere)
				Sub(hitToLightDir, hitPoint, shadowDir)
				shadowDir.Normalize()

				shadowRay.Orig = hitPoint // Should translate by epsilon along shadow Dir
				shadowRay.Dir = shadowDir
				obstructed := false

				// In our cornell box, the planes cannot cast shadows anyway, so just try to intersect spheres occluding
				// the path to the point light. (For soft shadows, pick a random point in the hemisphere of the light
				// sphere and sample until results are good. This is (of course) much slower.
				for n := range spheres {
					_, hit := IntersectSphereSIMD(shadowRay, spheres[n].Center, spheres[n].Radius)
					if hit {
						obstructed = true
						break
					}
				}

				// Compute final pixel color. Uses a really ugly trick to discern spheres from planes by offsetting
				// plane intersection index (intersectedIdx) with the number of spheres. intersectedIdx == 1000 is the
				// hard-coded light.
				if obstructed {
					fract = 0.0
					color = Vec4{0, 0, 0}
				} else if intersectedIdx < len(spheres) {
					// Sphere
					Sub(hitPoint, spheres[intersectedIdx].Center, hitNormal)
					hitNormal.Normalize()

					normalToReverseCamera := hitNormal[0]*-ray.Dir[0] + hitNormal[1]*-ray.Dir[1] + hitNormal[2]*-ray.Dir[2]
					fract = max(0.0, normalToReverseCamera)
					color = spheres[intersectedIdx].Color
				} else if intersectedIdx == 1000 {
					// Render light
					fract = 0.96
					color = light.Color
				} else {
					// Plane
					planeNormal := planes[intersectedIdx-len(spheres)].Normal
					normalToReverseCamera := planeNormal[0]*-ray.Dir[0] + planeNormal[1]*-ray.Dir[1] + planeNormal[2]*-ray.Dir[2]
					fract = max(0.0, normalToReverseCamera)
					color = planes[intersectedIdx-len(spheres)].Color
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
	dur := time.Since(st)
	fmt.Printf("duration: %v\n", dur)
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
