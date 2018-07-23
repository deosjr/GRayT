package main

import (
	"fmt"
	"math"
	m "model"
	"render"
)

var (
	width      uint = 1600
	height     uint = 1200
	numWorkers      = 10

	ex = m.Vector{1, 0, 0}
	ey = m.Vector{0, 1, 0}
	ez = m.Vector{0, 0, 1}
)

func main() {
	camera := m.NewPerspectiveCamera(width, height, 0.5*math.Pi)
	//camera := m.NewOrthographicCamera(width, height)

	scene := render.NewScene(camera)
	l1 := m.NewPointLight(m.Vector{-2, 2, 0}, m.NewColor(255, 255, 255), 800)
	l2 := m.NewPointLight(m.Vector{-0.1, 1, 0.1}, m.NewColor(255, 255, 255), 400)
	scene.AddLights(l1, l2)

	scene.Add(m.NewSphere(m.Vector{3, 1, 5}, 0.5, m.NewColor(255, 100, 0)))

	scene.Add(m.NewPlane(m.Vector{0, 0, 0}, ez, ex, m.NewColor(40, 200, 40)))
	scene.Add(m.NewPlane(m.Vector{-1, 0, -1}, ex, ey, m.NewColor(0, 40, 100)))

	object, err := render.LoadObj("bunny.obj", m.NewColor(160, 80, 0))
	if err != nil {
		fmt.Printf("Error reading file: %s \n", err.Error())
	}
	scene.Add(object)

	fmt.Println("Building BVH...")
	scene.Precompute()

	fmt.Println("Rendering...")

	// aw := render.NewAVI("out.avi", width, height)
	from, to := m.Vector{-0.2, 0.2, 0.2}, m.Vector{-0.08813016500000001, 0.14181918999999998, 0.011103720000000001}
	//from, to := m.Vector{-2, 2, 2}, m.Vector{-1, 2, 2}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")

	// for i := 0; i < 15; i++ {
	// 	camera.LookAt(from, to, ey)
	// 	film := render.Render(scene, numWorkers)
	// 	render.AddToAVI(aw, film)
	// 	from = from.Add(m.Vector{0.1, 0, 0.1})
	// }
	// render.SaveAVI(aw)
}
