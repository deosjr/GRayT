package render

import (
	"runtime"
	"testing"

	"model"
)

var (
	ex    = model.Vector{1, 0, 0}
	ey    = model.Vector{0, 1, 0}
	ez    = model.Vector{0, 0, 1}
	white = model.NewColor(255, 255, 255)
)

func sampleScene() *Scene {
	camera := model.NewCamera(160, 120)
	scene := NewScene(camera)
	l1 := model.NewPointLight(model.Vector{0, 4, 0}, model.NewColor(0, 0, 255), 1500)
	l2 := model.NewPointLight(model.Vector{-5, 5, 0}, model.NewColor(255, 0, 0), 1000)
	scene.AddLights(l1, l2)
	scene.Add(model.Sphere{model.Vector{0, 0, -5}, 1.0, white})
	scene.Add(model.Sphere{model.Vector{5, 0, -5}, 1.0, white})
	scene.Add(model.NewPlane(model.Vector{0, 0, -10}, ex, ey, white))
	scene.Add(model.NewPlane(model.Vector{0, -2, 0}, ez, ex, white))
	return scene
}

func benchmarkScene(scene *Scene, numWorkers int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		Render(scene, numWorkers)
	}
}

func BenchmarkNumWorkers1(b *testing.B) {
	benchmarkScene(sampleScene(), 1, b)
}

func BenchmarkNumWorkers10(b *testing.B) {
	benchmarkScene(sampleScene(), 10, b)
}

func BenchmarkNumWorkersNumCPU(b *testing.B) {
	benchmarkScene(sampleScene(), runtime.NumCPU(), b)
}
