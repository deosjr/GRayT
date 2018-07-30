package model

type Material interface {
	GetColor(*SurfaceInteraction) Color
}

type SurfaceInteraction struct {
	Point    Vector
	Object   Object
	AS       AccelerationStructure
	Incident Vector
	depth    int
}

type DiffuseMaterial struct {
	Color Color
}

func (m *DiffuseMaterial) GetColor(*SurfaceInteraction) Color {
	return m.Color
}

// temporary material to play around with
type PosFuncMat struct {
	Func func(Vector) Color
}

func (m *PosFuncMat) GetColor(si *SurfaceInteraction) Color {
	return m.Func(si.Point)
}

type ReflectiveMaterial struct {
	Scene *Scene
}

var maxReflectiveDepth = 5

func (m *ReflectiveMaterial) GetColor(si *SurfaceInteraction) Color {
	if si.depth == maxReflectiveDepth {
		return BACKGROUND_COLOR
	}
	i := si.Incident
	n := si.Object.SurfaceNormal(si.Point)
	reflection := i.Sub(n.Times(2 * i.Dot(n)))
	ray := NewRay(si.Point, reflection)
	return m.Scene.GetRayColor(ray)
}
