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

func main() {

	camera := model.NewCamera(WIDTH, HEIGHT)

	scene = model.NewScene(camera)
	scene.Add(model.Sphere{model.Vector{0, 0, 5}, 1.0})
	scene.Add(model.Sphere{model.Vector{5, 0, 5}, 1.0})
	scene.Add(model.NewPlane(ex, ey, model.Vector{0, 0, 6}))
	scene.Add(model.NewPlane(ex, ez, model.Vector{0, -2, 0}))

	ch := make(chan model.Question, NUMWORKERS)
	ans := make(chan model.Answer, NUMWORKERS)

	for i := 0; i < NUMWORKERS; i++ {
		go model.Worker(ch, ans)
	}

	go func() {
		for y := 0; y < int(HEIGHT); y++ {
			for x := 0; x < int(WIDTH); x++ {
				ch <- model.Question{scene, x, y}
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
		camera.Image.Set(a.X, a.Y, a.Color)
		numPixels--
	}

	camera.Image.Save()

}
