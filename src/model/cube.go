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
func (b AABB) Intersect(ray Ray) bool {
	t0, t1 := 0.0, math.MaxFloat64
	for _, dim := range Dimensions {
		// TODO: divide by zero?
		invRayDir := 1 / ray.Direction.Get(dim)
		tNear := (b.Pmin.Get(dim) - ray.Origin.Get(dim)) * invRayDir
		tFar := (b.Pmax.Get(dim) - ray.Origin.Get(dim)) * invRayDir
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
			return false
		}
	}
	return true
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

// Builds on quadrilateral definition:
// Let P1 - P4 be the top and
// let P5 - P8 be the bottom quadrilateral
// P1 corresponding to P5 etc

// TODO: Can be represented with only 2 points, like AABB
// This is actually a crude start at a mesh, come to think of it

type Cuboid struct {
	P1, P2, P3, P4, P5, P6, P7, P8 Vector
	Color                          Color
}

func (c Cuboid) Tesselate() []Object {
	F1 := Quadrilateral{c.P1, c.P2, c.P3, c.P4, c.Color}.Tesselate()
	F2 := Quadrilateral{c.P2, c.P1, c.P5, c.P6, c.Color}.Tesselate()
	F3 := Quadrilateral{c.P3, c.P2, c.P6, c.P7, c.Color}.Tesselate()
	F4 := Quadrilateral{c.P4, c.P3, c.P7, c.P8, c.Color}.Tesselate()
	F5 := Quadrilateral{c.P1, c.P4, c.P8, c.P5, c.Color}.Tesselate()
	F6 := Quadrilateral{c.P6, c.P5, c.P8, c.P7, c.Color}.Tesselate()
	return append(F1, append(F2, append(F3, append(F4, append(F5, F6...)...)...)...)...)
}

//   P1           P2
//    . --------- .
//    |           |
//    |           |
//    |           |
//    |           |
//    . --------- .
//   P4           P3

// cant represent in 2 points due to sidedness?

type Quadrilateral struct {
	P1, P2, P3, P4 Vector
	Color          Color
}

func (r Quadrilateral) Tesselate() []Object {
	return []Object{
		NewTriangle(r.P1, r.P2, r.P4, r.Color),
		NewTriangle(r.P4, r.P2, r.P3, r.Color),
	}
}
