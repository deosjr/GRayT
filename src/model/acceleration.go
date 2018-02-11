package model

import "math"

type AccelerationStructure interface {
	ClosestIntersection(ray Ray) *hit
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
func (na NaiveAcceleration) ClosestIntersection(ray Ray) *hit {
	var objectHit Object
	d := math.MaxFloat64
	for _, o := range na.objects {
		if distance, ok := o.Intersect(ray); ok && distance < d {
			d = distance
			objectHit = o
		}
	}
	if d == math.MaxFloat64 {
		return nil
	}
	return &hit{
		object: objectHit,
		point:  PointFromRay(ray, d),
	}
}
