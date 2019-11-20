package render

import (
	"math"
	"testing"

	"github.com/deosjr/GRayT/src/model"
)

func sampleScene(b *testing.B) *model.Scene {
	camera := model.NewPerspectiveCamera(1600, 1200, 0.5*math.Pi)
	scene := model.NewScene(camera)
	l1 := model.NewPointLight(model.Vector{-2, 2, 0}, model.NewColor(255, 255, 255), 300)
	l2 := model.NewPointLight(model.Vector{-0.1, 1, 0.1}, model.NewColor(255, 255, 255), 400)
	scene.AddLights(l1, l2)
	scene.Add(model.NewSphere(model.Vector{3, 1, 5}, 0.5, &model.DiffuseMaterial{Color: model.NewColor(255, 100, 0)}))

	scene.Precompute()

	from, to := model.Vector{0, 0, 0}, model.Vector{0, 0, 10}
	camera.LookAt(from, to, model.Vector{0, 1, 0})
	return scene
}

func benchmarkScene(scene *model.Scene, numWorkers int, b *testing.B) {
	params := Params{
		Scene:        scene,
		NumWorkers:   numWorkers,
		NumSamples:   1,
		AntiAliasing: false,
		TracerType:   model.WhittedStyle,
	}
	for i := 0; i < b.N; i++ {
		Render(params)
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
