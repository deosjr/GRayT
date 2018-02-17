package main

import (
	"math"

	m "model"
	"projects"
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
	l1 := m.NewPointLight(m.Vector{-2, 2, 0}, m.NewColor(255, 255, 255), 300)
	l2 := m.NewPointLight(m.Vector{-5, 5, 3}, m.NewColor(255, 255, 255), 600)
	scene.AddLights(l1, l2)

	scene.Add(m.NewSphere(m.Vector{3, 1, 5}, 0.5, m.NewColor(255, 100, 0)))

	// triangles
	r := m.Quadrilateral{
		m.Vector{0, -1, 6},
		m.Vector{4, -1, 3},
		m.Vector{0, -1, 0},
		m.Vector{-4, -1, 3},
		m.NewColor(255, 0, 0)}

	grid := projects.ToPointGrid(r, 0.1)
	grid = projects.PerlinHeightMap(grid)
	triangles := m.NewTriangleMesh(grid)
	scene.Add(triangles...)

	scene.Precompute()

	aw := render.NewAVI("out.avi", width, height)
	from, to := m.Vector{0, 0, 0}, m.Vector{0, 0, 10}
	for i := 0; i < 15; i++ {
		camera.LookAt(from, to, ey)
		film := render.Render(scene, numWorkers)
		render.AddToAVI(aw, film)
		from = from.Add(m.Vector{0.1, 0, 0.1})
	}
	render.SaveAVI(aw)
}
