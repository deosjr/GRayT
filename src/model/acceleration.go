package model

// floating point error margin
// TODO: setting too small drops shadows completely?
// setting to 0.1 or 0.5 shows shadows; setting too big gives weirdness
// see https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/ligth-and-shadows
// on shadow bias
var ERROR_MARGIN = 1E-10

type AccelerationStructure interface {
	ClosestIntersection(ray Ray, maxDistance float64) (hit, bool)
}

type NaiveAcceleration struct {
	objects []Object
}

func NewNaiveAcceleration(objects []Object) *NaiveAcceleration {
	return &NaiveAcceleration{objects: objects}
}

// Try and hit ALL objects EVERY time
func (na *NaiveAcceleration) ClosestIntersection(ray Ray, maxDistance float64) (hit, bool) {
	var hit hit
	var found bool
	distance := maxDistance
	for _, o := range na.objects {
		if h, ok := o.Intersect(ray); ok && h.distance < distance && h.distance > ERROR_MARGIN {
			distance = h.distance
			found = true
			hit = h
		}
	}
	return hit, found
}
