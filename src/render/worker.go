package render

import (
	"math"
	"model"
)

var BACKGROUND_COLOR = model.NewColor(100, 100, 100)
var MAX_RAY_DISTANCE = math.MaxFloat64

type question struct {
	x, y int
}

type answer struct {
	x, y  int
	color model.Color
}

func worker(scene *Scene, ch chan question, ans chan answer) {
	for q := range ch {
		ray := scene.Camera.PixelRay(q.x, q.y)
		hit := scene.AccelerationStructure.ClosestIntersection(ray, MAX_RAY_DISTANCE)
		if hit == nil {
			ans <- answer{q.x, q.y, BACKGROUND_COLOR}
			continue
		}

		color := model.NewColor(0, 0, 0)
		for _, l := range scene.Lights {
			c, ok := model.LightContribution(ray, hit, l, scene.AccelerationStructure)
			if !ok {
				continue
			}
			color = color.Add(c)
		}
		ans <- answer{q.x, q.y, color}
	}
}
