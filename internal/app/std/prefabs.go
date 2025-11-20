package std

func StandardSpheres() []Sphere {
	spheres := make([]Sphere, 0)

	// Origin
	spheres = append(spheres, Sphere{
		Center:        &Vec4{0, 0, 0, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         Vec4{1, 0, 0},
	})

	// 3 unit to the left (green)
	spheres = append(spheres, Sphere{
		Center:        &Vec4{-3, 0, 0, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         Vec4{0, 1, 0},
	})

	// 3 unit to the right (blue)
	spheres = append(spheres, Sphere{
		Center:        &Vec4{3, 0, 0, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         Vec4{0, 0, 1},
	})

	// 3 units up
	spheres = append(spheres, Sphere{
		Center:        &Vec4{0, 3, 0, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         Vec4{0, 1, 1},
	})

	// Closer to the camera
	spheres = append(spheres, Sphere{
		Center:        &Vec4{0, 0, 3, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         Vec4{1, 1, 0},
	})

	// Even closer
	spheres = append(spheres, Sphere{
		Center:        &Vec4{0, 0, 6, 0},
		Radius:        1,
		RadiusSquared: 1,
		Color:         Vec4{0.8, 0.8, 0.8},
	})

	return spheres
}

func CornellBox() []Plane {
	planes := make([]Plane, 0)
	// left wall
	leftWall := Plane{
		Point:  &Vec4{-6, 0, 0, 0},
		Normal: &Vec4{1, 0, 0, 0},
		Color:  Vec4{0.75, 0.25, 0.25},
	}

	// right wall
	rightWall := Plane{
		Point:  &Vec4{6, 0, 0, 0},
		Normal: &Vec4{-1, 0, 0, 0},
		Color:  Vec4{0.25, 0.25, 0.75},
	}

	// floor
	floor := Plane{
		Point:  &Vec4{0, -4, 0, 0},
		Normal: &Vec4{0, 1, 0, 0},
		Color:  Vec4{0.9, 0.8, 0.7},
	}

	// ceiling
	ceil := Plane{
		Point:  &Vec4{0, 5.5, 0, 0},
		Normal: &Vec4{0, -1, 0, 0},
		Color:  Vec4{0.9, 0.8, 0.7},
	}

	// back wall
	backWall := Plane{
		Point:  &Vec4{0, 0, 2, 0},
		Normal: &Vec4{0, 0, -1, 0},
		Color:  Vec4{0.9, 0.8, 0.7},
	}

	// front wall
	frontWall := Plane{
		Point:  &Vec4{0, 0, -2, 0},
		Normal: &Vec4{0, 0, 1, 0},
		Color:  Vec4{0.9, 0.8, 0.7},
	}

	return append(planes, leftWall, rightWall, floor, ceil, backWall, frontWall)

}
