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

// benchmarking cross product
func BenchmarkCrossSimd(b *testing.B) {
	u := testVArray([4]float32{1, 2, 3, 0})
	v := testVArray([4]float32{1, 2, 3, 0})
	for i := 0; i < b.N; i++ {
		Cross(u, v)
	}
}

func (u testVArray) Cross(v testVArray) [4]float32 {
	return [4]float32{
		u[1]*v[2] - u[2]*v[1],
		u[2]*v[0] - u[0]*v[2],
		u[0]*v[1] - u[1]*v[0],
	}
}

func BenchmarkCrossArray(b *testing.B) {
	u := testVArray([4]float32{1, 2, 3, 0})
	v := testVArray([4]float32{1, 2, 3, 0})
	for i := 0; i < b.N; i++ {
		u.Cross(v)
	}
}

func (u testVector) Cross(v testVector) testVector {
	return testVector{
		X: u.Y*v.Z - u.Z*v.Y,
		Y: u.Z*v.X - u.X*v.Z,
		Z: u.X*v.Y - u.Y*v.X,
	}
}

func BenchmarkCrossStruct(b *testing.B) {
	u := testVector{1, 2, 3}
	v := testVector{1, 2, 3}
	for i := 0; i < b.N; i++ {
		u.Cross(v)
	}
}

// benchmarking min
// verdict: simd wayy faster but only if using preconverted data
func BenchmarkMinSimd(b *testing.B) {
	u := testVector{1, 5, 3}
	v := testVector{6, 2, 4}
	for i := 0; i < b.N; i++ {
		uf := testVArray([4]float32{u.X, u.Y, u.Z, 0})
		vf := testVArray([4]float32{v.X, v.Y, v.Z, 0})
		Min(uf, vf)
	}
}

func BenchmarkMinSimdNoConversion(b *testing.B) {
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

// benchmarking cube intersection test:
// simd way faster but only if using preconverted data
type testCube struct {
	mins testVector
	maxs testVector
}

func intersectSimd(origins, dirs, mins, maxs [4]float32) (float32, bool) {
	nears, fars := BoxIntersect(origins, dirs, mins, maxs)

	var t0, t1 float32 = 0.0, math.MaxFloat32
	for dim := 0; dim < 3; dim++ {
		if nears[dim] > t0 {
			t0 = nears[dim]
		}
		if fars[dim] < t1 {
			t1 = fars[dim]
		}
		if t0 > t1 {
			return 0, false
		}
	}
	return t0, true
}

func BenchmarkBoxIntersectSimd(b *testing.B) {
	rayOrigin := testVector{0, 0, 0}
	rayDirection := testVector{1, 1, 1}
	cube := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	for i := 0; i < b.N; i++ {
		origins := [4]float32{rayOrigin.X, rayOrigin.Y, rayOrigin.Z, 0}
		dirs := [4]float32{rayDirection.X, rayDirection.Y, rayDirection.Z, 0}
		mins := [4]float32{cube.mins.X, cube.mins.Y, cube.mins.Z, 0}
		maxs := [4]float32{cube.maxs.X, cube.maxs.Y, cube.maxs.Z, 0}
		intersectSimd(origins, dirs, mins, maxs)
	}
}

func BenchmarkBoxIntersectSimdNoConversion(b *testing.B) {
	rayOrigin := testVector{0, 0, 0}
	rayDirection := testVector{1, 1, 1}
	cube := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	origins := [4]float32{rayOrigin.X, rayOrigin.Y, rayOrigin.Z, 0}
	dirs := [4]float32{rayDirection.X, rayDirection.Y, rayDirection.Z, 0}
	mins := [4]float32{cube.mins.X, cube.mins.Y, cube.mins.Z, 0}
	maxs := [4]float32{cube.maxs.X, cube.maxs.Y, cube.maxs.Z, 0}
	for i := 0; i < b.N; i++ {
		intersectSimd(origins, dirs, mins, maxs)
	}
}

func (b testCube) intersect(rayOrigin, rayDirection testVector) (float32, bool) {
	var t0 float32 = 0.0
	var t1 float32 = math.MaxFloat32
	invRayDirs := [3]float32{1.0 / rayDirection.X, 1.0 / rayDirection.Y, 1.0 / rayDirection.Z}
	rayOrigins := [3]float32{rayOrigin.X, rayOrigin.Y, rayOrigin.Z}
	bPmins := [3]float32{b.mins.X, b.mins.Y, b.mins.Z}
	bPmaxs := [3]float32{b.maxs.X, b.maxs.Y, b.maxs.Z}

	for dim := 0; dim < 3; dim++ {
		tNear := (bPmins[dim] - rayOrigins[dim]) * invRayDirs[dim]
		tFar := (bPmaxs[dim] - rayOrigins[dim]) * invRayDirs[dim]
		if tNear > tFar {
			tNear, tFar = tFar, tNear
		}
		if tNear > t0 {
			t0 = tNear
		}
		if tFar < t1 {
			t1 = tFar
		}
		if t0 > t1 {
			return 0, false
		}
	}
	return t0, true
}

func BenchmarkBoxIntersectStruct(b *testing.B) {
	rayOrigin := testVector{0, 0, 0}
	rayDirection := testVector{1, 1, 1}
	cube := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	for i := 0; i < b.N; i++ {
		cube.intersect(rayOrigin, rayDirection)
	}
}

func box4intersectStruct(cube1, cube2, cube3, cube4 testCube, rayOrigin, rayDirection testVector) (float32, bool) {
	var t0 float32 = 0.0
	d, ok := cube1.intersect(rayOrigin, rayDirection)
	if !ok {
		return 0.0, false
	}
	if t0 < d {
		t0 = d
	}
	d, ok = cube2.intersect(rayOrigin, rayDirection)
	if !ok {
		return 0.0, false
	}
	if t0 < d {
		t0 = d
	}
	d, ok = cube3.intersect(rayOrigin, rayDirection)
	if !ok {
		return 0.0, false
	}
	if t0 < d {
		t0 = d
	}
	d, ok = cube4.intersect(rayOrigin, rayDirection)
	if !ok {
		return 0.0, false
	}
	if t0 < d {
		t0 = d
	}
	return t0, true
}

func Benchmark4BoxIntersectStruct(b *testing.B) {
	rayOrigin := testVector{0, 0, 0}
	rayDirection := testVector{1, 1, 1}
	cube := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	for i := 0; i < b.N; i++ {
		box4intersectStruct(cube, cube, cube, cube, rayOrigin, rayDirection)
	}
}

// no noconversion version, since we always need to calculate these on the fly
// (except mayybe cube min/maxs if we know order in the tree?)
// NOTE: building the simd friendly ray data is shared between this and triangle intersects
// so possibly reduced 5 times (once for box intersects, once for the 4 triangles if its a leaf node)
// Running ahead of myself (see mBVH paper), for triangle data a cache can be used
func box4intersectSimd(cube1, cube2, cube3, cube4 testCube, rayOrigin, rayDirection testVector) (float32, bool) {
	o4x := [4]float32{rayOrigin.X, rayOrigin.X, rayOrigin.X, rayOrigin.X}
	o4y := [4]float32{rayOrigin.Y, rayOrigin.Y, rayOrigin.Y, rayOrigin.Y}
	o4z := [4]float32{rayOrigin.Z, rayOrigin.Z, rayOrigin.Z, rayOrigin.Z}
	d4x := [4]float32{rayDirection.X, rayDirection.X, rayDirection.X, rayDirection.X}
	d4y := [4]float32{rayDirection.Y, rayDirection.Y, rayDirection.Y, rayDirection.Y}
	d4z := [4]float32{rayDirection.Z, rayDirection.Z, rayDirection.Z, rayDirection.Z}
	min4x := [4]float32{cube1.mins.X, cube2.mins.X, cube3.mins.X, cube4.mins.X}
	min4y := [4]float32{cube1.mins.Y, cube2.mins.Y, cube3.mins.Y, cube4.mins.Y}
	min4z := [4]float32{cube1.mins.Z, cube2.mins.Z, cube3.mins.Z, cube4.mins.Z}
	max4x := [4]float32{cube1.maxs.X, cube2.maxs.X, cube3.maxs.X, cube4.maxs.X}
	max4y := [4]float32{cube1.maxs.Y, cube2.maxs.Y, cube3.maxs.Y, cube4.maxs.Y}
	max4z := [4]float32{cube1.maxs.Z, cube2.maxs.Z, cube3.maxs.Z, cube4.maxs.Z}
	t0s := Box4Intersect(o4x, o4y, o4z, d4x, d4y, d4z, min4x, min4y, min4z, max4x, max4y, max4z)
	var t0 float32 = 0.0
	for _, ti := range t0s {
		if ti > t0 {
			t0 = ti
		}
	}
	return t0, t0 != 0.0
}

func Benchmark4BoxIntersectSimd(b *testing.B) {
	rayOrigin := testVector{0, 0, 0}
	rayDirection := testVector{1, 1, 1}
	cube := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	for i := 0; i < b.N; i++ {
		box4intersectSimd(cube, cube, cube, cube, rayOrigin, rayDirection)
	}
}

func Test4BoxSimd(t *testing.T) {
	rayOrigin := testVector{0, 0, 0}
	rayDirection := testVector{1, 1, 1}
	cube1 := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	cube2 := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	cube3 := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	cube4 := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	got, gotHit := box4intersectSimd(cube1, cube2, cube3, cube4, rayOrigin, rayDirection)
	wantHit := true
	var want float32 = 2
	if wantHit != gotHit {
		t.Error("hit bool incorrect")
	}
	if math.Abs(float64(got-want)) > 0.001 {
		t.Errorf("got %f but want %f as t0", got, want)
	}
}
