package model

import (
	"math"
)

var BACKGROUND_COLOR = Color{Vector{10, 10, 10}}

type Question struct {
	Scene *Scene
	X, Y  int
}

type Answer struct {
	X, Y  int
	Color Color
}

func Worker(ch chan Question, ans chan Answer) {
	for q := range ch {
		ray := q.Scene.Camera.PixelRay(q.X, q.Y)

		var objectHit Object
		var intersection *Vector
		var minDistance = math.MaxFloat64
		for _, o := range q.Scene.Objects {
			if i, ok, distance := o.Intersect(ray); ok && distance >= 0 && distance < minDistance {
				minDistance = distance
				intersection = &i
				objectHit = o
			}
		}

		if intersection == nil {
			ans <- Answer{q.X, q.Y, BACKGROUND_COLOR}
			continue
		}

		color := Color{Vector{0, 0, 0}}

	Lights:
		for _, l := range q.Scene.Lights {
			segment := VectorFromTo(*intersection, l.Origin)
			shadow := NewRay(*intersection, segment)
			segmentLength := segment.Length()

			// floating point error margin
			// TODO: setting too small drops shadows completely?
			// setting to 0.1 or 0.5 shows shadows; setting too big gives weirdness
			// see https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/ligth-and-shadows
			// on shadow bias
			e := 1E10
			for _, o := range q.Scene.Objects {
				if _, ok, dis := o.Intersect(shadow); ok && dis > e && dis < segmentLength {
					continue Lights
				}
			}
			facingRatio := objectHit.SurfaceNormal(*intersection).Dot(segment)
			if facingRatio <= 0 {
				continue Lights
			}
			color = Color{color.Add(l.Color.Times(STANDARD_ALBEDO * l.Intensity * facingRatio))}
		}

		ans <- Answer{q.X, q.Y, color}
	}
}
