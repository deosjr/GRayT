package model

import (
	"math"
)

var BACKGROUND_COLOR = NewColor(10, 10, 10)

type question struct {
	scene *Scene
	x, y  int
}

type answer struct {
	x, y  int
	color Color
}

func worker(ch chan question, ans chan answer) {
	for q := range ch {
		ray := q.scene.Camera.PixelRay(q.x, q.y)

		var objectHit Object
		var intersection *Vector
		var minDistance = math.MaxFloat64
		for _, o := range q.scene.Objects {
			if i, ok, distance := o.Intersect(ray); ok && distance >= 0 && distance < minDistance {
				minDistance = distance
				intersection = &i
				objectHit = o
			}
		}

		if intersection == nil {
			ans <- answer{q.x, q.y, BACKGROUND_COLOR}
			continue
		}

		color := NewColor(0, 0, 0)

		for _, l := range q.scene.Lights {
			segment := VectorFromTo(*intersection, l.Origin())
			shadowRay := NewRay(*intersection, segment)
			segmentLength := segment.Length()
			if pointInShadow(shadowRay, q.scene.Objects, segmentLength) {
				continue
			}
			facingRatio := objectHit.SurfaceNormal(*intersection).Dot(segment)
			if facingRatio <= 0 {
				continue
			}
			color = color.Add(l.Color().Times(STANDARD_ALBEDO / math.Pi * l.Intensity(segmentLength) * facingRatio))
		}

		ans <- answer{q.x, q.y, color}
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
		if _, ok, distance := o.Intersect(shadowRay); ok && distance > e && distance < maxDistance {
			return true
		}
	}
	return false
}
