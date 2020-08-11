package model

import (
    "math"
    "math/rand"
)

// TODO: Textures. for starters, diffuse material can be a texture always returning same color
// PBRT calls this a ConstantTexture. Paves the way for generated texture patterns!

type Material interface {
	GetColor(*SurfaceInteraction) Color
	IsLight() bool
    Sample(r *rand.Rand, normal Vector) Vector
}

type material struct{}

func (material) IsLight() bool {
	return false
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
	incident Vector
	depth    int
	tracer   Tracer
}

func NewSurfaceInteraction(o Object, d float32, n Vector, r Ray) *SurfaceInteraction {
	return &SurfaceInteraction{
		object:   o,
		distance: d,
		normal:   n,
		ray:      r,
		Point:    PointFromRay(r, d),
		incident: r.Direction,
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
	Color Color
}

func (m *DiffuseMaterial) GetColor(si *SurfaceInteraction) Color {
	return m.Color
}

type RadiantMaterial struct {
	material
	Color Color
}

func (r *RadiantMaterial) GetColor(si *SurfaceInteraction) Color {
	facingRatio := si.normal.Dot(si.incident.Times(-1))
	if facingRatio <= 0 {
		return BLACK
	}
	return r.Color.Times(facingRatio)
}

func (*RadiantMaterial) IsLight() bool {
	return true
}

type ReflectiveMaterial struct {
	material
	Scene *Scene
}

func (m *ReflectiveMaterial) GetColor(si *SurfaceInteraction) Color {
	i := si.incident
	n := si.object.SurfaceNormal(si.Point)
	reflection := i.Sub(n.Times(2 * i.Dot(n)))
	ray := NewRay(si.Point, reflection)
	// TODO: retain maxdistance for tracing
	return si.tracer.GetRayColor(ray, m.Scene, si.depth+1) //.Times(1 - standardAlbedo) // simulates nonperfect reflection
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

// temporary material to play around with
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
