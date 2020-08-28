package model

import "math/rand"

// TODO: mesh can more optimally be used as SharedObject, since the transform
// can be precalculated. From PBRT Ch3.6:
//
//   Unlike the other shapes that leave the shape description in object space and 
//   then transform incoming rays from world space to object space, 
//   triangle meshes transform the shape into world space and thus save the work of 
//   transforming incoming rays into object space and the work of transforming the intersection’s 
//   geometric representation out to world space. 
//   This is a good idea because this operation can be performed once at start-up, 
//   avoiding transforming rays many times during rendering.

// start at mesh as a special case of complex object
// face-vertex mesh, matches .obj format description
type TriangleMesh struct {
    object
    as AccelerationStructure
	vertices map[int64]Vector
    // vertex normals used for interpolated normal mapping
    Normals map[int64]Vector
    // TODO: 2d vector instead of Vector?
    // u and v values associated to vertices, if any
    UV map[int64]Vector
}

// NOTE: the mesh is the object inheriting material, not the triangle
type TriangleInMesh struct {
	p0, p1, p2 int64
	Mesh       *TriangleMesh
}

type Face struct {
	V0, V1, V2 int64
}

func NewFace(v0, v1, v2 int64) Face {
	return Face{v0, v1, v2}
}

func NewTriangleMesh(vertices []Vector, faces []Face, mat Material) Object {
	vertexMap := map[int64]Vector{}
	for i, v := range vertices {
		vertexMap[int64(i)] = v
	}
	mesh := &TriangleMesh{
        object: object{mat},
		vertices: vertexMap,
	}
	triangles := make([]Object, len(faces))
	for i, f := range faces {
		triangles[i] = TriangleInMesh{
			p0:     f.V0,
			p1:     f.V1,
			p2:     f.V2,
			Mesh:   mesh,
		}
	}
    // TODO: make this work in triangleBVH (triangle interface? slower...)
    mesh.as = NewBVH(triangles, SplitSurfaceAreaHeuristic)
	return mesh
}

// input is a list of n+1 x m+1 points describing a rectangular mesh
// of n x m triangles, fully connected (but not circular)
func NewGridTriangleMesh(n, m int, vertices, normals, uvs []Vector, mat Material) Object {
    if len(vertices) != (n+1)*(m+1) {
        panic("incorrect number of vertices to mesh")
    }
	vertexMap := map[int64]Vector{}
	for i, v := range vertices {
		vertexMap[int64(i)] = v
	}
    var normalMap map[int64]Vector
    if len(normals) > 0 {
        normalMap = map[int64]Vector{}
        for i, v := range normals {
            normalMap[int64(i)] = v
        }
    }
    var uvMap map[int64]Vector
    if len(uvs) > 0 {
        uvMap = map[int64]Vector{}
        for i, v := range uvs {
            uvMap[int64(i)] = v
        }
    }
	mesh := &TriangleMesh{
        object:   object{mat},
		vertices: vertexMap,
        Normals:  normalMap,
        UV:       uvMap,
	}
    triangles := []Object{}
    for y:=0; y<m; y++ {
        for x:=0; x<n; x++ {
            llhc := int64(y*(m+1)+x)
            lrhc := int64(y*(m+1)+x+1)
            ulhc := int64((y+1)*(m+1)+x)
            urhc := int64((y+1)*(m+1)+x+1)
            t1 := TriangleInMesh{
                p0:   llhc,
                p1:   lrhc,
                p2:   ulhc,
                Mesh: mesh,
            }
            triangles = append(triangles, t1)
            t2 := TriangleInMesh{
                p0:   lrhc,
                p1:   urhc,
                p2:   ulhc,
                Mesh: mesh,
            }
            triangles = append(triangles, t2)
        }
    }
    // TODO: make this work in triangleBVH (triangle interface? slower...)
    mesh.as = NewBVH(triangles, SplitSurfaceAreaHeuristic)
	return mesh
}

func (m *TriangleMesh) Intersect(ray Ray) (*SurfaceInteraction, bool) {
	return m.as.ClosestIntersection(ray, MAX_RAY_DISTANCE)
}

// TODO: a prime candidate for caching
func (m *TriangleMesh) Bound(t Transform) AABB {
	return ObjectsBound(m.as.GetObjects(), t)
}

func (m *TriangleMesh) SurfaceNormal(point Vector) Vector {
	panic("Dont call this function!")
	return Vector{}
}

func (m *TriangleMesh) get(i int64) Vector {
	return m.vertices[i]
}
func (t TriangleInMesh) GetColor(si *SurfaceInteraction) Color {
    return t.Mesh.GetColor(si)
}

func (t TriangleInMesh) GetMaterial() Material {
    return t.Mesh.GetMaterial()
}

func (t TriangleInMesh) SampleDirection(r *rand.Rand, normal Vector) Vector {
    return t.Mesh.SampleDirection(r, normal)
}

func (t TriangleInMesh) IsLight() bool {
    return t.Mesh.IsLight()
}

func (t TriangleInMesh) Points() (p0, p1, p2 Vector) {
	p0 = t.Mesh.get(t.p0)
	p1 = t.Mesh.get(t.p1)
	p2 = t.Mesh.get(t.p2)
	return
}

func (t TriangleInMesh) PointIndices() (int64, int64, int64) {
    return t.p0, t.p1, t.p2
}

func (t TriangleInMesh) Bound(transform Transform) AABB {
	p0, p1, p2 := t.Points()
	return triangleBound(p0, p1, p2, transform)
}

func (t TriangleInMesh) Intersect(r Ray) (*SurfaceInteraction, bool) {
	p0, p1, p2 := t.Points()
	d, ok := triangleIntersect(p0, p1, p2, r)
	if !ok {
		return nil, false
	}
	n := triangleSurfaceNormal(p0, p1, p2)
	return NewSurfaceInteraction(t, d, n, r), true
}

func (t TriangleInMesh) IntersectOptimized(r Ray) (float32, bool) {
	p0, p1, p2 := t.Points()
	d, ok := triangleIntersect(p0, p1, p2, r)
	if !ok {
		return 0, false
	}
	return d, true
}

func (t TriangleInMesh) SurfaceNormal(Vector) Vector {
	p0, p1, p2 := t.Points()
	return triangleSurfaceNormal(p0, p1, p2)
}

func (t TriangleInMesh) Barycentric(p Vector) (float32, float32, float32) {
	p0, p1, p2 := t.Points()
    return barycentric(p0, p1, p2, p)
}
