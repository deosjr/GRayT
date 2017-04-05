package model

import "math"

type Sphere struct {
	Center Vector
	Radius float64
}

func (s Sphere) Intersect(r Ray) (Vector, bool, float64) {

	oc := VectorFromTo(s.Center, r.Origin)
	loc := r.Direction.Dot(oc)
	det := loc*loc - oc.Dot(oc) + s.Radius*s.Radius

	// Ray skims the sphere at det==0; ignored
	if det <= 0 {
		return Vector{}, false, 0
	}

	// only return closest intersection point
	d := -loc - math.Sqrt(det)
	return r.Origin.Add(r.Direction.Times(d)), true, d
}

func (s Sphere) SurfaceNormal(p Vector) Vector {
	return VectorFromTo(s.Center, p).Normalize()
}
