package model

import "math/rand"

// TODO: world to object coordinates and vice versa
// I think its only needed when caching common ray-object intersections?
// But I dont understand transformations well enough yet

// So pbrt uses 'primitive' for geometric primitives and complex objects both
// For me a primitive is always a geometric primitive, otherwise we'll talk about objects
// Primitives on their own can be objects, these are simple objects
// Complex objects consist of other objects. This includes meshes.
// Shared objects are objects referenced by multiple instances, and therefore need
// transformations from object space to world space (?)

type Object interface {
	Intersect(ray Ray) (si *SurfaceInteraction, hit bool)
	SurfaceNormal(point Vector) Vector
	GetColor(si *SurfaceInteraction) Color
	GetMaterial() Material
	SampleDirection(*rand.Rand, Vector) Vector
	IsLight() bool
	Bound(Transform) AABB
}

type object struct {
	Material Material
}

// to be replaced by more interesting materials info
func (o object) GetColor(si *SurfaceInteraction) Color {
	return o.Material.GetColor(si)
}

func (o object) SampleDirection(r *rand.Rand, normal Vector) Vector {
	return o.Material.Sample(r, normal)
}

func (o object) GetMaterial() Material {
	return o.Material
}

func (o object) IsLight() bool {
	return o.Material.IsLight()
}

func ObjectsBound(objects []Object, t Transform) AABB {
	bound := objects[0].Bound(t)
	for i := 1; i < len(objects); i++ {
		bound = bound.AddAABB(objects[i].Bound(t))
	}
	return bound
}

// in general, a complex object is just a set of objects.
// I use a bvh for the implementation
// SurfaceNormal and GetColor are part of later material functions
// these should always be called on the simple object that is hit,
// never on the aggregate object containing those (it doesnt have its own)
type ComplexObject struct {
	as AccelerationStructure
}

func NewComplexObject(objects []Object) Object {
	if len(objects) == 0 {
		panic("invalid objects, cant be empty")
	}
	return &ComplexObject{
		as: NewBVH(objects, SplitSurfaceAreaHeuristic),
	}
}

func NewTriangleComplexObject(triangles []Triangle) Object {
	if len(triangles) == 0 {
		panic("invalid triangles, cant be empty")
	}
	if SIMD_ENABLED {
		return &ComplexObject{
			as: NewTriangle4BVH(triangles),
		}
	}
	return &ComplexObject{
		as: NewTriangleBVH(triangles, SplitSurfaceAreaHeuristic),
	}
}

func (co *ComplexObject) Intersect(ray Ray) (*SurfaceInteraction, bool) {
	return co.as.ClosestIntersection(ray, MAX_RAY_DISTANCE)
}

func (co *ComplexObject) SurfaceNormal(point Vector) Vector {
	panic("Dont call this function!")
	return Vector{}
}

func (co *ComplexObject) GetColor(*SurfaceInteraction) Color {
	panic("Dont call this function!")
	return Color{}
}

func (co *ComplexObject) GetMaterial() Material {
	panic("Dont call this function!")
	return nil
}

func (co *ComplexObject) SampleDirection(r *rand.Rand, normal Vector) Vector {
	panic("Dont call this function!")
	return Vector{}
}

func (co *ComplexObject) IsLight() bool {
	panic("Dont call this function!")
	return false
}

// TODO: a prime candidate for caching
func (co *ComplexObject) Bound(t Transform) AABB {
	return ObjectsBound(co.as.GetObjects(), t)
}

func (co *ComplexObject) Objects() []Object {
	return co.as.GetObjects()
}

// a shared object stores a pointer to an object (type)
// and a transformation that places this instance of the object (token) in the scene.
// SurfaceNormal and GetColor are part of later material functions
// these should always be called on the simple object that is hit,
// never on the aggregate object containing those (it doesnt have its own)
// TODO: multiple instances can share geometry but differ in material? optionally?
type SharedObject struct {
	Object        Object
	ObjectToWorld Transform
	WorldToObject Transform
}

// o is the object being shared, originToPosition is the transform in
// world space from origin to the object's position
// Note: objects should be centered on origin or this will not work properly!
func NewSharedObject(o Object, originToPosition Transform) Object {
	// TODO: investigate: doubly shared objects?
	return &SharedObject{
		Object:        o,
		ObjectToWorld: originToPosition,
		WorldToObject: originToPosition.Inverse(),
	}
}

func (so *SharedObject) Intersect(ray Ray) (*SurfaceInteraction, bool) {
	// transform ray to object space
	r := so.WorldToObject.Ray(ray)
	si, ok := so.Object.Intersect(r)
	if !ok {
		return nil, false
	}
	// transform surface interaction info back to world space
	si.Point = so.ObjectToWorld.Point(si.UntransformedPoint)
	si.normal = so.ObjectToWorld.Normal(si.UntransformedNormal)
	return si, true
}

func (so *SharedObject) SurfaceNormal(Vector) Vector {
	panic("Dont call this function!")
	return Vector{}
}

func (so *SharedObject) GetColor(*SurfaceInteraction) Color {
	panic("Dont call this function!")
	return Color{}
}

func (so *SharedObject) GetMaterial() Material {
	panic("Dont call this function!")
	return nil
}

func (so *SharedObject) SampleDirection(r *rand.Rand, normal Vector) Vector {
	panic("Dont call this function!")
	return Vector{}
}

func (so *SharedObject) IsLight() bool {
	panic("Dont call this function!")
	return false
}

func (so *SharedObject) Bound(t Transform) AABB {
	transform := t.Mul(so.ObjectToWorld)
	return so.Object.Bound(transform)
}
