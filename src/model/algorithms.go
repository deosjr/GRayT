package model

import "math"

const standardAlbedo = 0.18

type hit struct {
	object   Object
	ray      Ray
	distance float64
}

func LightContribution(ray Ray, hit *hit, l Light, as AccelerationStructure) (Color, bool) {
	point := PointFromRay(hit.ray, hit.distance)
	segment := VectorFromTo(point, l.Origin())
	shadowRay := NewRay(point, segment)
	segmentLength := segment.Length()
	if pointInShadow(shadowRay, as, segmentLength) {
		return Color{}, false
	}
	facingRatio := hit.object.SurfaceNormal(point).Dot(VectorFromTo(point, ray.Origin))
	if facingRatio <= 0 {
		return Color{}, false
	}
	lightRatio := hit.object.SurfaceNormal(point).Dot(segment)
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
