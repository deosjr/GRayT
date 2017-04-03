package main

import (
	"model"
)

var (
	NUMWORKERS      = 10
	WIDTH      uint = 1600
	HEIGHT     uint = 1200

	scene *model.Scene

	ex = model.Vector{1, 0, 0}
	ey = model.Vector{0, 1, 0}
	ez = model.Vector{0, 0, 1}
)

type question struct {
	x, y int
}

type answer struct {
	x, y  int
	color uint8
}

func main() {

	camera := model.NewCamera(model.Vector{0, 0, 0}, ez, 1, WIDTH, HEIGHT)

	scene = model.NewScene(camera)
	scene.AddLight(0, 5, 0)
	scene.Add(model.Sphere{model.Vector{0, 0, 5}, 1.0})
	//scene.Add(model.NewPlane(ex, ey, model.Vector{0, 0, 6}))

	ch := make(chan question, NUMWORKERS)
	ans := make(chan answer, NUMWORKERS)

	for i := 0; i < NUMWORKERS; i++ {
		go worker(ch, ans)
	}

	go func() {
		for y := 0; y < int(HEIGHT); y++ {
			for x := 0; x < int(WIDTH); x++ {
				ch <- question{x, y}
			}
		}
		close(ch)
	}()

	numPixels := HEIGHT * WIDTH
	for {
		if numPixels == 0 {
			break
		}
		a := <-ans
		camera.Image.Set(a.x, a.y, a.color)
		numPixels--
	}

	camera.Image.Save()

}

func worker(ch chan question, ans chan answer) {
	for q := range ch {

		origin := scene.Camera.PixelVector(q.x, q.y)
		direction := model.VectorFromTo(scene.Camera.Origin, origin)
		ray := model.NewRay(origin, direction)

		var color uint8 = 50
		for _, o := range scene.Objects {
			if intersection, ok, _ := o.Intersect(ray); ok {

			Lights:
				for _, l := range scene.Lights {
					segment := model.VectorFromTo(intersection, l.Origin)
					shadow := model.NewRay(intersection, segment)
					segmentLength := segment.Length()
					for _, oo := range scene.Objects {
						if oo == o {
							continue
						}
						if _, ok, distance := oo.Intersect(shadow); ok && distance < segmentLength {
							continue Lights
						}
					}
					// Nothing blocking intersection point from getting light!
					color = uint8(-segmentLength * 100)
				}

			}
		}
		ans <- answer{q.x, q.y, color}
	}
}
