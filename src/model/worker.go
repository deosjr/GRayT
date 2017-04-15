package model

import (
	"math"
)

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
			segment := VectorFromTo(hit.point, l.Origin())
			shadowRay := NewRay(hit.point, segment)
			segmentLength := segment.Length()
			if pointInShadow(shadowRay, scene.Objects, segmentLength) {
				continue
			}
			facingRatio := hit.object.SurfaceNormal(hit.point).Dot(segment)
			if facingRatio <= 0 {
				continue
			}
			color = color.Add(l.Color().Times(STANDARD_ALBEDO / math.Pi * l.Intensity(segmentLength) * facingRatio))
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
		if distance, ok := o.Intersect(ray); ok && distance >= 0 && distance < d {
			d = distance
			objectHit = o
		}
	}
	if d == math.MaxFloat64 {
		return nil
	}
	return &hit{
		object: objectHit,
		point:  ray.Origin.Add(ray.Direction.Times(d)),
	}
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
