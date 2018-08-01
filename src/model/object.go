package model

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
	GetColor(si *SurfaceInteraction) Color
	Bound(Transform) AABB
}

type object struct {
	Material Material
}

// to be replaced by more interesting materials info
func (o object) GetColor(si *SurfaceInteraction) Color {
	return o.Material.GetColor(si)
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
	bvh   *BVH
}

func NewComplexObject(objects []Object) Object {
	if len(objects) == 0 {
		panic("invalid object, cant be empty")
	}
	return &ComplexObject{
		bvh:   NewBVH(objects, SplitMiddle),
	}
}

func (co *ComplexObject) Intersect(ray Ray) *hit {
	return co.bvh.ClosestIntersection(ray, MAX_RAY_DISTANCE)
}

func (co *ComplexObject) SurfaceNormal(point Vector) Vector {
	panic("Dont call this function!")
	return Vector{}
}

func (co *ComplexObject) GetColor(*SurfaceInteraction) Color {
	panic("Dont call this function!")
	return Color{}
}

// TODO: a prime candidate for caching
func (co *ComplexObject) Bound(t Transform) AABB {
	return ObjectsBound(co.bvh.objects, t)
}

// a shared object stores a pointer to an object (type)
// and a transformation that places this instance of the object (token) in the scene.
// SurfaceNormal and GetColor are part of later material functions
// these should always be called on the simple object that is hit,
// never on the aggregate object containing those (it doesnt have its own)
// TODO: multiple instances can share geometry but differ in material? optionally?
type SharedObject struct {
	object        Object
	objectToWorld Transform
}

// o is the object being shared, originToPosition is the transform in
// world space from origin to the object's position
// Note: objects should be centered on origin or this will not work properly!
func NewSharedObject(o Object, originToPosition Transform) Object {
	// TODO: investigate: doubly shared objects?
	if o.Bound(identity).Centroid().Length() != 0 {
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
	point := PointFromRay(r, hit.distance)
	hit.normal = so.objectToWorld.Normal(hit.object.SurfaceNormal(point))
	return hit
}

func (so *SharedObject) SurfaceNormal(Vector) Vector {
	panic("Dont call this function!")
	return Vector{}
}

func (so *SharedObject) GetColor(*SurfaceInteraction) Color {
	panic("Dont call this function!")
	return Color{}
}

func (so *SharedObject) Bound(t Transform) AABB {
	transform := t.Mul(so.objectToWorld)
	return so.object.Bound(transform)
}
