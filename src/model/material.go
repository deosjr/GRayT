package model

import "math"

type Material interface {
	GetColor(*SurfaceInteraction, Light) Color
}

type SurfaceInteraction struct {
	distance float64
	ray      Ray
	point    Vector
	normal   Vector
	object   Object
	as       AccelerationStructure
	incident Vector
	depth    int
	tracer   Tracer
}

func NewSurfaceInteraction(o Object, d float64, n Vector, r Ray) *SurfaceInteraction {
	return &SurfaceInteraction{
		object:   o,
		distance: d,
		normal:   n,
		ray:      r,
		point:    PointFromRay(r, d),
		incident: r.Direction,
	}
}

type DiffuseMaterial struct {
	Color Color
}

func (m *DiffuseMaterial) GetColor(si *SurfaceInteraction, l Light) Color {
	facingRatio := si.normal.Dot(si.incident.Times(-1))
	lightSegment := l.GetLightSegment(si.point)
	lightRatio := si.normal.Dot(lightSegment.Normalize())
	factors := standardAlbedo / math.Pi * l.Intensity(lightSegment.Length()) * facingRatio * lightRatio
	lightColor := l.Color().Times(factors)
	return m.Color.Product(lightColor)
}

type RadiantMaterial struct {
	Color Color
}

func (r *RadiantMaterial) GetColor(si *SurfaceInteraction, l Light) Color {
	return r.Color
}

type ReflectiveMaterial struct {
	Scene *Scene
}

func (m *ReflectiveMaterial) GetColor(si *SurfaceInteraction, l Light) Color {
	i := si.incident
	n := si.object.SurfaceNormal(si.point)
	reflection := i.Sub(n.Times(2 * i.Dot(n)))
	ray := NewRay(si.point, reflection)
	// TODO: retain maxdistance for tracing
	return si.tracer.GetRayColor(ray, m.Scene, si.depth+1) //.Times(1 - standardAlbedo) // simulates nonperfect reflection
}

type NormalMappingMaterial struct {
	WrappedMaterial Material
	NormalFunc      func(*SurfaceInteraction) Vector
}

func (m *NormalMappingMaterial) GetColor(si *SurfaceInteraction, l Light) Color {
	si.normal = m.NormalFunc(si)
	return m.WrappedMaterial.GetColor(si, l)
}

// temporary material to play around with
type PosFuncMat struct {
	Func func(*SurfaceInteraction, Light) Color
}

func (m *PosFuncMat) GetColor(si *SurfaceInteraction, l Light) Color {
	return m.Func(si, l)
}

var DebugNormalMaterial = &PosFuncMat{
	Func: func(si *SurfaceInteraction, _ Light) Color {
		n := si.normal.Times(0.5).Add(Vector{0.5, 0.5, 0.5})
		return Color{n.X, n.Y, n.Z}
	},
}
