package main

import (
	"model"
)

var (
	WIDTH      uint = 1600
	HEIGHT     uint = 1200
	NUMWORKERS      = 10

	ex = model.Vector{1, 0, 0}
	ey = model.Vector{0, 1, 0}
	ez = model.Vector{0, 0, 1}
)

func main() {

	camera := model.NewCamera(WIDTH, HEIGHT)

	scene := model.NewScene(camera)
	scene.AddLight(model.Vector{0, 4, 0}, model.NewColor(0, 0, 255), 1500)
	scene.AddLight(model.Vector{-5, 5, 0}, model.NewColor(255, 0, 0), 1000)
	scene.Add(model.Sphere{model.Vector{0, -1, -5}, 1.0})
	scene.Add(model.Sphere{model.Vector{3, 0, -5}, 1.0})
	scene.Add(model.Sphere{model.Vector{-3, 1, -5}, 1.0})
	scene.Add(model.NewPlane(ex, ey, model.Vector{0, 0, -10}))
	scene.Add(model.NewPlane(ez, ex, model.Vector{0, -2, 0}))

	img := model.Render(scene, NUMWORKERS)
	img.Save("out.png")

}
