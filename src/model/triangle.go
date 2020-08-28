package model

import (
	"math/rand"
)

type Triangle struct {
	object
	P0 Vector
	P1 Vector
	P2 Vector
}

func triangleBound(p0, p1, p2 Vector, t Transform) AABB {
	tp0 := t.Point(p0)
	tp1 := t.Point(p1)
	tp2 := t.Point(p2)
	return NewAABB(tp0, tp1).AddPoint(tp2)
}

// Moller-Trumbore intersection algorithm
func triangleIntersect(p0, p1, p2 Vector, ray Ray) (float32, bool) {
	e1 := p1.Sub(p0)
	e2 := p2.Sub(p0)
	pvec := ray.Direction.Cross(e2)
	det := e1.Dot(pvec)

	if det < 1e-8 && det > -1e-8 {
		return 0, false
	}
	inv_det := 1.0 / det

	tvec := ray.Origin.Sub(p0)
	u := tvec.Dot(pvec) * inv_det
	if u < 0 || u > 1 {
		return 0, false
	}

	qvec := tvec.Cross(e1)
	v := ray.Direction.Dot(qvec) * inv_det
	if v < 0 || u+v > 1 {
		return 0, false
	}
	return e2.Dot(qvec) * inv_det, true
}

func triangleSurfaceNormal(p0, p1, p2 Vector) Vector {
	return VectorFromTo(p0, p1).Cross(VectorFromTo(p0, p2)).Normalize()
}

func NewTriangle(p0, p1, p2 Vector, m Material) Triangle {
	return Triangle{
		object: object{m},
		P0:     p0,
		P1:     p1,
		P2:     p2,
	}
}

func (t Triangle) Bound(transform Transform) AABB {
	return triangleBound(t.P0, t.P1, t.P2, transform)
}

func (t Triangle) Intersect(r Ray) (*SurfaceInteraction, bool) {
	d, ok := triangleIntersect(t.P0, t.P1, t.P2, r)
	if !ok {
		return nil, false
	}
	n := triangleSurfaceNormal(t.P0, t.P1, t.P2)
	return NewSurfaceInteraction(t, d, n, r), true
}

func (t Triangle) IntersectOptimized(r Ray) (float32, bool) {
	d, ok := triangleIntersect(t.P0, t.P1, t.P2, r)
	if !ok {
		return 0, false
	}
	return d, true
}

func (t Triangle) SurfaceNormal(Vector) Vector {
	return triangleSurfaceNormal(t.P0, t.P1, t.P2)
}

func (t Triangle) Sample(random *rand.Rand) Vector {
	u := t.P1.Sub(t.P0)
	v := t.P2.Sub(t.P0)
	a, b := random.Float32(), random.Float32()
	if a+b > 1 {
		a, b = 1-a, 1-b
	}
	return t.P0.Add(u.Times(a)).Add(v.Times(b))
}

func triangleSurfaceArea(p0, p1, p2 Vector) float32 {
	return p1.Sub(p0).Cross(p2.Sub(p0)).Length() * 0.5
}

func (t Triangle) SurfaceArea() float32 {
	return triangleSurfaceArea(t.P0, t.P1, t.P2)
}

// for point P outside triangle T, this might not be very meaningful
// TODO: use (and understand) the edge function to compute this more efficiently
// PBRT does this _in_ the triangle intersect function!
func barycentric(p0, p1, p2, p Vector) (float32, float32, float32) {
	area := triangleSurfaceArea(p0, p1, p2)
	l0 := triangleSurfaceArea(p1, p2, p) / area
	l1 := triangleSurfaceArea(p2, p0, p) / area
	l2 := triangleSurfaceArea(p0, p1, p) / area
	return l0, l1, l2
}

func (t Triangle) Barycentric(p Vector) (float32, float32, float32) {
	return barycentric(t.P0, t.P1, t.P2, p)
}

func trianglesToSimd(t1, t2, t3, t4 Triangle) ([4]float32, [4]float32, [4]float32, [4]float32, [4]float32, [4]float32, [4]float32, [4]float32, [4]float32) {
	p0x := [4]float32{t1.P0.X, t2.P0.X, t3.P0.X, t4.P0.X}
	p0y := [4]float32{t1.P0.Y, t2.P0.Y, t3.P0.Y, t4.P0.Y}
	p0z := [4]float32{t1.P0.Z, t2.P0.Z, t3.P0.Z, t4.P0.Z}
	p1x := [4]float32{t1.P1.X, t2.P1.X, t3.P1.X, t4.P1.X}
	p1y := [4]float32{t1.P1.Y, t2.P1.Y, t3.P1.Y, t4.P1.Y}
	p1z := [4]float32{t1.P1.Z, t2.P1.Z, t3.P1.Z, t4.P1.Z}
	p2x := [4]float32{t1.P2.X, t2.P2.X, t3.P2.X, t4.P2.X}
	p2y := [4]float32{t1.P2.Y, t2.P2.Y, t3.P2.Y, t4.P2.Y}
	p2z := [4]float32{t1.P2.Z, t2.P2.Z, t3.P2.Z, t4.P2.Z}
	return p0x, p0y, p0z, p1x, p1y, p1z, p2x, p2y, p2z
}
