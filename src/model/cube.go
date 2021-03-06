package model

import (
    "image"
    "math"
)

// Axis-aligned bounding box
type AABB struct {
	Pmin, Pmax Vector
}

func NewAABB(p1, p2 Vector) AABB {
	return AABB{
		Pmin: VectorMin(p1, p2),
		Pmax: VectorMax(p1, p2),
	}
}

func (b AABB) AddPoint(p Vector) AABB {
	return AABB{
		Pmin: VectorMin(b.Pmin, p),
		Pmax: VectorMax(b.Pmax, p),
	}
}

func (b AABB) AddAABB(b2 AABB) AABB {
	return AABB{
		Pmin: VectorMin(b.Pmin, b2.Pmin),
		Pmax: VectorMax(b.Pmax, b2.Pmax),
	}
}

func (b AABB) Centroid() Vector {
	return b.Pmin.Add(b.Pmax).Times(0.5)
}

func (b AABB) Intersect(ray Ray) (tMin float32, hit bool) {
	var t0 float32 = 0.0
	var t1 float32 = math.MaxFloat32
	invRayDirs := [3]float32{1.0 / ray.Direction.X, 1.0 / ray.Direction.Y, 1.0 / ray.Direction.Z}
	rayOrigins := [3]float32{ray.Origin.X, ray.Origin.Y, ray.Origin.Z}
	bPmins := [3]float32{b.Pmin.X, b.Pmin.Y, b.Pmin.Z}
	bPmaxs := [3]float32{b.Pmax.X, b.Pmax.Y, b.Pmax.Z}

	for dim := 0; dim < 3; dim++ {
		tNear := (bPmins[dim] - rayOrigins[dim]) * invRayDirs[dim]
		tFar := (bPmaxs[dim] - rayOrigins[dim]) * invRayDirs[dim]
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
	}
	return t0, true
}

func (b AABB) SurfaceArea() float32 {
	d := b.Pmax.Sub(b.Pmin)
	return 2 * (d.X*d.Y + d.Y*d.Z + d.X*d.Z)
}

// used in surface area heuristic: maps a centroid's relative position
// to the box's min and max corners between 0 and 1
// assumption is that p lies within b
func (b AABB) Offset(p Vector) Vector {
	o := p.Sub(b.Pmin)
	// guarding against division by 0
	if b.Pmax.X > b.Pmin.X {
		o.X = o.X / (b.Pmax.X - b.Pmin.X)
	}
	if b.Pmax.Y > b.Pmin.Y {
		o.Y = o.Y / (b.Pmax.Y - b.Pmin.Y)
	}
	if b.Pmax.Z > b.Pmin.Z {
		o.Z = o.Z / (b.Pmax.Z - b.Pmin.Z)
	}
	return o
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

func boundsToSimd(b1, b2, b3, b4 AABB) ([4]float32, [4]float32, [4]float32, [4]float32, [4]float32, [4]float32) {
	min4x := [4]float32{b1.Pmin.X, b2.Pmin.X, b3.Pmin.X, b4.Pmin.X}
	min4y := [4]float32{b1.Pmin.Y, b2.Pmin.Y, b3.Pmin.Y, b4.Pmin.Y}
	min4z := [4]float32{b1.Pmin.Z, b2.Pmin.Z, b3.Pmin.Z, b4.Pmin.Z}
	max4x := [4]float32{b1.Pmax.X, b2.Pmax.X, b3.Pmax.X, b4.Pmax.X}
	max4y := [4]float32{b1.Pmax.Y, b2.Pmax.Y, b3.Pmax.Y, b4.Pmax.Y}
	max4z := [4]float32{b1.Pmax.Z, b2.Pmax.Z, b3.Pmax.Z, b4.Pmax.Z}
	return min4x, min4y, min4z, max4x, max4y, max4z
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

func (c Cuboid) Tesselate() []Triangle {
	pmin := c.cuboid.Pmin
	pmax := c.cuboid.Pmax

	p1 := Vector{pmin.X, pmax.Y, pmax.Z}
	p2 := Vector{pmax.X, pmax.Y, pmax.Z}
	p3 := Vector{pmax.X, pmax.Y, pmin.Z}
	p4 := Vector{pmin.X, pmax.Y, pmin.Z}

	p5 := Vector{pmin.X, pmin.Y, pmax.Z}
	p6 := Vector{pmax.X, pmin.Y, pmax.Z}
	p7 := Vector{pmax.X, pmin.Y, pmin.Z}
	p8 := Vector{pmin.X, pmin.Y, pmin.Z}

	triangles := make([]Triangle, 12)
	triangles[0], triangles[1] = QuadrilateralToTriangles(p1, p2, p3, p4, c.material)
	triangles[2], triangles[3] = QuadrilateralToTriangles(p2, p1, p5, p6, c.material)
	triangles[4], triangles[5] = QuadrilateralToTriangles(p3, p2, p6, p7, c.material)
	triangles[6], triangles[7] = QuadrilateralToTriangles(p4, p3, p7, p8, c.material)
	triangles[8], triangles[9] = QuadrilateralToTriangles(p1, p4, p8, p5, c.material)
	triangles[10], triangles[11] = QuadrilateralToTriangles(p6, p5, p8, p7, c.material)
	return triangles
}

func (c Cuboid) TesselateInsideOut() []Triangle {
	pmin := c.cuboid.Pmin
	pmax := c.cuboid.Pmax

	p1 := Vector{pmin.X, pmax.Y, pmax.Z}
	p2 := Vector{pmax.X, pmax.Y, pmax.Z}
	p3 := Vector{pmax.X, pmax.Y, pmin.Z}
	p4 := Vector{pmin.X, pmax.Y, pmin.Z}

	p5 := Vector{pmin.X, pmin.Y, pmax.Z}
	p6 := Vector{pmax.X, pmin.Y, pmax.Z}
	p7 := Vector{pmax.X, pmin.Y, pmin.Z}
	p8 := Vector{pmin.X, pmin.Y, pmin.Z}

	triangles := make([]Triangle, 12)
	triangles[0], triangles[1] = QuadrilateralToTriangles(p4, p3, p2, p1, c.material)
	triangles[2], triangles[3] = QuadrilateralToTriangles(p6, p5, p1, p2, c.material)
	triangles[4], triangles[5] = QuadrilateralToTriangles(p7, p6, p2, p3, c.material)
	triangles[6], triangles[7] = QuadrilateralToTriangles(p8, p7, p3, p4, c.material)
	triangles[8], triangles[9] = QuadrilateralToTriangles(p5, p8, p4, p1, c.material)
	triangles[10], triangles[11] = QuadrilateralToTriangles(p7, p8, p5, p6, c.material)
	return triangles
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
		NewTriangle(p2, p3, p4, m)
}

// Similar to cube.tesselate, but returns a triangle mesh with uv mapping
// img is mapped onto the cube as follows:
//     |---| 
//     | 2 |
// |---|---|---|---|
// | 1 | 3 | 5 | 6 |
// |---|---|---|---|
//     | 4 |
//     |---|
// 1 is the front of the cube, facing in +Z direction
// 2 is right, 3 is bottom, 4 is left, 5 is back, 6 is top
// each square has topleft at topleft, so oriented the same
// therefore NOT neatly wrapping the whole cross around the cube!
// works for arbitrary resolution as long as aspect ratio is 4:3
func CubeMesh(size float32, img image.Image) Object {
    /*
    */
    f := func(p0,p1,p2,p3 int64) (Face, Face) {
        return Face{p0, p2, p1}, Face{p1, p2, p3}
    }
    // unit cube centered around origin
    // TODO: should use ScaleUniform but cant get it to work properly yet
    span := size / 2.0
    min, max := -span, span

    // see ilkinulas.github.io/development/unity/2016/05/06/uv-mapping.html for details
	p0 := Vector{max, max, max}
	p1 := Vector{max, min, max}
	p2 := Vector{min, max, max}
	p3 := Vector{min, min, max}
	p4 := Vector{max, min, min}
	p5 := Vector{min, min, min}
	p6 := Vector{max, max, min}
	p7 := Vector{min, max, min}
    // duplicate vertices because they can be uv mapped differently
    vertices := []Vector{p0, p1, p2, p3, p4, p5, p6, p7, p0, p2, p0, p6, p2, p7}

	faces := make([]Face, 12)
	faces[0], faces[1] = f(0, 1, 2, 3)
	faces[2], faces[3] = f(10, 11, 1, 4)
	faces[4], faces[5] = f(1, 4, 3, 5)
	faces[6], faces[7] = f(3, 5, 12, 13)
	faces[8], faces[9] = f(4, 6, 5, 7)
	faces[10], faces[11] = f(6, 8, 7, 9)
    uvmap := map[int64]Vector{
        0:  Vector{0, 2.0/3.0, 0},
        1:  Vector{0.25, 2.0/3.0, 0},
        2:  Vector{0, 1.0/3.0, 0},
        3:  Vector{0.25, 1.0/3.0, 0},
        4:  Vector{0.5, 2.0/3.0, 0},
        5:  Vector{0.5, 1.0/3.0, 0},
        6:  Vector{0.75, 2.0/3.0, 0},
        7:  Vector{0.75, 1.0/3.0, 0},
        8:  Vector{1, 2.0/3.0, 0},
        9:  Vector{1, 1.0/3.0, 0},
        10: Vector{0.25, 1, 0},
        11: Vector{0.5, 1, 0},
        12: Vector{0.25, 0, 0},
        13: Vector{0.5, 0, 0},
    }

    mat := NewDiffuseMaterial(NewImageTexture(img, TriangleMeshUVFunc))
    obj := NewTriangleMesh(vertices, faces, mat)
    obj.(*TriangleMesh).UV = uvmap
    return obj
}
