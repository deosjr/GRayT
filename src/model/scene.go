package model

const (
	standardAlbedo   = 0.18
	MAX_RAY_DISTANCE = 1000000.0
)

var BACKGROUND_COLOR Color

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
	for _, light := range s.Lights {
		point := PointFromRay(ray, hit.distance)
		if pointInShadow(light, point, s.AccelerationStructure) {
			continue
		}
		facingRatio := hit.normal.Dot(ray.Direction.Times(-1))
		if facingRatio <= 0 {
			continue
		}

		si := &SurfaceInteraction{
			Point:  point,
			Normal: hit.normal,
			Object: hit.object,
			AS:     s.AccelerationStructure,
			// already normalized
			Incident: ray.Direction,
		}
		objectColor := hit.object.GetColor(si, light)
		color = color.Add(objectColor)
	}
	return color
}

func pointInShadow(light Light, point Vector, as AccelerationStructure) bool {
	lightSegment := light.GetLightSegment(point)
	shadowRay := NewRay(point, lightSegment)
	maxDistance := lightSegment.Length()
	if hit := as.ClosestIntersection(shadowRay, maxDistance); hit != nil {
		return true
	}
	return false
}

func SetBackgroundColor(c Color) {
	BACKGROUND_COLOR = c
}
