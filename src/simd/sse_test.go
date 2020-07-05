package simd

import (
	"math"
	"math/rand"
	"testing"
	"time"
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

func box4intersectStruct(cube1, cube2, cube3, cube4 testCube, rayOrigin, rayDirection testVector) [4]float32 {
	var t0s [4]float32
	d, ok := cube1.intersect(rayOrigin, rayDirection)
	if ok {
		t0s[0] = d
	}
	d, ok = cube2.intersect(rayOrigin, rayDirection)
	if ok {
		t0s[1] = d
	}
	d, ok = cube3.intersect(rayOrigin, rayDirection)
	if ok {
		t0s[2] = d
	}
	d, ok = cube4.intersect(rayOrigin, rayDirection)
	if ok {
		t0s[3] = d
	}
	return t0s
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
func box4intersectSimd(cube1, cube2, cube3, cube4 testCube, rayOrigin, rayDirection testVector) [4]float32 {
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
	return Box4Intersect(o4x, o4y, o4z, d4x, d4y, d4z, min4x, min4y, min4z, max4x, max4y, max4z)
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
	for d, tt := range []struct {
		ro testVector
		rd testVector
		c1 testCube
		c2 testCube
		c3 testCube
		c4 testCube
	}{
		{
			ro: testVector{0, 0, 0},
			rd: testVector{1, 1, 1},
			c1: testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}},
			c2: testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}},
			c3: testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}},
			c4: testCube{mins: testVector{2, 2, 2}, maxs: testVector{3, 3, 3}},
		},
		// random generated test, fixed by adding t0=0.0 at start of calculations
		// i.e. XORPS the first register and add MAXPS with first t0 value
		{
			ro: testVector{0, 0, 0},
			rd: testVector{-3.2323341, -33.804115, 29.499321},
			c1: testCube{mins: testVector{18.107834, -25.848492, -18.847755}, maxs: testVector{-30.707195, 13.312214, 45.28969}},
			c2: testCube{mins: testVector{29.642372, -6.9531517, -47.928432}, maxs: testVector{-11.721863, -8.825039, 17.887497}},
			c3: testCube{mins: testVector{-10.571224, 4.5797653, -24.498928}, maxs: testVector{39.838516, -2.975563, -2.0600662}},
			c4: testCube{mins: testVector{16.159958, 30.34137, -41.693176}, maxs: testVector{-14.749153, -45.38565, 9.922146}},
		},
		// random generated test
		{
			ro: testVector{0, 0, 0},
			rd: testVector{-28.309319, -24.20735, 31.777344},
			c1: testCube{mins: testVector{33.618423, -12.316433, -12.135181}, maxs: testVector{-35.066486, 12.27599, -6.2505836}},
			c2: testCube{mins: testVector{-40.74101, 1.7597961, -21.449244}, maxs: testVector{-22.05585, -35.47505, 45.336975}},
			c3: testCube{mins: testVector{-39.31895, 22.843544, -13.831745}, maxs: testVector{-27.322006, 7.1916885, 9.056133}},
			c4: testCube{mins: testVector{43.012024, -47.114372, -29.980309}, maxs: testVector{39.411858, 24.975937, 31.67485}},
		},
	} {
		want := box4intersectStruct(tt.c1, tt.c2, tt.c3, tt.c4, tt.ro, tt.rd)
		got := box4intersectSimd(tt.c1, tt.c2, tt.c3, tt.c4, tt.ro, tt.rd)
		for i := range want {
			if math.Abs(float64(got[i]-want[i])) > 0.001 {
				t.Errorf("%d) got %f but want %f as t0s", d, got, want)
				break
			}
		}
	}
}

func TestRandom4BoxSimd(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	f := func() float32 {
		return rand.Float32()*100 - 50
	}

	for {
		rayOrigin := testVector{0, 0, 0}
		rayDirection := testVector{f(), f(), f()}
		cube1 := testCube{mins: testVector{f(), f(), f()}, maxs: testVector{f(), f(), f()}}
		cube2 := testCube{mins: testVector{f(), f(), f()}, maxs: testVector{f(), f(), f()}}
		cube3 := testCube{mins: testVector{f(), f(), f()}, maxs: testVector{f(), f(), f()}}
		cube4 := testCube{mins: testVector{f(), f(), f()}, maxs: testVector{f(), f(), f()}}

		got := box4intersectStruct(cube1, cube2, cube3, cube4, rayOrigin, rayDirection)

		gotSimd := box4intersectSimd(cube1, cube2, cube3, cube4, rayOrigin, rayDirection)

		if got != gotSimd {
			t.Errorf("got %v gotSimd %v", got, gotSimd)
			t.Errorf("ro %v rd %v c1 %v c2 %v c3 %v c4 %v", rayOrigin, rayDirection, cube1, cube2, cube3, cube4)
			break
		}
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
		// 2 randomly generated test cases that fail on rounding errors
		// when using RCPPS
		{
			p0: testVector{X: -281.573, Y: 410.15726, Z: -141.82341},
			p1: testVector{X: 351.67438, Y: 434.09607, Z: 440.06393},
			p2: testVector{X: -125.7391, Y: -127.68364, Z: 318.1397},
			ro: testVector{X: 174.38438, Y: 304.48764, Z: 272.30704},
			rd: testVector{X: 0.2938059, Y: -0.20797282, Z: 0.9329659},
		},
		{
			p0: testVector{X: -63.931763, Y: -56.17392, Z: -442.17444},
			p1: testVector{X: -420.88693, Y: -469.2147, Z: 304.75366},
			p2: testVector{X: -396.92105, Y: 244.8274, Z: -53.85065},
			ro: testVector{X: -114.9804, Y: -66.87874, Z: 238.39647},
			rd: testVector{X: -0.7497336, Y: -0.6617179, Z: -0.0053925477},
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
	for d, tt := range []struct {
		triangles [4]testTriangle
		ro        testVector
		rd        testVector
	}{
		{
			triangles: [4]testTriangle{
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
			},
			ro: testVector{0, 0, 0},
			rd: testVector{0, 0, -1},
		},
		// 2 randomly generated tests that fail rounding errors
		// when using RCPPS
		{
			triangles: [4]testTriangle{
				{
					p0: testVector{X: -460.37598, Y: -324.83853, Z: 33.135204},
					p1: testVector{X: -299.03824, Y: 229.80937, Z: 30.185402},
					p2: testVector{X: 280.91913, Y: 301.7894, Z: 316.75595},
				},
				{
					p0: testVector{X: 154.19307, Y: -309.96835, Z: 367.1002},
					p1: testVector{X: 132.89345, Y: -417.66684, Z: -19.636272},
					p2: testVector{X: 220.65219, Y: -303.69104, Z: 69.31779},
				},
				{
					p0: testVector{X: -281.573, Y: 410.15726, Z: -141.82341},
					p1: testVector{X: 351.67438, Y: 434.09607, Z: 440.06393},
					p2: testVector{X: -125.7391, Y: -127.68364, Z: 318.1397},
				},
				{
					p0: testVector{X: -87.283554, Y: 354.5092, Z: 102.55295},
					p1: testVector{X: -248.82996, Y: 164.36073, Z: 357.0556},
					p2: testVector{X: 274.93723, Y: 408.07104, Z: -367.10382},
				},
			},
			ro: testVector{X: 174.38438, Y: 304.48764, Z: 272.30704},
			rd: testVector{X: 0.2938059, Y: -0.20797282, Z: 0.9329659},
		},
		{
			triangles: [4]testTriangle{
				{
					p0: testVector{X: -63.931763, Y: -56.17392, Z: -442.17444},
					p1: testVector{X: -420.88693, Y: -469.2147, Z: 304.75366},
					p2: testVector{X: -396.92105, Y: 244.8274, Z: -53.85065},
				},
				{
					p0: testVector{X: -435.83505, Y: -152.09586, Z: 348.71805},
					p1: testVector{X: 262.57355, Y: 355.4922, Z: 183.94792},
					p2: testVector{X: -308.7444, Y: -129.70721, Z: -449.11},
				},
				{
					p0: testVector{X: 335.11725, Y: -427.42883, Z: 373.0301},
					p1: testVector{X: 472.27707, Y: 99.24653, Z: -127.22689},
					p2: testVector{X: 466.58768, Y: -464.82974, Z: -466.27313},
				},
				{
					p0: testVector{X: 496.35965, Y: -27.85486, Z: -202.64835},
					p1: testVector{X: -351.3991, Y: -428.6788, Z: 364.8578},
					p2: testVector{X: 57.72078, Y: 35.152046, Z: 78.13752},
				},
			},
			ro: testVector{X: -114.9804, Y: -66.87874, Z: 238.39647},
			rd: testVector{X: -0.7497336, Y: -0.6617179, Z: -0.0053925477},
		},
	} {
		out := [4]float32{}
		hit := [4]bool{}
		var p0x, p0y, p0z, p1x, p1y, p1z, p2x, p2y, p2z [4]float32
		for i, tr := range tt.triangles {
			out[i], hit[i] = tr.intersect(tt.ro, tt.rd)
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
		rox := [4]float32{tt.ro.X, tt.ro.X, tt.ro.X, tt.ro.X}
		roy := [4]float32{tt.ro.Y, tt.ro.Y, tt.ro.Y, tt.ro.Y}
		roz := [4]float32{tt.ro.Z, tt.ro.Z, tt.ro.Z, tt.ro.Z}
		rdx := [4]float32{tt.rd.X, tt.rd.X, tt.rd.X, tt.rd.X}
		rdy := [4]float32{tt.rd.Y, tt.rd.Y, tt.rd.Y, tt.rd.Y}
		rdz := [4]float32{tt.rd.Z, tt.rd.Z, tt.rd.Z, tt.rd.Z}
		outf := Triangle4Intersect(p0x, p0y, p0z, p1x, p1y, p1z, p2x, p2y, p2z, rox, roy, roz, rdx, rdy, rdz)
		for i := 0; i < 4; i++ {
			hitf := outf[i] != 0.0

			if hitf != hit[i] {
				t.Errorf("%d-%d): got %v want %v", d, i, hitf, hit[i])
			}
			if hit[i] && math.Abs(float64(out[i]-outf[i])) > 0.001 {
				t.Errorf("%d-%d): got %e want %e", d, i, outf[i], out[i])
			}
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
