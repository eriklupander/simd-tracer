package std

type Plane struct {
	Point  *Vec4
	Normal *Vec4
	Color  Vec4
}

// From https://lousodrome.net/blog/light/2020/07/03/intersection-of-a-ray-and-a-plane/
func IntersectPlane2(normal *Vec4, pointOnPlane *Vec4, orig *Vec4, dir *Vec4) (float32, bool) {
	lower := dir.DotProduct(normal)
	if lower > -0.00001 {
		return 0.0, false
	}
	var tmp1 Vec4
	Sub(pointOnPlane, orig, &tmp1)

	upper := tmp1.DotProduct(normal)

	t := upper / lower

	return t, t > 0.0
}
