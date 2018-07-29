package main

import (
	"fmt"
	"math"
	"math/rand"
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
	//camera := m.NewOrthographicCamera(width, height)

	scene := render.NewScene(camera)
	l1 := m.NewPointLight(m.Vector{-2, 2, 0}, m.NewColor(255, 255, 255), 500)
	l2 := m.NewPointLight(m.Vector{-0.1, 1, -1}, m.NewColor(255, 255, 255), 400)
	scene.AddLights(l1, l2)

	c := m.Cuboid{
		m.Vector{-0.1, 0.1, -0.1},
		m.Vector{0.1, 0.1, -0.1},
		m.Vector{0.1, 0.1, 0.1},
		m.Vector{-0.1, 0.1, 0.1},
		m.Vector{-0.1, -0.1, -0.1},
		m.Vector{0.1, -0.1, -0.1},
		m.Vector{0.1, -0.1, 0.1},
		m.Vector{-0.1, -0.1, 0.1},
		m.NewColor(255, 0, 0),
	}

	object := m.NewComplexObject(c.Tesselate())
	//objHeight := -object.Bound().Pmin.Y
	for i := 0; i < 100; i++ {
		x := rand.Float64() * 2
		y := rand.Float64() * 2
		z := -rand.Float64() - 2
		shared := m.NewSharedObject(object, m.Translate(m.Vector{x, y, z}))
		scene.Add(shared)
	}

	fmt.Println("Building BVH...")
	scene.Precompute()

	fmt.Println("Rendering...")

	//aw := render.NewAVI("out.avi", width, height)
	from, to := m.Vector{0.2, 0.2, 0}, m.Vector{-1, 0, -10}
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
