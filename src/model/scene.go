package model

import "math"

const (
	standardAlbedo   = 0.18
	MAX_RAY_DISTANCE = math.MaxFloat64
)

var BACKGROUND_COLOR = NewColor(0, 50, 100)

type hit struct {
	object   Object
	normal   Vector
	distance float64
}

func NewHit(o Object, d float64) *hit {
	return &hit{
		object:   o,
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
		point := PointFromRay(ray, hit.distance)
		lightSegment := l.GetLightSegment(point)

		if pointInShadow(point, lightSegment, s.AccelerationStructure) {
			continue
		}
		facingRatio := hit.normal.Dot(VectorFromTo(point, ray.Origin))
		if facingRatio <= 0 {
			continue
		}

		if !lightFound {
			lightFound = true
			si := &SurfaceInteraction{
				Point:    point,
				Normal:   hit.normal,
				Object:   hit.object,
				AS:       s.AccelerationStructure,
				// already normalized
				Incident: ray.Direction,
			}
			objectColor = hit.object.GetColor(si)
		}

		lightRatio := l.LightRatio(point, hit.normal)
		// TODO: current lighting weirdness issue is at least
		// in part due to this formula only applying to diffuse 
		// surfaces. At the moment it's also applied to reflective ones!
		factors := standardAlbedo / math.Pi * l.Intensity(lightSegment.Length()) * facingRatio * lightRatio
		lightColor := l.Color().Times(factors)
		color = color.Add(objectColor.Product(lightColor))
	}
	return color
}

func pointInShadow(point Vector, lightSegment Vector, as AccelerationStructure) bool {
	shadowRay := NewRay(point, lightSegment)
	maxDistance := lightSegment.Length()
	if hit := as.ClosestIntersection(shadowRay, maxDistance); hit != nil {
		return true
	}
	return false
}
