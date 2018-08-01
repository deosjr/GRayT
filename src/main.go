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

	l1 := m.NewDistantLight(m.Vector{1, -1, 1}, m.NewColor(255, 255, 255), 20)
	// l2 := m.NewPointLight(m.Vector{1, 2, 3}, m.NewColor(255, 255, 255), 200)
	scene.AddLights(l1)

	diffMat := &m.DiffuseMaterial{m.NewColor(50, 10, 100)}
	reflMat := &m.ReflectiveMaterial{scene}

	rectangle := m.NewQuadrilateral(
		m.Vector{-1, 0, -1},
		m.Vector{1, 0, -1},
		m.Vector{1, 0, 1},
		m.Vector{-1, 0, 1},
		reflMat)
	plane := rectangle.Tesselate()
	translation := m.Translate(m.Vector{0, 0, 3})//.Mul(m.ScaleUniform(2))
	// TODO: current issues very clear when adding a 
	// scaling transformation to the 'plane' object
	// I still suspect normal transformations to be off too
	scene.Add(m.NewSharedObject(plane, translation))

	// NOTE: dividing the plane up into more small rectangles works just fine
	// somehow the size of the object determines light intensity... ?
	// scene.Add(m.NewSharedObject(plane, m.Translate(m.Vector{1, 0, 3})))
	// scene.Add(m.NewSharedObject(plane, m.Translate(m.Vector{-1, 0, 3})))
	// scene.Add(m.NewSharedObject(plane, m.Translate(m.Vector{0, 0, 2})))
	// scene.Add(m.NewSharedObject(plane, m.Translate(m.Vector{-1, 0, 2})))
	// scene.Add(m.NewSharedObject(plane, m.Translate(m.Vector{1, 0, 2})))

	box := m.NewAABB(m.Vector{-0.1, -0.1, -0.1}, m.Vector{0.1, 0.1, 0.1})
	c := m.NewCuboid(box, diffMat)
	object := c.Tesselate()

	rotation := m.RotateY(math.Pi / 4) 
	translation = m.Translate(m.Vector{0, 0.5, 3})
	shared := m.NewSharedObject(object, translation.Mul(rotation))
	scene.Add(shared)

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
