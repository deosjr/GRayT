package model

import "math"

var BACKGROUND_COLOR = NewColor(10, 10, 10)

type question struct {
	x, y int
}

type answer struct {
	x, y  int
	color Color
}

func worker(scene *Scene, ch chan question, ans chan answer) {
	for q := range ch {
		ray := scene.Camera.PixelRay(q.x, q.y)
		hit := closestIntersection(ray, scene.Objects)
		if hit == nil {
			ans <- answer{q.x, q.y, BACKGROUND_COLOR}
			continue
		}

		color := NewColor(0, 0, 0)
		for _, l := range scene.Lights {
			c, ok := lightContribution(ray, hit, color, l, scene.Objects)
			if !ok {
				continue
			}
			color = c
		}
		ans <- answer{q.x, q.y, color}
	}
}

type hit struct {
	object Object
	point  Vector
}

func closestIntersection(ray Ray, objects []Object) *hit {
	var objectHit Object
	d := math.MaxFloat64
	for _, o := range objects {
		if distance, ok := o.Intersect(ray); ok && distance < d {
			d = distance
			objectHit = o
		}
	}
	if d == math.MaxFloat64 {
		return nil
	}
	return &hit{
		object: objectHit,
		point:  PointFromRay(ray, d),
	}
}

func lightContribution(ray Ray, hit *hit, color Color, l Light, objects []Object) (Color, bool) {
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
	// TODO: this seems fishy. First light will contribute more color than second one?
	factors := STANDARD_ALBEDO / math.Pi * l.Intensity(segmentLength) * facingRatio * lightRatio
	lightColor := color.Add(l.Color().Times(factors))
	return hit.object.GetColor().Product(lightColor), true
}

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
