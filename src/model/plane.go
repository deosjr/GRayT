package model

import "math"

type Plane struct {
	object
	Point  Vector
	Normal Vector
}

func NewPlane(p Vector, u, v Vector, c Color) Plane {
	n := u.Cross(v).Normalize()
	return Plane{
		object: object{c},
		Point:  p,
		Normal: n,
	}
}

// NOTE: planes dont work nicely with BVH since their bounding box is infinite
// They force a lot more intersection tests, slowing everything down

func (p Plane) Bound() AABB {
	return NewAABB(
		Vector{
			-math.MaxFloat64,
			-math.MaxFloat64,
			-math.MaxFloat64},
		Vector{
			math.MaxFloat64,
			math.MaxFloat64,
			math.MaxFloat64},
	)
}

func (p Plane) Intersect(r Ray) (float64, bool) {
	ln := r.Direction.Dot(p.Normal)
	if ln == 0 {
		// line and plane parallel
		return 0, false
	}
	d := VectorFromTo(r.Origin, p.Point).Dot(p.Normal) / ln
	if d <= 0 {
		return 0, false
	}
	return d, true
}

func (p Plane) SurfaceNormal(Vector) Vector {
	// TODO: determine direction?!
	// solution: I think 2d faces have only one side that can be lit;
	// related to direction of surface normal.
	// with facingRatio this means other side is always pure black
	return p.Normal.Times(1)
}
