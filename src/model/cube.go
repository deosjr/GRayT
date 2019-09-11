package model

import "math"

// Axis-aligned bounding box
type AABB struct {
	Pmin, Pmax Vector
}

func NewAABB(p1, p2 Vector) AABB {
	return AABB{
		Pmin: Vector{
			math.Min(p1.X, p2.X),
			math.Min(p1.Y, p2.Y),
			math.Min(p1.Z, p2.Z)},
		Pmax: Vector{
			math.Max(p1.X, p2.X),
			math.Max(p1.Y, p2.Y),
			math.Max(p1.Z, p2.Z)},
	}
}

func (b AABB) AddPoint(p Vector) AABB {
	return AABB{
		Pmin: Vector{
			math.Min(b.Pmin.X, p.X),
			math.Min(b.Pmin.Y, p.Y),
			math.Min(b.Pmin.Z, p.Z)},
		Pmax: Vector{
			math.Max(b.Pmax.X, p.X),
			math.Max(b.Pmax.Y, p.Y),
			math.Max(b.Pmax.Z, p.Z)},
	}
}

func (b AABB) AddAABB(b2 AABB) AABB {
	return AABB{
		Pmin: Vector{
			math.Min(b.Pmin.X, b2.Pmin.X),
			math.Min(b.Pmin.Y, b2.Pmin.Y),
			math.Min(b.Pmin.Z, b2.Pmin.Z)},
		Pmax: Vector{
			math.Max(b.Pmax.X, b2.Pmax.X),
			math.Max(b.Pmax.Y, b2.Pmax.Y),
			math.Max(b.Pmax.Z, b2.Pmax.Z)},
	}
}

func (b AABB) Centroid() Vector {
	return b.Pmin.Add(b.Pmax).Times(0.5)
}

// Unoptimised, analytic solution for now
// TODO: either optimise or check assumptions on NaN / divide by zero (inf) logic
// This function is one of the main bottlenecks
// writing the dimension loop explicitly saves a lot of time
func (b AABB) Intersect(ray Ray) (tMin float64, hit bool) {
	t0, t1 := 0.0, math.MaxFloat64
	invRayDir := 1 / ray.Direction.X
	tNear := (b.Pmin.X - ray.Origin.X) * invRayDir
	tFar := (b.Pmax.X - ray.Origin.Y) * invRayDir
	if tNear > tFar {
		tNear, tFar = tFar, tNear
	}
	// TODO: correct for error margin in tFar
	if tNear > t0 {
		t0 = tNear
	}
	if tFar < t1 {
		t1 = tFar
	}
	if t0 > t1 {
		return 0, false
	}
	invRayDir = 1 / ray.Direction.Y
	tNear = (b.Pmin.Y - ray.Origin.Y) * invRayDir
	tFar = (b.Pmax.Y - ray.Origin.Y) * invRayDir
	if tNear > tFar {
		tNear, tFar = tFar, tNear
	}
	// TODO: correct for error margin in tFar
	if tNear > t0 {
		t0 = tNear
	}
	if tFar < t1 {
		t1 = tFar
	}
	if t0 > t1 {
		return 0, false
	}
	invRayDir = 1 / ray.Direction.Z
	tNear = (b.Pmin.Z - ray.Origin.Z) * invRayDir
	tFar = (b.Pmax.Z - ray.Origin.Z) * invRayDir
	if tNear > tFar {
		tNear, tFar = tFar, tNear
	}
	// TODO: correct for error margin in tFar
	if tNear > t0 {
		t0 = tNear
	}
	if tFar < t1 {
		t1 = tFar
	}
	if t0 > t1 {
		return 0, false
	}
	return t0, true
}

// TODO: check efficiency
func (b AABB) MaximumExtent() Dimension {
	xExtent := b.Pmax.X - b.Pmin.X
	yExtent := b.Pmax.Y - b.Pmin.Y
	zExtent := b.Pmax.Z - b.Pmin.Z
	switch {
	case xExtent >= yExtent && xExtent >= zExtent:
		return X
	case yExtent >= xExtent && yExtent >= zExtent:
		return Y
	}
	return Z
}

// any cuboid is just an axis-aligned box with a rotation
type Cuboid struct {
	cuboid   AABB
	material Material
}

func NewCuboid(aabb AABB, m Material) Cuboid {
	return Cuboid{
		cuboid:   aabb,
		material: m,
	}
}

// Builds on quadrilateral definition:
// Let P1 - P4 be the top and
// let P5 - P8 be the bottom quadrilateral
// P1 corresponding to P5 etc

func (c Cuboid) Tesselate() Object {
	pmin := c.cuboid.Pmin
	pmax := c.cuboid.Pmax

	p1 := Vector{pmin.X, pmax.Y, pmin.Z}
	p2 := Vector{pmax.X, pmax.Y, pmin.Z}
	p3 := Vector{pmax.X, pmax.Y, pmax.Z}
	p4 := Vector{pmin.X, pmax.Y, pmax.Z}

	p5 := Vector{pmin.X, pmin.Y, pmin.Z}
	p6 := Vector{pmax.X, pmin.Y, pmin.Z}
	p7 := Vector{pmax.X, pmin.Y, pmax.Z}
	p8 := Vector{pmin.X, pmin.Y, pmax.Z}

	triangles := make([]Object, 12)
	triangles[0], triangles[1] = QuadrilateralToTriangles(p1, p2, p3, p4, c.material)
	triangles[2], triangles[3] = QuadrilateralToTriangles(p2, p1, p5, p6, c.material)
	triangles[4], triangles[5] = QuadrilateralToTriangles(p3, p2, p6, p7, c.material)
	triangles[6], triangles[7] = QuadrilateralToTriangles(p4, p3, p7, p8, c.material)
	triangles[8], triangles[9] = QuadrilateralToTriangles(p1, p4, p8, p5, c.material)
	triangles[10], triangles[11] = QuadrilateralToTriangles(p6, p5, p8, p7, c.material)
	return NewComplexObject(triangles)
}

//   P1           P2
//    . --------- .
//    |           |
//    |           |
//    |           |
//    |           |
//    . --------- .
//   P4           P3

type Quadrilateral struct {
	P1, P2, P3, P4 Vector
	material       Material
}

func NewQuadrilateral(p1, p2, p3, p4 Vector, m Material) Quadrilateral {
	return Quadrilateral{
		P1:       p1,
		P2:       p2,
		P3:       p3,
		P4:       p4,
		material: m,
	}
}

func (q Quadrilateral) Tesselate() (Triangle, Triangle) {
	t1, t2 := QuadrilateralToTriangles(q.P1, q.P2, q.P3, q.P4, q.material)
	return t1, t2
}

func QuadrilateralToTriangles(p1, p2, p3, p4 Vector, m Material) (Triangle, Triangle) {
	return NewTriangle(p1, p2, p4, m),
		NewTriangle(p4, p2, p3, m)
}
