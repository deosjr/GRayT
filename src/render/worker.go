package render

import "model"

var BACKGROUND_COLOR = model.NewColor(10, 10, 10)

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
		hit := model.ClosestIntersection(ray, scene.Objects)
		if hit == nil {
			ans <- answer{q.x, q.y, BACKGROUND_COLOR}
			continue
		}

		color := model.NewColor(0, 0, 0)
		for _, l := range scene.Lights {
			c, ok := model.LightContribution(ray, hit, l, scene.Objects)
			if !ok {
				continue
			}
			color = color.Add(c)
		}
		ans <- answer{q.x, q.y, color}
	}
}
