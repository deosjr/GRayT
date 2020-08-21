package main

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
)

func CornellBox() *m.Scene {
	// original cornell box data gives focal length of 35mm
	// vertical angle of view for 35mm f is roughly 37.8
	fov := float32((37.8 / 360.0) * 2 * math.Pi)
	camera := m.NewPerspectiveCamera(width, height, fov)
	scene := m.NewScene(camera)

	// use pointlight for whitted style ray tracer
	pointLight := m.NewPointLight(m.Vector{250, 500, 100}, m.NewColor(255, 255, 255), 50000000)
	scene.AddLights(pointLight)

	// use area light for path tracer
	var intensity float32 = 100.0
	lightMat := m.NewRadiantMaterial(m.ConstantTexture{Color: m.NewColor(255, 255, 255).Times(intensity)})

	triangles := []m.Triangle{}

	// light Y just lower than ceiling, otherwise we collide
	// todo: data gives a hole in the ceiling for the light, im just lazy..
	light := m.NewQuadrilateral(
		m.Vector{343, 548.7, 227},
		m.Vector{343, 548.7, 332},
		m.Vector{213, 548.7, 332},
		m.Vector{213, 548.7, 227},
		lightMat)
	t1, t2 := light.Tesselate()
	triangles = append(triangles, t1, t2)
	scene.Emitters = []m.Triangle{t1, t2}

	white := m.NewDiffuseMaterial(m.ConstantTexture{Color: m.NewColor(186, 186, 186)})
	green := m.NewDiffuseMaterial(m.ConstantTexture{Color: m.NewColor(31, 115, 38)})
	red := m.NewDiffuseMaterial(m.ConstantTexture{Color: m.NewColor(166, 13, 13)})

	floor := m.NewQuadrilateral(
		m.Vector{552.8, 0.0, 0.0},
		m.Vector{0.0, 0.0, 0.0},
		m.Vector{0.0, 0.0, 559.2},
		m.Vector{549.6, 0.0, 559.2},
		white)
	t1, t2 = floor.Tesselate()
	triangles = append(triangles, t1, t2)

	ceiling := m.NewQuadrilateral(
		m.Vector{556.0, 548.8, 0.0},
		m.Vector{556.0, 548.8, 559.2},
		m.Vector{0.0, 548.8, 559.2},
		m.Vector{0.0, 548.8, 0.0},
		white)
	t1, t2 = ceiling.Tesselate()
	triangles = append(triangles, t1, t2)

	backwall := m.NewQuadrilateral(
		m.Vector{549.6, 0.0, 559.2},
		m.Vector{0.0, 0.0, 559.2},
		m.Vector{0.0, 548.8, 559.2},
		m.Vector{556.0, 548.8, 559.2},
		white)
	t1, t2 = backwall.Tesselate()
	triangles = append(triangles, t1, t2)

	rightwall := m.NewQuadrilateral(
		m.Vector{0.0, 0.0, 559.2},
		m.Vector{0.0, 0.0, 0.0},
		m.Vector{0.0, 548.8, 0.0},
		m.Vector{0.0, 548.8, 559.2},
		green)
	t1, t2 = rightwall.Tesselate()
	triangles = append(triangles, t1, t2)

	leftwall := m.NewQuadrilateral(
		m.Vector{552.8, 0.0, 0.0},
		m.Vector{549.6, 0.0, 559.2},
		m.Vector{556.0, 548.8, 559.2},
		m.Vector{556.0, 548.8, 0.0},
		red)
	t1, t2 = leftwall.Tesselate()
	triangles = append(triangles, t1, t2)

	//shortblock
	shortblock1 := m.NewQuadrilateral(
		m.Vector{130, 165, 65},
		m.Vector{82, 165, 225},
		m.Vector{240, 165, 272},
		m.Vector{290, 165, 114},
		white)
	shortblock2 := m.NewQuadrilateral(
		m.Vector{290, 0, 114},
		m.Vector{290, 165, 114},
		m.Vector{240, 165, 272},
		m.Vector{240, 0, 272},
		white)
	shortblock3 := m.NewQuadrilateral(
		m.Vector{130, 0, 65},
		m.Vector{130, 165, 65},
		m.Vector{290, 165, 114},
		m.Vector{290, 0, 114},
		white)
	shortblock4 := m.NewQuadrilateral(
		m.Vector{82, 0, 225},
		m.Vector{82, 165, 225},
		m.Vector{130, 165, 65},
		m.Vector{130, 0, 65},
		white)
	shortblock5 := m.NewQuadrilateral(
		m.Vector{240, 0, 272},
		m.Vector{240, 165, 272},
		m.Vector{82, 165, 225},
		m.Vector{82, 0, 225},
		white)
	t1, t2 = shortblock1.Tesselate()
	t3, t4 := shortblock2.Tesselate()
	t5, t6 := shortblock3.Tesselate()
	t7, t8 := shortblock4.Tesselate()
	t9, t10 := shortblock5.Tesselate()
	triangles = append(triangles, t1, t2, t3, t4, t5, t6, t7, t8, t9, t10)

	//tallblock
	tallblock1 := m.NewQuadrilateral(
		m.Vector{423, 330, 247},
		m.Vector{265, 330, 296},
		m.Vector{314, 330, 456},
		m.Vector{472, 330, 406},
		white)
	tallblock2 := m.NewQuadrilateral(
		m.Vector{423, 0, 247},
		m.Vector{423, 330, 247},
		m.Vector{472, 330, 406},
		m.Vector{472, 0, 406},
		white)
	tallblock3 := m.NewQuadrilateral(
		m.Vector{472, 0, 406},
		m.Vector{472, 330, 406},
		m.Vector{314, 330, 456},
		m.Vector{314, 0, 456},
		white)
	tallblock4 := m.NewQuadrilateral(
		m.Vector{314, 0, 456},
		m.Vector{314, 330, 456},
		m.Vector{265, 330, 296},
		m.Vector{265, 0, 296},
		white)
	tallblock5 := m.NewQuadrilateral(
		m.Vector{265, 0, 296},
		m.Vector{265, 330, 296},
		m.Vector{423, 330, 247},
		m.Vector{423, 0, 247},
		white)
	t1, t2 = tallblock1.Tesselate()
	t3, t4 = tallblock2.Tesselate()
	t5, t6 = tallblock3.Tesselate()
	t7, t8 = tallblock4.Tesselate()
	t9, t10 = tallblock5.Tesselate()
	triangles = append(triangles, t1, t2, t3, t4, t5, t6, t7, t8, t9, t10)

	scene.Add(m.NewTriangleComplexObject(triangles))
	scene.Precompute()

	from, to := m.Vector{278, 273, -800}, m.Vector{278, 273, -799}
	camera.LookAt(from, to, ey)

	return scene
}
