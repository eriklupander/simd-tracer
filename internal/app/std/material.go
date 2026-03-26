package std

// Material describes the surface properties of a primitive.
type Material struct {
	Color          Vec4
	Emissive       bool
	Intensity      float64
	Reflectivity   float32
	RefractiveIndex float32
}
