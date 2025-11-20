package std

import (
	"fmt"
	"math"
	"time"
)

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func Render(width, height int, spheres []Sphere, planes []Plane, renderTo []byte) {
	hitToLightDir := &Vec4{2, 5.4, 0}
	light := Sphere{
		Center:        &Vec4{2, 5.5, 0},
		Radius:        2,
		RadiusSquared: 4,
		Color:         Vec4{1, 1, 1},
	}

	//cameraToWorld := NewMatrix44f( // Original
	//	0.945519, 0, -0.325569, 0,
	//	-0.179534, 0.834209, -0.521403, 0,
	//	0.271593, 0.551447, 0.78876, 0,
	//	4.208271, 8.374532, 17.932925, 1,
	//)
	//cameraToWorld := NewMatrix44f(
	//	1, 0, 0, 0, // rotation around X axis (tilt up/down)
	//	0, 1, 0, 0, // rotation around Y axis (tilt left/right)
	//	0, 0, 1, 0, // rotation around Z axis (to
	//	3, 4, 20, 1, // Camera position 10 units into the screen
	//)
	cameraOrigin := &Vec4{3, 4, 20}
	lookAt := &Vec4{3, 4, 0}
	up := &Vec4{0, 1, 0}
	vt := ViewTransform(cameraOrigin, lookAt, up)
	cameraToWorld := Inverse(vt)             // the inverse view-transform matrix is used in actual rendering.
	cameraToWorld = Transpose(cameraToWorld) // convert to column-major

	fov := 51.52
	scale := float32(math.Tan(deg2rad(fov * 0.5)))
	imageAspectRatio := float32(width) / float32(height)
	// [comment]
	// Don't forget to transform the ray origin (which is also the camera origin
	// by transforming the point with coordinates (0,0,0) to world-space using the
	// camera-to-world matrix.
	// [/comment]
	var orig = cameraToWorld.mulVec2(Vec4{0, 0, 0})

	var ray = &Ray{Orig: &orig}
	var dir = &Vec4{}
	var shadowRay = &Ray{}
	var shadowDir = &Vec4{}
	var hitPoint = &Vec4{}
	var hitNormal = &Vec4{}
	var cast = 0
	st := time.Now()
	for j := 0; j < height; j++ {

		for i := 0; i < width; i++ {

			//if j != 1200 || i != 1216 {
			//	continue
			//}

			// [comment]
			// Generate primary ray direction. Compute the x and y position
			// of the ray in screen space. This gives a point on the image plane
			// at z=1. From there, we simply compute the direction by normalized
			// the resulting vec3f variable. This is similar to taking the vector
			// between the point on the image plane and the camera origin, which
			// in camera space is (0,0,0):
			//
			// ray.Dir = normalize(Vec3f(x,y,-1) - Vec3f(0));
			// [/comment]

			x := (2*(float32(i)+0.5)/float32(width) - 1) * imageAspectRatio * scale
			y := (1 - 2*(float32(j)+0.5)/float32(height)) * scale // Optimize: This can be moved to before the inner for-loop

			// Don't forget to transform the ray direction using the camera-to-world matrix.
			cameraToWorld.mulDirScalar(x, y, -1, dir)
			dir.normalize()
			ray.Dir = dir

			minT := float32(3.4028235e38)
			intersectedIdx := -1

			for n := range spheres {
				cast++
				t, hit := IntersectSphereFullSIMD(ray, spheres[n].Center, spheres[n].RadiusSquared)
				if hit && t < minT {
					intersectedIdx = n
					minT = t
				}
			}
			for n := range planes {
				cast++
				planeT, hit := intersectPlane2(planes[n].Normal, planes[n].Point, ray.Orig, ray.Dir)
				if hit && planeT < minT {
					intersectedIdx = len(spheres) + n
					minT = planeT
				}
			}

			t, hit := IntersectSphereFullSIMD(ray, light.Center, light.RadiusSquared)
			if hit && t < minT {
				intersectedIdx = 1000
				minT = t
			}

			// Prepare RGBA image index to render to
			renderToByteIndex := (j*width + i) * 4

			if intersectedIdx > -1 {
				// Compute point of hit

				addMulVec4(ray.Orig, ray.Dir, minT, hitPoint)

				var fract float32
				var color Vec4

				// Temp, cast a shadow ray against point light represented by white sphere

				sub2(hitToLightDir, hitPoint, shadowDir)
				shadowDir.normalize()

				shadowRay.Orig = hitPoint // Should translate by epsilon along shadow Dir
				shadowRay.Dir = shadowDir
				obstructed := false
				for n := range spheres {
					cast++
					_, hit := IntersectSphereFullSIMD(shadowRay, spheres[n].Center, spheres[n].Radius)
					if hit {
						obstructed = true
						break
					}
				}
				if obstructed {
					fract = 0.0
					color = Vec4{0, 0, 0}
					// simulate ambient color by using correct color with low fract
				} else if intersectedIdx < len(spheres) {
					// Sphere

					sub2(hitPoint, spheres[intersectedIdx].Center, hitNormal)
					hitNormal.normalize()

					normalToReverseCamera := hitNormal[0]*-ray.Dir[0] + hitNormal[1]*-ray.Dir[1] + hitNormal[2]*-ray.Dir[2]
					fract = max(0.0, normalToReverseCamera)
					color = spheres[intersectedIdx].Color
				} else if intersectedIdx == 1000 {
					// Render light
					//var hitNormal = &Vec3{}
					//sub2(hitPoint, light.Center, hitNormal)
					//hitNormal.normalize()
					//
					//normalToReverseCamera := hitNormal[0]*-ray.Dir[0] + hitNormal[1]*-ray.Dir[1] + hitNormal[2]*-ray.Dir[2]
					//fract = max(0.0, normalToReverseCamera)
					fract = 0.96
					color = light.Color
				} else {
					// Plane
					planeNormal := planes[intersectedIdx-len(spheres)].Normal
					normalToReverseCamera := planeNormal[0]*-ray.Dir[0] + planeNormal[1]*-ray.Dir[1] + planeNormal[2]*-ray.Dir[2]
					fract = max(0.0, normalToReverseCamera)
					color = planes[intersectedIdx-len(spheres)].Color
				}

				renderTo[renderToByteIndex] = uint8(color[0] * fract * 255)
				renderTo[renderToByteIndex+1] = uint8(color[1] * fract * 255)
				renderTo[renderToByteIndex+2] = uint8(color[2] * fract * 255)
				renderTo[renderToByteIndex+3] = 0xFF // No transparency
			} else {
				renderTo[renderToByteIndex] = 0x00
				renderTo[renderToByteIndex+1] = 0x00
				renderTo[renderToByteIndex+2] = 0x00
				renderTo[renderToByteIndex+3] = 0xFF // No transparency
			}
		}
	}
	dur := time.Since(st)
	fmt.Printf("duration: %v\n", dur)
	fmt.Printf("rays cast: %d\n", cast)
	fmt.Printf("intersections/s: %f\n", float64(cast)/dur.Seconds())
}

//
//func solveQuadratic(a, b, c float32) (x0, x1 float32, ret bool) {
//	discr := b*b - 4*a*c
//	if discr < 0 {
//		return 0.0, 0.0, false
//	} else if discr == 0 {
//		x0 = -0.5 * b / a
//		x1 = -0.5 * b / a
//	} else {
//		q := float32(0.0)
//		if b > 0 {
//			q = -0.5 * (b + float32(math.Sqrt(float64(discr))))
//		} else {
//			q = -0.5 * (b - float32(math.Sqrt(float64(discr))))
//		}
//
//		x0 = q / a
//		x1 = c / q
//	}
//
//	return x0, x1, true
//}
