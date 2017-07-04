package model

type Triangle struct {
	P0    Vector
	P1    Vector
	P2    Vector
	Color Color
}

func (t Triangle) GetColor() Color {
	return t.Color
}

// Normal points towards side where points are numbered counter-clockwise
func (t Triangle) Intersect(r Ray) (float64, bool) {
	// n = (P1 - P0) x (P2 - P0)
	n := t.P1.Sub(t.P0).Cross(t.P2.Sub(t.P0))
	// now we have a plane with point P0 and normal n
	// so let's use plane intersection logic
	ln := r.Direction.Dot(n)
	if ln == 0 {
		// line and plane parallel
		return 0, false
	}
	d := t.P0.Sub(r.Origin).Dot(n) / ln
	if d <= 0 {
		return 0, false
	}

	// we have an intersection with the plane,
	// now we need to decide whether the point (x) will be in the triangle
	x := PointFromRay(r, d)

	// (P1 - P0) x (X - P0) . n >= 0
	if t.P1.Sub(t.P0).Cross(x.Sub(t.P0)).Dot(n) < 0 {
		return 0, false
	}

	// (P2 - P1) x (X - P1) . n >= 0
	if t.P2.Sub(t.P1).Cross(x.Sub(t.P1)).Dot(n) < 0 {
		return 0, false
	}

	// (P0 - P2) x (X - P2) . n >= 0
	if t.P0.Sub(t.P2).Cross(x.Sub(t.P2)).Dot(n) < 0 {
		return 0, false
	}

	return d, true
}

func (t Triangle) SurfaceNormal(Vector) Vector {
	return t.P1.Sub(t.P0).Cross(t.P2.Sub(t.P0)).Times(1)
}
