package model

type Material interface {
	GetColor(*SurfaceInteraction) Color
	IsLight() bool
}

type material struct{}

func (material) IsLight() bool {
	return false
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
