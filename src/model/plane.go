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
