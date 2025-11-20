package std

// ViewTransform creates a column-minor view-transformation matrix. When used in simd-tracer, which uses col-major,
// remember to call transpose first.
func ViewTransform(from, to, up *Vec4) Matrix44 {
	// Create a new matrix from the identity matrix.
	vt := NewIdentityMatrix44()

	// Sub creates the initial vector between the eye and what we're looking at.
	var forward Vec4
	sub2(to, from, &forward)

	forward.normalize()

	// Normalize the up vector
	up.normalize()

	// Use the cross product to get the "third" axis (in this case, not the forward or up one)
	var left Vec4
	forward.crossProduct(up, &left)

	// Again, use cross product between the just computed left and forward to get the "true" up.
	var trueUp Vec4
	left.crossProduct(&forward, &trueUp)

	// copy each axis into the matrix
	vt[0][0] = left[0]
	vt[0][1] = left[1]
	vt[0][2] = left[2]

	vt[1][0] = trueUp[0]
	vt[1][1] = trueUp[1]
	vt[1][2] = trueUp[2]

	vt[2][0] = -forward[0]
	vt[2][1] = -forward[1]
	vt[2][2] = -forward[2]

	// finally, move the view matrix opposite the camera position to emulate that the camera has moved.
	translationM := newTranslationMatrix(-from[0], -from[1], -from[2])
	return multiplyMatricies(vt, translationM)
}
