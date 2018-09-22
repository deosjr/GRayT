package model

// TODO: optimizations
// - mesh is not a mesh at all atm,
// just a collection of triangles.
// - speed: SIMD instructions

type TriangleMesh struct {
	pointGrid    [][]Vector
	xSize, ySize int
}

type TriangleInMesh struct {
	object
	x, y int
	mesh *TriangleMesh
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
	return VectorFromTo(p0, p2).Cross(VectorFromTo(p0, p1)).Normalize().Times(1)
}

func (t TriangleInMesh) points() (p0, p1, p2 Vector) {
	return t.mesh.get(t.x, t.y)
}
func (t TriangleInMesh) Bound(transform Transform) AABB {
	p0, p1, p2 := t.points()
	return triangleBound(p0, p1, p2, transform)
}
func (t TriangleInMesh) Intersect(r Ray) *hit {
	p0, p1, p2 := t.points()
	d, ok := triangleIntersect(p0, p1, p2, r)
	if !ok {
		return nil
	}
	hit := NewHit(t, d)
	hit.normal = t.SurfaceNormal(PointFromRay(r, d))
	return hit
}
func (t TriangleInMesh) SurfaceNormal(Vector) Vector {
	p0, p1, p2 := t.points()
	return triangleSurfaceNormal(p0, p1, p2)
}

// Assumption: for now this is a fully filled grid mesh
// so it doesnt represent things like cubes well (which wrap)
func NewTriangleMesh(grid [][]Vector) []Object {
	mesh := &TriangleMesh{
		pointGrid: grid,
		xSize:     len(grid[0]),
		ySize:     len(grid),
	}
	triangles := make([]Object, 2*(mesh.xSize-1)*(mesh.ySize-1))
	i := 0

	c1 := &DiffuseMaterial{NewColor(50, 200, 0)}
	c2 := &DiffuseMaterial{NewColor(50, 150, 50)}

	for y := 0; y < mesh.ySize-1; y++ {
		for x := 0; x < 2*(mesh.xSize-1); x += 2 {
			triangles[i] = TriangleInMesh{
				mesh:   mesh,
				x:      x,
				y:      y,
				object: object{c1},
			}
			i++
			triangles[i] = TriangleInMesh{
				mesh:   mesh,
				x:      x + 1,
				y:      y,
				object: object{c2},
			}
			i++
		}
	}
	return triangles
}

//   P1       P2
//    . ----- .
//    |    /  |
//    |  /    |
//    . ----- .
//   P4       P3
func (m *TriangleMesh) get(x, y int) (p0, p1, p2 Vector) {
	oldX := x
	x = x / 2
	if oldX%2 == 0 {
		// P1, P4, P2
		return m.pointGrid[y][x], m.pointGrid[y+1][x], m.pointGrid[y][x+1]
	}
	// P4, P3, P2
	return m.pointGrid[y+1][x], m.pointGrid[y+1][x+1], m.pointGrid[y][x+1]
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
func (t Triangle) Intersect(r Ray) *hit {
	d, ok := triangleIntersect(t.P0, t.P1, t.P2, r)
	if !ok {
		return nil
	}
	hit := NewHit(t, d)
	hit.normal = t.SurfaceNormal(PointFromRay(r, d))
	return hit
}
func (t Triangle) SurfaceNormal(Vector) Vector {
	return triangleSurfaceNormal(t.P0, t.P1, t.P2)
}
