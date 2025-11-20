package std

// intersectPlane was fixed by removing a previously existing check for a positive denominator
//func intersectPlane(normal *Vec4, pointOnPlane *Vec4, Orig *Vec4, Dir *Vec4) (float32, bool) {
//	// Assuming vectors are all normalized
//
//	denom := normal.dotProduct(Dir)
//
//	// if denom > 1e-6 { // did NOT work, denom is typically negative with our setup given a normal pointing upwards {0,1,0}
//	// and a pointOnPlane of (for example) 0,-1,0. What is wrong?
//	var p0l0 = Vec4{}
//	sub2(pointOnPlane, Orig, &p0l0)
//
//	t := p0l0.dotProduct(normal) / denom
//
//	return t, t >= 0
//}

type Plane struct {
	Point  *Vec4
	Normal *Vec4
	Color  Vec4
}

// From https://lousodrome.net/blog/light/2020/07/03/intersection-of-a-ray-and-a-plane/
func intersectPlane2(normal *Vec4, pointOnPlane *Vec4, orig *Vec4, dir *Vec4) (float32, bool) {
	lower := dir.dotProduct(normal)
	if lower > -0.00001 {
		return 0.0, false
	}
	var tmp1 Vec4
	sub2(pointOnPlane, orig, &tmp1)

	upper := tmp1.dotProduct(normal)

	t := upper / lower

	return t, t > 0.0
}
