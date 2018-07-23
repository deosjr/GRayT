package model

import "math"

type Sphere struct {
	object
	Center Vector
	Radius float64
}

func NewSphere(o Vector, r float64, c Color) Sphere {
	return Sphere{
		object: object{c},
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
	return &hit{
		object:   s,
		ray:      r,
		distance: d,
	}
}

func (s Sphere) SurfaceNormal(p Vector) Vector {
	return VectorFromTo(s.Center, p).Normalize()
}
