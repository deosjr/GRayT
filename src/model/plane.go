package model

type Plane struct {
	Point  Vector
	Normal Vector
	Color  Color
}

func NewPlane(p Vector, u, v Vector, c Color) Plane {
	n := u.Cross(v).Normalize()
	return Plane{
		Point:  p,
		Normal: n,
		Color:  c,
	}
}

func (p Plane) GetColor() Color {
	return p.Color
}

func (p Plane) Intersect(r Ray) (float64, bool) {
	ln := r.Direction.Dot(p.Normal)
	if ln == 0 {
		// line and plane parallel
		return 0, false
	}
	d := p.Point.Sub(r.Origin).Dot(p.Normal) / ln
	return d, true
}

func (p Plane) SurfaceNormal(Vector) Vector {
	// TODO: determine direction?!
	// solution: I think 2d faces have only one side that can be lit;
	// related to direction of surface normal.
	// with facingRatio this means other side is always pure black
	return p.Normal.Times(1)
}
