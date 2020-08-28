package model

import (
	"math"
	"math/rand"
)

// TODO: currently src/model/tracer.go has a lot of switch cases for materials.
// these should be consolidated as material methods in some form again
// unsupported (properly) in path tracers: reflective material

type Material interface {
	IsLight() bool
	Sample(r *rand.Rand, normal Vector) Vector
	GetColor(si *SurfaceInteraction) Color
}

type material struct {
	texture Texture
}

func (material) IsLight() bool {
	return false
}

func (m material) GetColor(si *SurfaceInteraction) Color {
	return m.texture.GetColor(si)
}

// default sampling for all material right now is the same
func (material) Sample(r *rand.Rand, normal Vector) Vector {
	return randomInHemisphere(r, normal)
}

// this is actually slower than the very naive method before..
func randomInHemisphere(random *rand.Rand, normal Vector) Vector {
	// uniform hemisphere sampling: pbrt 774
	// samples from hemisphere with z-axis = up direction
	z := random.Float64()
	det := 1 - z*z
	var r float64 = 0.0
	if det > 0 {
		r = math.Sqrt(det)
	}
	phi := 2 * math.Pi * random.Float64()
	v := Vector{float32(r * math.Cos(phi)), float32(r * math.Sin(phi)), float32(z)}

	ez := Vector{0, 0, 1}
	rotationVector := ez.Cross(normal)
	theta := math.Acos(float64(ez.Dot(normal)))
	return Rotate(theta, rotationVector).Vector(v)
}

type SurfaceInteraction struct {
	distance float32
	ray      Ray
	Point    Vector
	normal   Vector
	object   Object
	as       AccelerationStructure
	depth    int
	tracer   Tracer

	// for sharedobjects
	UntransformedPoint  Vector
	UntransformedNormal Vector
}

func NewSurfaceInteraction(o Object, d float32, n Vector, r Ray) *SurfaceInteraction {
	p := PointFromRay(r, d)
	return &SurfaceInteraction{
		object:   o,
		distance: d,
		normal:   n,
		ray:      r,
		Point:    p,
		// can be overriden by sharedobjects
		UntransformedPoint:  p,
		UntransformedNormal: n,
	}
}

func (si *SurfaceInteraction) GetNormal() Vector {
	return si.normal
}

func (si *SurfaceInteraction) GetObject() Object {
	return si.object
}

type DiffuseMaterial struct {
	material
}

func NewDiffuseMaterial(t Texture) *DiffuseMaterial {
	return &DiffuseMaterial{
		material: material{
			texture: t,
		},
	}
}

type RadiantMaterial struct {
	material
}

func NewRadiantMaterial(t Texture) *RadiantMaterial {
	return &RadiantMaterial{
		material: material{
			texture: t,
		},
	}
}

func (*RadiantMaterial) IsLight() bool {
	return true
}

type ReflectiveMaterial struct {
	material
	Scene *Scene
}

type NormalMappingMaterial struct {
	material
	WrappedMaterial Material
	NormalFunc      func(*SurfaceInteraction) Vector
}

// TODO: this is a bit of a hack, no? where should this normal mapping happen?
func (m *NormalMappingMaterial) GetColor(si *SurfaceInteraction) Color {
	si.normal = m.NormalFunc(si)
	return m.WrappedMaterial.GetColor(si)
}

// only works for triangles in mesh
func InterpolatedNormalMappingMaterial(mat Material) *NormalMappingMaterial {
	return &NormalMappingMaterial{
		WrappedMaterial: mat,
		NormalFunc: func(si *SurfaceInteraction) Vector {
			tr := si.GetObject().(TriangleInMesh)
			p := si.UntransformedPoint
			l0, l1, l2 := tr.Barycentric(p)
			p0, p1, p2 := tr.PointIndices()
			nl0 := tr.Mesh.Normals[p0]
			nl1 := tr.Mesh.Normals[p1]
			nl2 := tr.Mesh.Normals[p2]
			return nl0.Times(l0).Add(nl1.Times(l1)).Add(nl2.Times(l2))
		},
	}
}

// used in whitted style raytracer to ignore light contribution when debugging
type PosFuncMat struct {
	material
	Func func(*SurfaceInteraction) Color
}

func (m *PosFuncMat) GetColor(si *SurfaceInteraction) Color {
	return m.Func(si)
}

var DebugNormalMaterial = &PosFuncMat{
	Func: func(si *SurfaceInteraction) Color {
		n := si.normal.Times(0.5).Add(Vector{0.5, 0.5, 0.5})
		return Color{n.X, n.Y, n.Z}
	},
}
