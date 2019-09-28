package simd

import "math"
import "testing"

// testing three different implementations:
// - vectors as [4]float32 with last value 0, using sse instructions with go assembly
// - vectors as [4]float32 type aliased with last value 0, using pure go
// - vectors as structs, using 3 float32 fields in pure go

type testVector struct {
	X, Y, Z float32
}

type testVArray [4]float32

// benchmarking super simple piecewise addition
// verdict: vector is fastest
func BenchmarkAddSimd(b *testing.B) {
	u := testVArray([4]float32{1, 2, 3, 0})
	v := testVArray([4]float32{1, 2, 3, 0})
	for i := 0; i < b.N; i++ {
		Add(u, v)
	}
}

func (u testVArray) Add(v testVArray) testVArray {
	return [4]float32{u[0] + v[0], u[1] + v[1], u[2] + v[2], 0}
}

func BenchmarkAddArray(b *testing.B) {
	u := testVArray([4]float32{1, 2, 3, 0})
	v := testVArray([4]float32{1, 2, 3, 0})
	for i := 0; i < b.N; i++ {
		u.Add(v)
	}
}

func (u testVector) Add(v testVector) testVector {
	return testVector{
		X: u.X + v.X,
		Y: u.Y + v.Y,
		Z: u.Z + v.Z,
	}
}

func BenchmarkAddStruct(b *testing.B) {
	u := testVector{1, 2, 3}
	v := testVector{1, 2, 3}
	for i := 0; i < b.N; i++ {
		u.Add(v)
	}
}

// benchmarking dot product
// verdict: simd way slow, array slightly faster than vector
func BenchmarkDotSimd(b *testing.B) {
	u := testVArray([4]float32{1, 2, 3, 0})
	v := testVArray([4]float32{1, 2, 3, 0})
	for i := 0; i < b.N; i++ {
		Dot(u, v)
	}
}

func (u testVArray) Dot(v testVArray) float32 {
	return u[0]*v[0] + u[1]*v[1] + u[2]*v[2]
}

func BenchmarkDotArray(b *testing.B) {
	u := testVArray([4]float32{1, 2, 3, 0})
	v := testVArray([4]float32{1, 2, 3, 0})
	for i := 0; i < b.N; i++ {
		u.Dot(v)
	}
}

func (u testVector) Dot(v testVector) float32 {
	return u.X*v.X + u.Y*v.Y + u.Z*v.Z
}

func BenchmarkDotStruct(b *testing.B) {
	u := testVector{1, 2, 3}
	v := testVector{1, 2, 3}
	for i := 0; i < b.N; i++ {
		u.Dot(v)
	}
}

// benchmarking min
// verdict: simd wayy faster
func BenchmarkMinSimd(b *testing.B) {
	u := testVArray([4]float32{1, 5, 3, 0})
	v := testVArray([4]float32{6, 2, 4, 0})
	for i := 0; i < b.N; i++ {
		Min(u, v)
	}
}

func (u testVArray) Min(v testVArray) testVArray {
	xmin := float32(math.Min(float64(u[0]), float64(v[0])))
	ymin := float32(math.Min(float64(u[1]), float64(v[1])))
	zmin := float32(math.Min(float64(u[2]), float64(v[2])))
	return [4]float32{xmin, ymin, zmin, 0}
}

func BenchmarkMinArray(b *testing.B) {
	u := testVArray([4]float32{1, 5, 3, 0})
	v := testVArray([4]float32{6, 2, 4, 0})
	for i := 0; i < b.N; i++ {
		u.Min(v)
	}
}

func (u testVector) Min(v testVector) testVector {
	xmin := float32(math.Min(float64(u.X), float64(v.X)))
	ymin := float32(math.Min(float64(u.Y), float64(v.Y)))
	zmin := float32(math.Min(float64(u.Z), float64(v.Z)))
	return testVector{xmin, ymin, zmin}
}

func BenchmarkMinStruct(b *testing.B) {
	u := testVector{1, 5, 3}
	v := testVector{6, 2, 4}
	for i := 0; i < b.N; i++ {
		u.Min(v)
	}
}
