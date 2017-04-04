package model

import (
	"math"
)

type Question struct {
	Scene *Scene
	X, Y  int
}

type Answer struct {
	X, Y  int
	Color uint8
}

func Worker(ch chan Question, ans chan Answer) {
	for q := range ch {
		ray := q.Scene.Camera.PixelRay(q.X, q.Y)

		var color uint8 = 50
		var intersection *Vector
		var minDistance = math.MaxFloat64
		for _, o := range q.Scene.Objects {
			if i, ok, distance := o.Intersect(ray); ok && distance >= 0 && distance < minDistance {
				minDistance = distance
				intersection = &i
			}
		}

		if intersection == nil {
			ans <- Answer{q.X, q.Y, color}
			continue
		}

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
			// Nothing blocking intersection point from getting light!
			color = uint8(-segmentLength * 100)
		}

		ans <- Answer{q.X, q.Y, color}
	}
}
