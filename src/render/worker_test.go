package render

import (
	"math"
	"testing"

	"model"
)

func sampleScene(b *testing.B) *Scene {
	camera := model.NewPerspectiveCamera(1600, 1200, 0.5*math.Pi)
	scene := NewScene(camera)
	l1 := model.NewPointLight(model.Vector{-2, 2, 0}, model.NewColor(255, 255, 255), 300)
	l2 := model.NewPointLight(model.Vector{-0.1, 1, 0.1}, model.NewColor(255, 255, 255), 400)
	scene.AddLights(l1, l2)
	scene.Add(model.NewSphere(model.Vector{3, 1, 5}, 0.5, model.NewColor(255, 100, 0)))

	triangles, err := LoadObj("../bunny.obj", model.NewColor(255, 0, 0))
	if err != nil {
		b.Fatalf("Error in benchmark: %s", err.Error())
	}
	scene.Add(triangles...)

	scene.Precompute()

	from, to := model.Vector{0, 0, 0}, model.Vector{0, 0, 10}
	camera.LookAt(from, to, model.Vector{0, 1, 0})
	return scene
}

func benchmarkScene(scene *Scene, numWorkers int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		Render(scene, numWorkers)
	}
}

// Numbers slightly exaggerated by order of tests

func BenchmarkNumWorkers1(b *testing.B) {
	benchmarkScene(sampleScene(b), 1, b)
}

func BenchmarkNumWorkers10(b *testing.B) {
	benchmarkScene(sampleScene(b), 10, b)
}

func BenchmarkNumWorkers50(b *testing.B) {
	benchmarkScene(sampleScene(b), 50, b)
}

func BenchmarkNumWorkers100(b *testing.B) {
	benchmarkScene(sampleScene(b), 100, b)
}
