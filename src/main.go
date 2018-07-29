package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

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

	// mat := &m.DiffuseMaterial{m.NewColor(255, 0, 0)}
	mat := &m.PosFuncMat{
		func(p m.Vector) m.Color {
			return m.NewColor(uint8((p.X+2)*60), uint8((p.Y)*60), uint8((-p.Z+1)*60))
		},
	}

	box := m.NewAABB(m.Vector{-0.1, -0.1, -0.1}, m.Vector{0.1, 0.1, 0.1})
	c := m.NewCuboid(box, mat)

	object := c.Tesselate()
	//objHeight := -object.Bound().Pmin.Y
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		x := -2 + rand.Float64()*4
		y := rand.Float64() * 4
		z := -1 - rand.Float64()*4
		shared := m.NewSharedObject(object, m.Translate(m.Vector{x, y, z}))
		scene.Add(shared)
	}

	fmt.Println("Building BVH...")
	scene.Precompute()

	fmt.Println("Rendering...")

	// aw := render.NewAVI("out.avi", width, height)
	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, -10}
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
