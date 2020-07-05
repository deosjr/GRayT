package model

// floating point error margin
// TODO: setting too small drops shadows completely?
// setting to 0.1 or 0.5 shows shadows; setting too big gives weirdness
// see https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/ligth-and-shadows
// on shadow bias
var ERROR_MARGIN float32 = 1e-3
var SIMD_ENABLED bool = false

type AccelerationStructure interface {
	GetObjects() []Object
	ClosestIntersection(ray Ray, maxDistance float32) (*SurfaceInteraction, bool)
}

type NaiveAcceleration struct {
	objects []Object
}

func NewNaiveAcceleration(objects []Object) *NaiveAcceleration {
	return &NaiveAcceleration{objects: objects}
}

func (na NaiveAcceleration) GetObjects() []Object {
	return na.objects
}

// Try and hit ALL objects EVERY time
func (na *NaiveAcceleration) ClosestIntersection(ray Ray, maxDistance float32) (*SurfaceInteraction, bool) {
	var found bool
	var surfaceInteraction *SurfaceInteraction
	distance := maxDistance
	for _, o := range na.objects {
		if si, ok := o.Intersect(ray); ok && si.distance < distance && si.distance > ERROR_MARGIN {
			distance = si.distance
			surfaceInteraction = si
			found = true
		}
	}
	if !found {
		return nil, false
	}
	return surfaceInteraction, true
}
