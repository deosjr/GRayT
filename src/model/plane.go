package model

type Plane struct {
	Point  Vector
	Normal Vector
}

func NewPlane(u, v Vector, p Vector) Plane {
	n := u.Cross(v).Normalize()
	return Plane{
		Point:  p,
		Normal: n,
	}
}

func (p Plane) Intersect(r Ray) (Vector, bool, float64) {
	ln := r.Direction.Dot(p.Normal)
	if ln == 0 {
		// line and plane parallel
		return Vector{}, false, 0
	}
	d := p.Point.Sub(r.Origin).Dot(p.Normal) / ln
	intersection := r.Origin.Add(r.Direction.Times(d))
	return intersection, true, d
}

func (p Plane) SurfaceNormal(Vector) Vector {
	// TODO: determine direction?!
	// solution: I think 2d faces have only one side that can be lit;
	// related to direction of surface normal.
	// with facingRatio this means other side is always pure black
	return p.Normal.Times(1)
}
