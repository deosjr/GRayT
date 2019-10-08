package simd

import (
	"math"
	"math/rand"
	"testing"
)

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

// benchmarking dot product
func BenchmarkDotSimd(b *testing.B) {
	u := testVector{1, 2, 3}
	v := testVector{1, 2, 3}
	uf := [4]float32{u.X, u.Y, u.Z, 0}
	vf := [4]float32{v.X, v.Y, v.Z, 0}
	for i := 0; i < b.N; i++ {
		Dot(uf, vf)
	}
}

func TestDotProductSimd(t *testing.T) {
	u := testVector{1, 2, 3}
	v := testVector{1, 2, 3}
	uf := [4]float32{u.X, u.Y, u.Z, 0}
	vf := [4]float32{v.X, v.Y, v.Z, 0}
	got := Dot(uf, vf)
	if got != 14.0 {
		t.Errorf("got %f but wanted %f", got, 14.0)
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

func Benchmark4BoxIntersectSimdNoConversion(b *testing.B) {
	rayOrigin := testVector{0, 0, 0}
	rayDirection := testVector{1, 1, 1}
	cube1 := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	cube2 := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	cube3 := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
	cube4 := testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}}
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
	for i := 0; i < b.N; i++ {
		Box4Intersect(o4x, o4y, o4z, d4x, d4y, d4z, min4x, min4y, min4z, max4x, max4y, max4z)
	}
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

// benchmarking normalize function
func (u testVector) normalize() testVector {
	l := float32(math.Sqrt(float64(u.X*u.X + u.Y*u.Y + u.Z*u.Z)))
	if l == 0 {
		return testVector{0, 0, 0}
	}
	r := 1.0 / l
	return testVector{u.X * r, u.Y * r, u.Z * r}
}

func normalize4(a, b, c, d testVector) (e, f, g, h testVector) {
	x4 := [4]float32{a.X, b.X, c.X, d.X}
	y4 := [4]float32{a.Y, b.Y, c.Y, d.Y}
	z4 := [4]float32{a.Z, b.Z, c.Z, d.Z}
	Normalize4(x4, y4, z4)
	e = testVector{x4[0], y4[0], z4[0]}
	f = testVector{x4[1], y4[1], z4[1]}
	g = testVector{x4[2], y4[2], z4[2]}
	h = testVector{x4[3], y4[3], z4[3]}
	return
}

func BenchmarkNormalizeStruct(b *testing.B) {
	v := testVector{4, 3, 2}
	for i := 0; i < b.N; i++ {
		v.normalize()
	}
}

func BenchmarkNormalizeSimd(b *testing.B) {
	v := testVector{4, 3, 2}
	for i := 0; i < b.N; i++ {
		vf := [4]float32{v.X, v.Y, v.Z, 0}
		Normalize(vf)
	}
}

func Benchmark4NormalizeStruct(b *testing.B) {
	v1 := testVector{rand.Float32(), rand.Float32(), rand.Float32()}
	v2 := testVector{rand.Float32(), rand.Float32(), rand.Float32()}
	v3 := testVector{rand.Float32(), rand.Float32(), rand.Float32()}
	v4 := testVector{rand.Float32(), rand.Float32(), rand.Float32()}
	for i := 0; i < b.N; i++ {
		v1.normalize()
		v2.normalize()
		v3.normalize()
		v4.normalize()
	}
}

func Benchmark4NormalizeSimd(b *testing.B) {
	v1 := testVector{rand.Float32(), rand.Float32(), rand.Float32()}
	v2 := testVector{rand.Float32(), rand.Float32(), rand.Float32()}
	v3 := testVector{rand.Float32(), rand.Float32(), rand.Float32()}
	v4 := testVector{rand.Float32(), rand.Float32(), rand.Float32()}
	for i := 0; i < b.N; i++ {
		normalize4(v1, v2, v3, v4)
	}
}

type testTriangle struct {
	p0, p1, p2 testVector
}

// Moller-Trumbore intersection algorithm
func (t testTriangle) intersect(rayOrigin, rayDirection testVector) (float32, bool) {
	e1 := t.p1.Sub(t.p0)
	e2 := t.p2.Sub(t.p0)
	pvec := rayDirection.Cross(e2)
	det := e1.Dot(pvec)

	if det < 1e-8 && det > -1e-8 {
		return 0, false
	}
	inv_det := 1.0 / det

	tvec := rayOrigin.Sub(t.p0)
	u := tvec.Dot(pvec) * inv_det
	if u < 0 || u > 1 {
		return 0, false
	}

	qvec := tvec.Cross(e1)
	v := rayDirection.Dot(qvec) * inv_det
	if v < 0 || u+v > 1 {
		return 0, false
	}
	return e2.Dot(qvec) * inv_det, true
}

func (u testVector) Sub(v testVector) testVector {
	return testVector{
		X: u.X - v.X,
		Y: u.Y - v.Y,
		Z: u.Z - v.Z,
	}
}

func TestTriangleIntersect(t *testing.T) {
	for i, tt := range []struct {
		ro testVector
		rd testVector
		p0 testVector
		p1 testVector
		p2 testVector
	}{
		{
			ro: testVector{0, 0, 0},
			rd: testVector{0, 0, -1},
			p0: testVector{-1, 0, -1},
			p1: testVector{1, 0, -1},
			p2: testVector{1, 1, -1},
		},
		{
			ro: testVector{5, 5, 5},
			rd: testVector{0, 0, -1},
			p0: testVector{-1, 0, -1},
			p1: testVector{1, 0, -1},
			p2: testVector{1, 1, -1},
		},
		{
			ro: testVector{0, 0, 0},
			rd: testVector{0, 0, -1},
			p0: testVector{-2, 0, -2},
			p1: testVector{2, 0, -2},
			p2: testVector{2, 2, -2},
		},
	} {
		triangle := testTriangle{tt.p0, tt.p1, tt.p2}
		out, hit := triangle.intersect(tt.ro, tt.rd)

		rof := [4]float32{tt.ro.X, tt.ro.Y, tt.ro.Z, 0}
		rdf := [4]float32{tt.rd.X, tt.rd.Y, tt.rd.Z, 0}
		p0f := [4]float32{tt.p0.X, tt.p0.Y, tt.p0.Z, 0}
		p1f := [4]float32{tt.p1.X, tt.p1.Y, tt.p1.Z, 0}
		p2f := [4]float32{tt.p2.X, tt.p2.Y, tt.p2.Z, 0}
		outf := TriangleIntersect(p0f, p1f, p2f, rof, rdf)
		hitf := outf != 0.0

		if hitf != hit {
			t.Errorf("%d): got %t want %t", i, hitf, hit)
		}
		if hit && math.Abs(float64(out-outf)) > 0.001 {
			t.Errorf("%d): got %e want %e", i, outf, out)
		}
	}
}

func BenchmarkTriangleIntersectStruct(b *testing.B) {
	ro := testVector{0, 0, 0}
	rd := testVector{0, 0, -1}
	p0 := testVector{-1, 0, -1}
	p1 := testVector{1, 0, -1}
	p2 := testVector{1, 1, -1}
	triangle := testTriangle{p0, p1, p2}
	for i := 0; i < b.N; i++ {
		triangle.intersect(ro, rd)
	}
}

func BenchmarkTriangleIntersectSimd(b *testing.B) {
	ro := testVector{0, 0, 0}
	rd := testVector{0, 0, -1}
	p0 := testVector{-1, 0, -1}
	p1 := testVector{1, 0, -1}
	p2 := testVector{1, 1, -1}
	for i := 0; i < b.N; i++ {
		rof := [4]float32{ro.X, ro.Y, ro.Z, 0}
		rdf := [4]float32{rd.X, rd.Y, rd.Z, 0}
		p0f := [4]float32{p0.X, p0.Y, p0.Z, 0}
		p1f := [4]float32{p1.X, p1.Y, p1.Z, 0}
		p2f := [4]float32{p2.X, p2.Y, p2.Z, 0}
		TriangleIntersect(p0f, p1f, p2f, rof, rdf)
	}
}

func BenchmarkTriangleIntersectNoConversionSimd(b *testing.B) {
	ro := testVector{0, 0, 0}
	rd := testVector{0, 0, -1}
	p0 := testVector{-1, 0, -1}
	p1 := testVector{1, 0, -1}
	p2 := testVector{1, 1, -1}
	rof := [4]float32{ro.X, ro.Y, ro.Z, 0}
	rdf := [4]float32{rd.X, rd.Y, rd.Z, 0}
	p0f := [4]float32{p0.X, p0.Y, p0.Z, 0}
	p1f := [4]float32{p1.X, p1.Y, p1.Z, 0}
	p2f := [4]float32{p2.X, p2.Y, p2.Z, 0}
	for i := 0; i < b.N; i++ {
		TriangleIntersect(p0f, p1f, p2f, rof, rdf)
	}
}

func Test4TriangleIntersect(t *testing.T) {
	triangles := [4]testTriangle{
		{
			p0: testVector{-1, 0, -1},
			p1: testVector{1, 0, -1},
			p2: testVector{1, 1, -1},
		},
		{
			p0: testVector{-2, 0, -2},
			p1: testVector{2, 0, -2},
			p2: testVector{2, 2, -2},
		},
		{
			p0: testVector{-1, 0, -1},
			p1: testVector{1, 0, -1},
			p2: testVector{1, 1, -1},
		},
		{
			p0: testVector{-1, 0, -1},
			p1: testVector{1, 0, -1},
			p2: testVector{1, 1, -1},
		},
	}
	ro := testVector{0, 0, 0}
	rd := testVector{0, 0, -1}
	out := [4]float32{}
	hit := [4]bool{}
	var p0x, p0y, p0z, p1x, p1y, p1z, p2x, p2y, p2z [4]float32
	for i, tr := range triangles {
		out[i], hit[i] = tr.intersect(ro, rd)
		p0x[i] = tr.p0.X
		p0y[i] = tr.p0.Y
		p0z[i] = tr.p0.Z
		p1x[i] = tr.p1.X
		p1y[i] = tr.p1.Y
		p1z[i] = tr.p1.Z
		p2x[i] = tr.p2.X
		p2y[i] = tr.p2.Y
		p2z[i] = tr.p2.Z
	}
	rox := [4]float32{ro.X, ro.X, ro.X, ro.X}
	roy := [4]float32{ro.Y, ro.Y, ro.Y, ro.Y}
	roz := [4]float32{ro.Z, ro.Z, ro.Z, ro.Z}
	rdx := [4]float32{rd.X, rd.X, rd.X, rd.X}
	rdy := [4]float32{rd.Y, rd.Y, rd.Y, rd.Y}
	rdz := [4]float32{rd.Z, rd.Z, rd.Z, rd.Z}
	outf := Triangle4Intersect(p0x, p0y, p0z, p1x, p1y, p1z, p2x, p2y, p2z, rox, roy, roz, rdx, rdy, rdz)
	for i := 0; i < 4; i++ {
		hitf := outf[i] != 0.0

		if hitf != hit[i] {
			t.Errorf("%d): got %v want %v", i, hitf, hit[i])
		}
		if hit[i] && math.Abs(float64(out[i]-outf[i])) > 0.001 {
			t.Errorf("%d): got %e want %e", i, outf[i], out[i])
		}
	}
}

func Benchmark4TriangleStruct(b *testing.B) {
	ro := testVector{0, 0, 0}
	rd := testVector{0, 0, -1}
	p0 := testVector{-1, 0, -1}
	p1 := testVector{1, 0, -1}
	p2 := testVector{1, 1, -1}
	tr := testTriangle{p0, p1, p2}
	triangles := [4]testTriangle{tr, tr, tr, tr}
	for i := 0; i < b.N; i++ {
		for _, tr := range triangles {
			tr.intersect(ro, rd)
		}
	}
}

func Benchmark4TriangleIntersectSimdNoConversion(b *testing.B) {
	ro := testVector{0, 0, 0}
	rd := testVector{0, 0, -1}
	p0 := testVector{-1, 0, -1}
	p1 := testVector{1, 0, -1}
	p2 := testVector{1, 1, -1}
	tr := testTriangle{p0, p1, p2}
	triangles := [4]testTriangle{tr, tr, tr, tr}
	rox := [4]float32{ro.X, ro.X, ro.X, ro.X}
	roy := [4]float32{ro.Y, ro.Y, ro.Y, ro.Y}
	roz := [4]float32{ro.Z, ro.Z, ro.Z, ro.Z}
	rdx := [4]float32{rd.X, rd.X, rd.X, rd.X}
	rdy := [4]float32{rd.Y, rd.Y, rd.Y, rd.Y}
	rdz := [4]float32{rd.Z, rd.Z, rd.Z, rd.Z}
	p0x := [4]float32{triangles[0].p0.X, triangles[1].p0.X, triangles[2].p0.X, triangles[3].p0.X}
	p0y := [4]float32{triangles[0].p0.Y, triangles[1].p0.Y, triangles[2].p0.Y, triangles[3].p0.Y}
	p0z := [4]float32{triangles[0].p0.Z, triangles[1].p0.Z, triangles[2].p0.Z, triangles[3].p0.Z}
	p1x := [4]float32{triangles[0].p1.X, triangles[1].p1.X, triangles[2].p1.X, triangles[3].p1.X}
	p1y := [4]float32{triangles[0].p1.Y, triangles[1].p1.Y, triangles[2].p1.Y, triangles[3].p1.Y}
	p1z := [4]float32{triangles[0].p1.Z, triangles[1].p1.Z, triangles[2].p1.Z, triangles[3].p1.Z}
	p2x := [4]float32{triangles[0].p2.X, triangles[1].p2.X, triangles[2].p2.X, triangles[3].p2.X}
	p2y := [4]float32{triangles[0].p2.Y, triangles[1].p2.Y, triangles[2].p2.Y, triangles[3].p2.Y}
	p2z := [4]float32{triangles[0].p2.Z, triangles[1].p2.Z, triangles[2].p2.Z, triangles[3].p2.Z}
	for i := 0; i < b.N; i++ {
		Triangle4Intersect(p0x, p0y, p0z, p1x, p1y, p1z, p2x, p2y, p2z, rox, roy, roz, rdx, rdy, rdz)
	}
}

func Benchmark4TriangleIntersectSimd(b *testing.B) {
	ro := testVector{0, 0, 0}
	rd := testVector{0, 0, -1}
	p0 := testVector{-1, 0, -1}
	p1 := testVector{1, 0, -1}
	p2 := testVector{1, 1, -1}
	tr := testTriangle{p0, p1, p2}
	triangles := [4]testTriangle{tr, tr, tr, tr}
	for i := 0; i < b.N; i++ {
		rox := [4]float32{ro.X, ro.X, ro.X, ro.X}
		roy := [4]float32{ro.Y, ro.Y, ro.Y, ro.Y}
		roz := [4]float32{ro.Z, ro.Z, ro.Z, ro.Z}
		rdx := [4]float32{rd.X, rd.X, rd.X, rd.X}
		rdy := [4]float32{rd.Y, rd.Y, rd.Y, rd.Y}
		rdz := [4]float32{rd.Z, rd.Z, rd.Z, rd.Z}
		p0x := [4]float32{triangles[0].p0.X, triangles[1].p0.X, triangles[2].p0.X, triangles[3].p0.X}
		p0y := [4]float32{triangles[0].p0.Y, triangles[1].p0.Y, triangles[2].p0.Y, triangles[3].p0.Y}
		p0z := [4]float32{triangles[0].p0.Z, triangles[1].p0.Z, triangles[2].p0.Z, triangles[3].p0.Z}
		p1x := [4]float32{triangles[0].p1.X, triangles[1].p1.X, triangles[2].p1.X, triangles[3].p1.X}
		p1y := [4]float32{triangles[0].p1.Y, triangles[1].p1.Y, triangles[2].p1.Y, triangles[3].p1.Y}
		p1z := [4]float32{triangles[0].p1.Z, triangles[1].p1.Z, triangles[2].p1.Z, triangles[3].p1.Z}
		p2x := [4]float32{triangles[0].p2.X, triangles[1].p2.X, triangles[2].p2.X, triangles[3].p2.X}
		p2y := [4]float32{triangles[0].p2.Y, triangles[1].p2.Y, triangles[2].p2.Y, triangles[3].p2.Y}
		p2z := [4]float32{triangles[0].p2.Z, triangles[1].p2.Z, triangles[2].p2.Z, triangles[3].p2.Z}
		Triangle4Intersect(p0x, p0y, p0z, p1x, p1y, p1z, p2x, p2y, p2z, rox, roy, roz, rdx, rdy, rdz)
	}
}
