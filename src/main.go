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
	scene.AddLight(model.Vector{0, 4, 0}, model.NewColor(255, 255, 255), 500)
	scene.AddLight(model.Vector{-5, 5, 0}, model.NewColor(255, 255, 255), 500)
	scene.Add(model.Sphere{model.Vector{0, -1, -5}, 1.0, model.NewColor(255, 0, 0)})
	scene.Add(model.Sphere{model.Vector{3, 0, -5}, 1.0, model.NewColor(100, 100, 100)})
	scene.Add(model.Sphere{model.Vector{-3, 1, -5}, 1.0, model.NewColor(0, 0, 255)})
	scene.Add(model.NewPlane(model.Vector{0, 0, -10}, ex, ey, model.NewColor(50, 200, 240)))
	scene.Add(model.NewPlane(model.Vector{0, -2, 0}, ez, ex, model.NewColor(45, 200, 45)))
	scene.Add(model.Triangle{model.Vector{3, 0, -4}, model.Vector{4, 1, -4}, model.Vector{3, 1, -4}, model.NewColor(0, 255, 0)})

	img := model.Render(scene, NUMWORKERS)
	img.Save("out.png")

}
