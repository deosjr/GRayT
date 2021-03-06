package model

import "math"

type Sphere struct {
	object
	Center Vector
	Radius float32
}

func NewSphere(o Vector, r float32, m Material) Sphere {
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

func (s Sphere) Intersect(r Ray) (*SurfaceInteraction, bool) {

	oc := VectorFromTo(s.Center, r.Origin)
	loc := r.Direction.Dot(oc)
	det := loc*loc - oc.Dot(oc) + s.Radius*s.Radius

	// Ray skims the sphere at det==0; ignored
	if det <= 0 {
		return nil, false
	}

	// only return closest intersection point
	d := -loc - float32(math.Sqrt(float64(det)))
	if d <= 0 {
		return nil, false
	}
	n := s.SurfaceNormal(PointFromRay(r, d))
	return NewSurfaceInteraction(s, d, n, r), true
}

func (s Sphere) SurfaceNormal(p Vector) Vector {
	return VectorFromTo(s.Center, p).Normalize()
}
