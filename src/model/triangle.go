package model

// TODO: optimizations
// - memory: vertex sharing
// - speed: SIMD instructions

type Triangle struct {
	object
	P0 Vector
	P1 Vector
	P2 Vector
}

func NewTriangle(p0, p1, p2 Vector, c Color) Triangle {
	return Triangle{
		object: object{c},
		P0:     p0,
		P1:     p1,
		P2:     p2,
	}
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

func (t Triangle) SurfaceNormal(Vector) Vector {
	return VectorFromTo(t.P0, t.P1).Cross(VectorFromTo(t.P0, t.P2)).Normalize().Times(1)
}
