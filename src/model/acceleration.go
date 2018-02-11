package model

// floating point error margin
// TODO: setting too small drops shadows completely?
// setting to 0.1 or 0.5 shows shadows; setting too big gives weirdness
// see https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/ligth-and-shadows
// on shadow bias
var ERROR_MARGIN = 1E-10

type AccelerationStructure interface {
	ClosestIntersection(ray Ray, maxDistance float64) *hit
}

type hit struct {
	object Object
	point  Vector
}

type NaiveAcceleration struct {
	objects []Object
}

func NewNaiveAcceleration(objects []Object) NaiveAcceleration {
	return NaiveAcceleration{objects: objects}
}

// Try and hit ALL objects EVERY time
func (na NaiveAcceleration) ClosestIntersection(ray Ray, maxDistance float64) *hit {
	var objectHit Object
	d := maxDistance
	for _, o := range na.objects {
		if distance, ok := o.Intersect(ray); ok && distance < d && distance > ERROR_MARGIN {
			d = distance
			objectHit = o
		}
	}
	if d == maxDistance {
		return nil
	}
	return &hit{
		object: objectHit,
		point:  PointFromRay(ray, d),
	}
}
