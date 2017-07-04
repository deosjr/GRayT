package model

import (
	"runtime"
	"testing"
)

var (
	ex    = Vector{1, 0, 0}
	ey    = Vector{0, 1, 0}
	ez    = Vector{0, 0, 1}
	white = NewColor(255, 255, 255)
)

func sampleScene() *Scene {
	camera := NewCamera(160, 120)
	scene := NewScene(camera)
	scene.AddLight(Vector{0, 4, 0}, NewColor(0, 0, 255), 1500)
	scene.AddLight(Vector{-5, 5, 0}, NewColor(255, 0, 0), 1000)
	scene.Add(Sphere{Vector{0, 0, -5}, 1.0, white})
	scene.Add(Sphere{Vector{5, 0, -5}, 1.0, white})
	scene.Add(NewPlane(Vector{0, 0, -10}, ex, ey, white))
	scene.Add(NewPlane(Vector{0, -2, 0}, ez, ex, white))
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
