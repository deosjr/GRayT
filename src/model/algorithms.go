package model

import "math"

const standardAlbedo = 0.18

type hit struct {
	object   *Object
	ray      Ray
	distance float64
}

func NewHit(o Object, r Ray, d float64) *hit {
	return &hit{
		object:   &o,
		ray:      r,
		distance: d,
	}
}

func LightContribution(ray Ray, hit *hit, l Light, as AccelerationStructure) (Color, bool) {
	point := PointFromRay(hit.ray, hit.distance)
	segment := VectorFromTo(point, l.Origin())
	shadowRay := NewRay(point, segment)
	segmentLength := segment.Length()
	if pointInShadow(shadowRay, as, segmentLength) {
		return Color{}, false
	}
	object := *hit.object
	facingRatio := object.SurfaceNormal(point).Dot(VectorFromTo(point, ray.Origin))
	if facingRatio <= 0 {
		return Color{}, false
	}
	lightRatio := object.SurfaceNormal(point).Dot(segment)
	factors := standardAlbedo / math.Pi * l.Intensity(segmentLength) * facingRatio * lightRatio
	lightColor := l.Color().Times(factors)
	return object.GetColor().Product(lightColor), true
}

func pointInShadow(shadowRay Ray, as AccelerationStructure, maxDistance float64) bool {
	if hit := as.ClosestIntersection(shadowRay, maxDistance); hit != nil {
		return true
	}
	return false
}
