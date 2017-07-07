package model

// TODO: optimizations
// - memory: vertex sharing
// - speed: SIMD instructions

type Triangle struct {
	P0    Vector
	P1    Vector
	P2    Vector
	Color Color
}

func (t Triangle) GetColor() Color {
	return t.Color
}

// Moller-Trumbore intersection algorithm
func (t Triangle) Intersect(r Ray) (float64, bool) {
	e1 := t.P1.Sub(t.P0)
	e2 := t.P2.Sub(t.P0)
	pvec := r.Direction.Cross(e2)
	det := e1.Dot(pvec)

	if det < 1e-8 && det > -1e-8 {
		return 0, false
	}
	inv_det := 1.0 / det

	tvec := r.Origin.Sub(t.P0)
	u := tvec.Dot(pvec) * inv_det
	if u < 0 || u > 1 {
		return 0, false
	}

	qvec := tvec.Cross(e1)
	v := r.Direction.Dot(qvec) * inv_det
	if v < 0 || u+v > 1 {
		return 0, false
	}
	return e2.Dot(qvec) * inv_det, true
}

// TODO: currently unused!
// Naive implementation using plane intersection
// Normal points towards side where points are numbered counter-clockwise
func (t Triangle) intersect(r Ray) (float64, bool) {
	// n = (P1 - P0) x (P2 - P0)
	n := VectorFromTo(t.P0, t.P1).Cross(VectorFromTo(t.P0, t.P2)).Normalize()
	// now we have a plane with point P0 and normal n
	// so let's use plane intersection logic
	ln := r.Direction.Dot(n)
	if ln == 0 {
		// line and plane parallel
		return 0, false
	}
	d := VectorFromTo(r.Origin, t.P0).Dot(n) / ln
	if d <= 0 {
		return 0, false
	}

	// we have an intersection with the plane,
	// now we need to decide whether the point (x) will be in the triangle
	x := PointFromRay(r, d)

	// (P1 - P0) x (X - P0) . n >= 0
	if VectorFromTo(t.P0, t.P1).Cross(VectorFromTo(t.P0, x)).Dot(n) < 0 {
		return 0, false
	}

	// (P2 - P1) x (X - P1) . n >= 0
	if VectorFromTo(t.P1, t.P2).Cross(VectorFromTo(t.P1, x)).Dot(n) < 0 {
		return 0, false
	}

	// (P0 - P2) x (X - P2) . n >= 0
	if VectorFromTo(t.P2, t.P0).Cross(VectorFromTo(t.P2, x)).Dot(n) < 0 {
		return 0, false
	}

	return d, true
}

func (t Triangle) SurfaceNormal(Vector) Vector {
	return VectorFromTo(t.P0, t.P1).Cross(VectorFromTo(t.P0, t.P2)).Normalize().Times(1)
}
