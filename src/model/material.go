package model

import "math"

type Material interface {
	GetColor(*SurfaceInteraction, Light) Color
}

type SurfaceInteraction struct {
	Point    Vector
	Normal   Vector
	Object   Object
	AS       AccelerationStructure
	Incident Vector
	depth    int
	tracer   Tracer
}

type DiffuseMaterial struct {
	Color Color
}

func (m *DiffuseMaterial) GetColor(si *SurfaceInteraction, l Light) Color {
	facingRatio := si.Normal.Dot(si.Incident.Times(-1))
	lightSegment := l.GetLightSegment(si.Point)
	lightRatio := si.Normal.Dot(lightSegment.Normalize())
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
	i := si.Incident
	n := si.Object.SurfaceNormal(si.Point)
	reflection := i.Sub(n.Times(2 * i.Dot(n)))
	ray := NewRay(si.Point, reflection)
	// TODO: retain maxdistance for tracing
	return si.tracer.GetRayColor(ray, m.Scene, si.depth+1) //.Times(1 - standardAlbedo) // simulates nonperfect reflection
}

type NormalMappingMaterial struct {
	WrappedMaterial Material
	NormalFunc      func(*SurfaceInteraction) Vector
}

func (m *NormalMappingMaterial) GetColor(si *SurfaceInteraction, l Light) Color {
	si.Normal = m.NormalFunc(si)
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
		n := si.Normal.Times(0.5).Add(Vector{0.5, 0.5, 0.5})
		return Color{n.X, n.Y, n.Z}
	},
}
