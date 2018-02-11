package model

import "math"

const standardAlbedo = 0.18

// TODO: world to object coordinates and vice versa
// I think its only needed when caching common ray-object intersections?
// But I dont understand transformations well enough yet

// TODO: The term Object is conflated with Primitive right now
// Objects should be complex (recognisable) objects, which
// can be split or tesselated to their primitives

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

func LightContribution(ray Ray, hit *hit, l Light, objects []Object) (Color, bool) {
	segment := VectorFromTo(hit.point, l.Origin())
	shadowRay := NewRay(hit.point, segment)
	segmentLength := segment.Length()
	if pointInShadow(shadowRay, objects, segmentLength) {
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

// TODO: use BVH here too
func pointInShadow(shadowRay Ray, objects []Object, maxDistance float64) bool {
	// floating point error margin
	// TODO: setting too small drops shadows completely?
	// setting to 0.1 or 0.5 shows shadows; setting too big gives weirdness
	// see https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/ligth-and-shadows
	// on shadow bias
	e := 1E-10
	for _, o := range objects {
		if distance, ok := o.Intersect(shadowRay); ok && distance > e && distance < maxDistance {
			return true
		}
	}
	return false
}
