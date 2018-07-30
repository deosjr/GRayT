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

	fmt.Println("Creating scene...")
	camera := m.NewPerspectiveCamera(width, height, 0.5*math.Pi)

	scene := m.NewScene(camera)
	l1 := m.NewPointLight(m.Vector{-2, 2, 0}, m.NewColor(255, 255, 255), 500)
	l2 := m.NewPointLight(m.Vector{-0.1, 3, 0.5}, m.NewColor(255, 255, 255), 600)
	scene.AddLights(l1, l2)

	diffMat := &m.DiffuseMaterial{m.NewColor(255, 0, 0)}
	reflMat := &m.ReflectiveMaterial{scene}

	scene.Add(m.NewSphere(m.Vector{-1, 2, 2}, 0.5, reflMat))
	scene.Add(m.NewSphere(m.Vector{1, 2, 2}, 0.5, reflMat))
	scene.Add(m.NewSphere(m.Vector{0, 1, 1}, 0.5, diffMat))

	fmt.Println("Building BVH...")
	scene.Precompute()

	fmt.Println("Rendering...")

	// aw := render.NewAVI("out.avi", width, height)
	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")

	// for i := 0; i < 30; i++ {
	// 	camera.LookAt(from, to, ey)
	// 	film := render.Render(scene, numWorkers)
	// 	render.AddToAVI(aw, film)
	// 	from = from.Add(m.Vector{0, 0, -0.05})
	// }
	// render.SaveAVI(aw)
}
