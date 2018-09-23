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
}

type DiffuseMaterial struct {
	Color Color
}

func (m *DiffuseMaterial) GetColor(si *SurfaceInteraction, l Light) Color {
	lightRatio := l.LightRatio(si.Point, si.Normal)
	facingRatio := si.Normal.Dot(si.Incident.Times(-1))
	lightSegment := l.GetLightSegment(si.Point)
	factors := standardAlbedo / math.Pi * l.Intensity(lightSegment.Length()) * facingRatio * lightRatio
	lightColor := l.Color().Times(factors)
	return m.Color.Product(lightColor)
}

type ReflectiveMaterial struct {
	Scene *Scene
}

var maxRayDepth = 5

func (m *ReflectiveMaterial) GetColor(si *SurfaceInteraction, l Light) Color {
	if si.depth == maxRayDepth {
		return BACKGROUND_COLOR
	}
	i := si.Incident
	n := si.Object.SurfaceNormal(si.Point)
	reflection := i.Sub(n.Times(2 * i.Dot(n)))
	ray := NewRay(si.Point, reflection)
	// TODO: retain maxdistance for tracing
	return m.Scene.GetRayColor(ray) //.Times(1 - standardAlbedo) // simulates nonperfect reflection
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
	Func func(Vector) Color
}

func (m *PosFuncMat) GetColor(si *SurfaceInteraction) Color {
	return m.Func(si.Point)
}
