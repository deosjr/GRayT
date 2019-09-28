package model

import "math"
import "math/rand"

// - speed: SIMD instructions

type TriangleMesh struct {
	vertices map[int64]Vector
}

type TriangleInMesh struct {
	object
	p0, p1, p2 int64
	mesh       *TriangleMesh
}

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
func triangleIntersect(p0, p1, p2 Vector, ray Ray) (float64, bool) {
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
	return VectorFromTo(p0, p2).Cross(VectorFromTo(p0, p1)).Normalize() //.Times(-1)
}

func (t TriangleInMesh) points() (p0, p1, p2 Vector) {
	p0 = t.mesh.get(t.p0)
	p1 = t.mesh.get(t.p1)
	p2 = t.mesh.get(t.p2)
	return
}
func (t TriangleInMesh) Bound(transform Transform) AABB {
	p0, p1, p2 := t.points()
	return triangleBound(p0, p1, p2, transform)
}
func (t TriangleInMesh) Intersect(r Ray) (*SurfaceInteraction, bool) {
	p0, p1, p2 := t.points()
	d, ok := triangleIntersect(p0, p1, p2, r)
	if !ok {
		return nil, false
	}
	n := triangleSurfaceNormal(p0, p1, p2)
	return NewSurfaceInteraction(t, d, n, r), true
}
func (t TriangleInMesh) IntersectOptimized(r Ray) (float64, bool) {
	p0, p1, p2 := t.points()
	d, ok := triangleIntersect(p0, p1, p2, r)
	if !ok {
		return 0, false
	}
	return d, true
}
func (t TriangleInMesh) SurfaceNormal(Vector) Vector {
	p0, p1, p2 := t.points()
	return triangleSurfaceNormal(p0, p1, p2)
}

type Face struct {
	V0, V1, V2 int64
}

func NewFace(v0, v1, v2 int64) Face {
	return Face{v0, v1, v2}
}

func NewTriangleMesh(vertices []Vector, faces []Face, m Material) Object {
	vertexMap := map[int64]Vector{}
	for i, v := range vertices {
		vertexMap[int64(i)] = v
	}
	mesh := &TriangleMesh{
		vertices: vertexMap,
	}
	triangles := make([]Object, len(faces))
	for i, f := range faces {
		triangles[i] = TriangleInMesh{
			object: object{m},
			p0:     f.V0,
			p1:     f.V1,
			p2:     f.V2,
			mesh:   mesh,
		}
	}
	return NewComplexObject(triangles)
}

func (m *TriangleMesh) get(i int64) Vector {
	return m.vertices[i]
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
func (t Triangle) IntersectOptimized(r Ray) (float64, bool) {
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
	a, b := random.Float64(), random.Float64()
	if a+b > 1 {
		a, b = 1-a, 1-b
	}
	return t.P0.Add(u.Times(a)).Add(v.Times(b))
}

// Heron's formula
func (t Triangle) Area() float64 {
	a := t.P1.Sub(t.P0).Length()
	b := t.P2.Sub(t.P0).Length()
	c := t.P1.Sub(t.P2).Length()
	s := (a + b + c) / 2
	return math.Sqrt(s * (s - a) * (s - b) * (s - c))
}
