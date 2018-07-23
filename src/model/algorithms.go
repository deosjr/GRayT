package model

import "math"

const standardAlbedo = 0.18

// TODO: world to object coordinates and vice versa
// I think its only needed when caching common ray-object intersections?
// But I dont understand transformations well enough yet

// So PBR uses 'primitive' for geometric primitives and complex objects both
// For me a primitive is always a geometric primitive, otherwise we'll talk about objects
// Primitives on their own can be objects, these are simple objects
// Complex objects consist of other objects. This includes meshes.
// Shared objects are objects referenced by multiple instances, and therefore need
// transformations from object space to world space (?)

type Object interface {
	Intersect(Ray) (distance float64, ok bool)
	SurfaceNormal(point Vector) Vector
	GetColor() Color
	Bound() AABB
}

type object struct {
	Color Color
}

func (o object) GetColor() Color {
	return o.Color
}

func LightContribution(ray Ray, hit *hit, l Light, as AccelerationStructure) (Color, bool) {
	segment := VectorFromTo(hit.point, l.Origin())
	shadowRay := NewRay(hit.point, segment)
	segmentLength := segment.Length()
	if pointInShadow(shadowRay, as, segmentLength) {
		return Color{}, false
	}
	facingRatio := hit.object.SurfaceNormal(hit.point).Dot(VectorFromTo(hit.point, ray.Origin))
	if facingRatio <= 0 {
		return Color{}, false
	}
	lightRatio := hit.object.SurfaceNormal(hit.point).Dot(segment)
	factors := standardAlbedo / math.Pi * l.Intensity(segmentLength) * facingRatio * lightRatio
	lightColor := l.Color().Times(factors)
	return hit.object.GetColor().Product(lightColor), true
}

func pointInShadow(shadowRay Ray, as AccelerationStructure, maxDistance float64) bool {
	if hit := as.ClosestIntersection(shadowRay, maxDistance); hit != nil {
		return true
	}
	return false
}
