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

func (s Sphere) Bound(t Transform) AABB {
	c := t.Point(s.Center)
	return NewAABB(
		Vector{
			c.X - s.Radius,
			c.Y - s.Radius,
			c.Z - s.Radius,
		},
		Vector{
			c.X + s.Radius,
			c.Y + s.Radius,
			c.Z + s.Radius,
		},
	)
}

func (s Sphere) Intersect(r Ray) (float64, bool) {

	oc := VectorFromTo(s.Center, r.Origin)
	loc := r.Direction.Dot(oc)
	det := loc*loc - oc.Dot(oc) + s.Radius*s.Radius

	// Ray skims the sphere at det==0; ignored
	if det <= 0 {
		return 0, false
	}

	// only return closest intersection point
	d := -loc - math.Sqrt(det)
	if d <= 0 {
		return 0, false
	}
	return d, true
}

func (s Sphere) SurfaceNormal(p Vector) Vector {
	return VectorFromTo(s.Center, p).Normalize()
}
