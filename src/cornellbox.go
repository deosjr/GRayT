package main

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
)

func CornellBox() *m.Scene {
	camera := m.NewPerspectiveCamera(width, height, 0.5*math.Pi)
	scene := m.NewScene(camera)

	// use pointlight for whitted style ray tracer
	pointLight := m.NewPointLight(m.Vector{250, 500, 100}, m.NewColor(255, 255, 255), 50000000)
	scene.AddLights(pointLight)

	// use area light for path tracer
	intensity := 5.0
	lightMat := &m.RadiantMaterial{Color: m.NewColor(255, 255, 255).Times(intensity)}

	/*
		light := m.NewQuadrilateral(
			m.Vector{213, 548.8, 332},
			m.Vector{343, 548.8, 332},
			m.Vector{343, 548.8, 227},
			m.Vector{213, 548.8, 227},
			lightMat)
		scene.Add(light.Tesselate())
	*/

	white := &m.DiffuseMaterial{Color: m.NewColor(186, 186, 186)}
	green := &m.DiffuseMaterial{Color: m.NewColor(31, 115, 38)}
	red := &m.DiffuseMaterial{Color: m.NewColor(166, 13, 13)}

	floor := m.NewQuadrilateral(
		m.Vector{0.0, 0.0, 559.2},
		m.Vector{0.0, 0.0, 0.0},
		m.Vector{552.8, 0.0, 0.0},
		m.Vector{549.6, 0.0, 559.2},
		white)
	scene.Add(floor.Tesselate())

	ceiling := m.NewQuadrilateral(
		m.Vector{556.0, 548.8, 0.0},
		m.Vector{0.0, 548.8, 0.0},
		m.Vector{0.0, 548.8, 559.2},
		m.Vector{556.0, 548.8, 559.2},
		lightMat) //white
	scene.Add(ceiling.Tesselate())

	backwall := m.NewQuadrilateral(
		m.Vector{0.0, 548.8, 559.2},
		m.Vector{0.0, 0.0, 559.2},
		m.Vector{549.6, 0.0, 559.2},
		m.Vector{556.0, 548.8, 559.2},
		white)
	scene.Add(backwall.Tesselate())

	rightwall := m.NewQuadrilateral(
		m.Vector{0.0, 548.8, 0.0},
		m.Vector{0.0, 0.0, 0.0},
		m.Vector{0.0, 0.0, 559.2},
		m.Vector{0.0, 548.8, 559.2},
		green)
	scene.Add(rightwall.Tesselate())

	leftwall := m.NewQuadrilateral(
		m.Vector{556.0, 548.8, 559.2},
		m.Vector{549.6, 0.0, 559.2},
		m.Vector{552.8, 0.0, 0.0},
		m.Vector{556.0, 548.8, 0.0},
		red)
	scene.Add(leftwall.Tesselate())

	//shortblock
	shortblock1 := m.NewQuadrilateral(
		m.Vector{240, 165, 272},
		m.Vector{82, 165, 225},
		m.Vector{130, 165, 65},
		m.Vector{290, 165, 114},
		white)
	scene.Add(shortblock1.Tesselate())
	shortblock2 := m.NewQuadrilateral(
		m.Vector{240, 165, 272},
		m.Vector{290, 165, 114},
		m.Vector{290, 0, 114},
		m.Vector{240, 0, 272},
		white)
	scene.Add(shortblock2.Tesselate())
	shortblock3 := m.NewQuadrilateral(
		m.Vector{290, 165, 114},
		m.Vector{130, 165, 65},
		m.Vector{130, 0, 65},
		m.Vector{290, 0, 114},
		white)
	scene.Add(shortblock3.Tesselate())
	shortblock4 := m.NewQuadrilateral(
		m.Vector{130, 165, 65},
		m.Vector{82, 165, 225},
		m.Vector{82, 0, 225},
		m.Vector{130, 0, 65},
		white)
	scene.Add(shortblock4.Tesselate())
	shortblock5 := m.NewQuadrilateral(
		m.Vector{82, 165, 225},
		m.Vector{240, 165, 272},
		m.Vector{240, 0, 272},
		m.Vector{82, 0, 225},
		white)
	scene.Add(shortblock5.Tesselate())

	//tallblock
	tallblock1 := m.NewQuadrilateral(
		m.Vector{314, 330, 456},
		m.Vector{265, 330, 296},
		m.Vector{423, 330, 247},
		m.Vector{472, 330, 406},
		white)
	scene.Add(tallblock1.Tesselate())
	tallblock2 := m.NewQuadrilateral(
		m.Vector{472, 330, 406},
		m.Vector{423, 330, 247},
		m.Vector{423, 0, 247},
		m.Vector{472, 0, 406},
		white)
	scene.Add(tallblock2.Tesselate())
	tallblock3 := m.NewQuadrilateral(
		m.Vector{314, 330, 456},
		m.Vector{472, 330, 406},
		m.Vector{472, 0, 406},
		m.Vector{314, 0, 456},
		white)
	scene.Add(tallblock3.Tesselate())
	tallblock4 := m.NewQuadrilateral(
		m.Vector{265, 330, 296},
		m.Vector{314, 330, 456},
		m.Vector{314, 0, 456},
		m.Vector{265, 0, 296},
		white)
	scene.Add(tallblock4.Tesselate())
	tallblock5 := m.NewQuadrilateral(
		m.Vector{423, 330, 247},
		m.Vector{265, 330, 296},
		m.Vector{265, 0, 296},
		m.Vector{423, 0, 247},
		white)
	scene.Add(tallblock5.Tesselate())

	scene.Precompute()

	from, to := m.Vector{278, 273, -800}, m.Vector{278, 273, -799}
	camera.LookAt(from, to, ey)

	return scene
}
