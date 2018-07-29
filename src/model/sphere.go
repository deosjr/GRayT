package model

import "math"

type Sphere struct {
	object
	Center Vector
	Radius float64
}

func NewSphere(o Vector, r float64, m Material) Sphere {
	return Sphere{
		object: object{m},
		Center: o,
		Radius: r,
	}
}

func (s Sphere) Bound() AABB {
	return NewAABB(
		Vector{
			s.Center.X - s.Radius,
			s.Center.Y - s.Radius,
			s.Center.Z - s.Radius,
		},
		Vector{
			s.Center.X + s.Radius,
			s.Center.Y + s.Radius,
			s.Center.Z + s.Radius,
		},
	)
}

func (s Sphere) Intersect(r Ray) *hit {

	oc := VectorFromTo(s.Center, r.Origin)
	loc := r.Direction.Dot(oc)
	det := loc*loc - oc.Dot(oc) + s.Radius*s.Radius

	// Ray skims the sphere at det==0; ignored
	if det <= 0 {
		return nil
	}

	// only return closest intersection point
	d := -loc - math.Sqrt(det)
	if d <= 0 {
		return nil
	}
	return NewHit(s, r, d)
}

func (s Sphere) SurfaceNormal(p Vector) Vector {
	return VectorFromTo(s.Center, p).Normalize()
}
