package main

import (
	m "model"
	"projects"
	"render"
)

var (
	WIDTH      uint = 1600
	HEIGHT     uint = 1200
	NUMWORKERS      = 10

	ex = m.Vector{1, 0, 0}
	ey = m.Vector{0, 1, 0}
	ez = m.Vector{0, 0, 1}
)

func main() {

	camera := m.NewCamera(WIDTH, HEIGHT)

	scene := render.NewScene(camera)
	l1 := m.NewPointLight(m.Vector{2, 2, 0}, m.NewColor(255, 255, 255), 300)
	l2 := m.NewPointLight(m.Vector{-5, 5, -3}, m.NewColor(255, 255, 255), 300)
	scene.AddLights(l1, l2)
	// background
	scene.Add(m.NewPlane(m.Vector{0, 0, -10}, ex, ey, m.NewColor(50, 200, 240)))
	// floor
	scene.Add(m.NewPlane(m.Vector{0, -2, 0}, ez, ex, m.NewColor(45, 200, 45)))

	scene.Add(m.Sphere{m.Vector{-2, 1, -4}, 1.0, m.NewColor(0, 0, 255)})

	c := m.Cuboid{
		m.Vector{1.5, 1, -4},
		m.Vector{2, 1, -4},
		m.Vector{2, 1, -3.5},
		m.Vector{1.5, 1, -3.5},
		m.Vector{1.5, 0.5, -4},
		m.Vector{2, 0.5, -4},
		m.Vector{2, 0.5, -3.5},
		m.Vector{1.5, 0.5, -3.5},
		m.NewColor(255, 0, 0),
	}
	scene.Add(c.Tesselate()...)

	// triangles
	r := m.Quadrilateral{
		m.Vector{-1, -1, -4},
		m.Vector{1, -1, -4},
		m.Vector{1, -1, -2},
		m.Vector{-1, -1, -2},
		m.NewColor(255, 0, 0)}

	scene.Add(projects.GridToTriangles(projects.ToPointGrid(r, 0.2))...)

	img := render.Render(scene, NUMWORKERS)
	img.Save("out.png")
}
