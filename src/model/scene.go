package model

import "math"

const (
	standardAlbedo   = 0.18
	MAX_RAY_DISTANCE = math.MaxFloat64
)

var BACKGROUND_COLOR = NewColor(0, 50, 100)

type hit struct {
	object   Object
	ray      Ray
	distance float64
}

func NewHit(o Object, r Ray, d float64) *hit {
	return &hit{
		object:   o,
		ray:      r,
		distance: d,
	}
}

type Scene struct {
	Objects []Object
	Lights  []Light
	Camera  Camera

	AccelerationStructure AccelerationStructure
}

func NewScene(camera Camera) *Scene {
	return &Scene{
		Objects: []Object{},
		Lights:  []Light{},
		Camera:  camera,
	}
}

func (s *Scene) Add(o ...Object) {
	s.Objects = append(s.Objects, o...)
}

func (s *Scene) AddLights(l ...Light) {
	s.Lights = append(s.Lights, l...)
}

func (s *Scene) Precompute() {
	s.AccelerationStructure = NewBVH(s.Objects, SplitMiddle)
}

func (s *Scene) GetColor(x, y int) Color {
	ray := s.Camera.PixelRay(x, y)
	return s.GetRayColor(ray)
}

func (s *Scene) GetRayColor(ray Ray) Color {
	hit := s.AccelerationStructure.ClosestIntersection(ray, MAX_RAY_DISTANCE)
	if hit == nil {
		return BACKGROUND_COLOR
	}

	color := NewColor(0, 0, 0)
	var objectColor Color
	lightFound := false
	for _, l := range s.Lights {
		point := PointFromRay(hit.ray, hit.distance)
		segment := l.VectorFromPoint(point)
		shadowRay := NewRay(point, segment)

		if pointInShadow(shadowRay, s.AccelerationStructure, segment.Length()) {
			continue
		}
		facingRatio := hit.object.SurfaceNormal(point).Dot(VectorFromTo(point, ray.Origin))
		if facingRatio <= 0 {
			continue
		}

		if !lightFound {
			lightFound = true
			si := &SurfaceInteraction{
				Point:    point,
				Object:   hit.object,
				AS:       s.AccelerationStructure,
				Incident: ray.Direction.Normalize(),
			}
			objectColor = hit.object.GetColor(si)
		}

		lightRatio := hit.object.SurfaceNormal(point).Dot(segment)
		factors := standardAlbedo / math.Pi * l.Intensity(segment.Length()) * facingRatio * lightRatio
		lightColor := l.Color().Times(factors)
		color = color.Add(objectColor.Product(lightColor))
	}
	return color
}

func pointInShadow(shadowRay Ray, as AccelerationStructure, maxDistance float64) bool {
	if hit := as.ClosestIntersection(shadowRay, maxDistance); hit != nil {
		return true
	}
	return false
}
