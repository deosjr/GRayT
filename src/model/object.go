package model

import (
	"math"
)

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
	Intersect(Ray) *hit
	SurfaceNormal(point Vector) Vector
	GetColor() Color
	Bound() AABB
}

type object struct {
	Color Color
}

// to be replaced by more interesting materials info
func (o object) GetColor() Color {
	return o.Color
}

func ObjectsBound(objects []Object) AABB {
	bound := objects[0].Bound()
	for i := 1; i < len(objects); i++ {
		bound = bound.AddAABB(objects[i].Bound())
	}
	return bound
}

// in general, a complex object is just a set of objects.
// I use a bvh for the implementation
// SurfaceNormal and GetColor are part of later material functions
// these should always be called on the simple object that is hit,
// never on the aggregate object containing those (it doesnt have its own)
type ComplexObject struct {
	bvh   BVH
	bound AABB
}

func NewComplexObject(objects []Object) Object {
	if len(objects) == 0 {
		panic("invalid object, cant be empty")
	}
	return &ComplexObject{
		bvh:   NewBVH(objects, SplitMiddle),
		bound: ObjectsBound(objects),
	}
}

func (co *ComplexObject) Intersect(ray Ray) *hit {
	// TODO: not happy to redeclare max intersect distance here
	return co.bvh.ClosestIntersection(ray, math.MaxFloat64)
}

func (co *ComplexObject) SurfaceNormal(point Vector) Vector {
	panic("Dont call this function!")
	return Vector{}
}

func (co *ComplexObject) GetColor() Color {
	panic("Dont call this function!")
	return Color{}
}

func (co *ComplexObject) Bound() AABB {
	return co.bound
}

// a shared object stores a pointer to an object (type)
// and a transformation that places this instance of the object (token) in the scene.
// TODO: multiple instances can share geometry but differ in material? optionally?
type SharedObject struct {
	object        Object
	objectToWorld Transform
}

// o is the object being shared, originToPosition is the transform in
// world space from origin to the object's position
// Note: objects should be centered on origin or this will not work properly!
func NewSharedObject(o Object, originToPosition Transform) Object {
	if o.Bound().Centroid().Length() != 0 {
		panic("shared object should be centered on the origin!")
	}
	return &SharedObject{
		object:        o,
		objectToWorld: originToPosition,
	}
}

func (so *SharedObject) Intersect(ray Ray) *hit {
	// transform ray to object space
	r := so.objectToWorld.Inverse().Ray(ray)

	hit := so.object.Intersect(r)
	if hit == nil {
		return nil
	}

	// transform hit info back to world space
	hit.ray = so.objectToWorld.Ray(r)
	return hit
}

func (so *SharedObject) SurfaceNormal(point Vector) Vector {
	return so.object.SurfaceNormal(point)
}

func (so *SharedObject) GetColor() Color {
	return so.object.GetColor()
}

func (so *SharedObject) Bound() AABB {
	b := so.object.Bound()
	bmin := so.objectToWorld.Vector(b.Pmin)
	bmax := so.objectToWorld.Vector(b.Pmax)
	return NewAABB(bmin, bmax)
}
