package model

import (
	"math"
	"math/rand"
	"time"
)

const (
	standardAlbedo   = 0.18
	MAX_RAY_DISTANCE = 1000000.0
	MAX_RAY_DEPTH    = 5
	INVPI            = 1 / math.Pi
)

var BACKGROUND_COLOR Color
var BLACK = NewColor(0, 0, 0)

type Tracer interface {
	GetRayColor(Ray, *Scene, int) Color
	Random() *rand.Rand
}

type tracer struct {
	random *rand.Rand
}

func (t tracer) Random() *rand.Rand {
	return t.random
}

type TracerType uint

const (
	WhittedStyle TracerType = iota
	Path
	PathNextEventEstimate
)

type whittedRayTracer struct {
	tracer
}

func NewWhittedRayTracer() Tracer {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &whittedRayTracer{tracer{random: r}}
}

func (wrt whittedRayTracer) GetRayColor(ray Ray, scene *Scene, depth int) Color {
	if depth == MAX_RAY_DEPTH {
		return BLACK
	}

	si, ok := scene.AccelerationStructure.ClosestIntersection(ray, MAX_RAY_DISTANCE)
	if !ok {
		return BACKGROUND_COLOR
	}

	si.as = scene.AccelerationStructure
	si.depth = depth
	si.tracer = wrt

	color := NewColor(0, 0, 0)
	material := si.object.GetMaterial()
	for _, light := range scene.Lights {
	    lightSegment := light.GetLightSegment(si.Point)
	    maxDistance := lightSegment.Length()
		if pointInShadow(si.Point, lightSegment, maxDistance, si.as) {
			continue
		}
		facingRatio := si.normal.Dot(si.incident.Times(-1))
		if facingRatio <= 0 {
			continue
		}

		var objectColor Color
		switch mat := material.(type) {
		case *RadiantMaterial:
			objectColor = mat.Color
        case *NormalMappingMaterial:
            // NOTE: normal mapping only wraps diffuse now
            si.normal = mat.NormalFunc(si)
			lightRatio := si.normal.Dot(lightSegment.Normalize())
			factors := standardAlbedo * INVPI * light.Intensity(lightSegment.Length()) * lightRatio
			lightColor := light.Color().Times(factors)
			objectColor = mat.WrappedMaterial.(*DiffuseMaterial).Color.Product(lightColor)
		case *DiffuseMaterial:
			lightRatio := si.normal.Dot(lightSegment.Normalize())
			factors := standardAlbedo * INVPI * light.Intensity(lightSegment.Length()) * lightRatio
			lightColor := light.Color().Times(factors)
			objectColor = mat.Color.Product(lightColor)
		case *ReflectiveMaterial:
			i := si.incident
			n := si.object.SurfaceNormal(si.Point)
			reflection := i.Sub(n.Times(2 * i.Dot(n)))
			newRay := NewRay(si.Point, reflection)
			// TODO: retain maxdistance for tracing
			objectColor = wrt.GetRayColor(newRay, scene, depth+1) //.Times(1 - standardAlbedo) // simulates nonperfect reflection
		case *PosFuncMat:
			objectColor = mat.GetColor(si)
		}
		color = color.Add(objectColor.Times(facingRatio))
	}
	return color
}

func pointInShadow(point, segment Vector, maxDistance float32, as AccelerationStructure) bool {
	shadowRay := NewRay(point, segment)
	if _, ok := as.ClosestIntersection(shadowRay, maxDistance); ok {
		return true
	}
	return false
}

type pathTracer struct {
	tracer
}

func NewPathTracer() Tracer {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &pathTracer{tracer{random: r}}
}

const pdf = 2.0 * math.Pi

func (pt *pathTracer) GetRayColor(ray Ray, scene *Scene, depth int) Color {
	if depth == MAX_RAY_DEPTH {
		return BLACK
	}

	si, ok := scene.AccelerationStructure.ClosestIntersection(ray, MAX_RAY_DISTANCE)
	if !ok {
		return BLACK
	}

	si.as = scene.AccelerationStructure
	si.depth = depth
	si.tracer = pt

	if si.object.IsLight() {
		return si.object.GetMaterial().(*RadiantMaterial).Color
	}
	surfaceDiffuseColor := si.object.GetColor(si)
	brdf := surfaceDiffuseColor.Times(INVPI)

	// random new ray
	randomDirection := randomInHemisphere(pt.random, si.normal)
	newRay := NewRay(si.Point, randomDirection)
	cos := si.normal.Dot(randomDirection)
	recursiveColor := pt.GetRayColor(newRay, scene, depth+1)
	//pdf := 1.0 / (2.0 * math.Pi)
	//TODO: albedo? just multiplying with standardalbedo leads to horrible results
	// probably because light intensity does not make sense yet
	sampleColor := recursiveColor.Times(cos * pdf).Product(brdf)
	return sampleColor
}

type pathTracerNEE struct {
	tracer
}

func NewPathTracerNEE() Tracer {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &pathTracerNEE{tracer{random: r}}
}

func (pt *pathTracerNEE) GetRayColor(ray Ray, scene *Scene, depth int) Color {
	if depth == MAX_RAY_DEPTH {
		return BLACK
	}

	si, ok := scene.AccelerationStructure.ClosestIntersection(ray, MAX_RAY_DISTANCE)
	if !ok {
		return BLACK
	}

	si.as = scene.AccelerationStructure
	si.depth = depth
	si.tracer = pt

	// only count light if we immediately hit it
	// direct light sampling counts the rest
	if si.object.IsLight() {
		if depth == 0 {
			return si.object.GetMaterial().(*RadiantMaterial).Color
		}
		return BLACK
	}
	surfaceDiffuseColor := si.object.GetColor(si)
	brdf := surfaceDiffuseColor.Times(INVPI)

	// direct light sampling
	direct := NewColor(0, 0, 0)
	light := scene.randomEmitter(pt.random)
	lpoint := light.Sample(pt.random)
	nl := light.SurfaceNormal(lpoint)
	l := VectorFromTo(si.Point, lpoint)
	lightFacing := si.normal.Dot(l.Normalize())
	dist := l.Length()
	lightCos := nl.Dot(l.Normalize().Times(-1))
	if lightFacing > 0 && lightCos > 0 && !pointInShadow(si.Point, l, dist, si.as) {
		lightPDF := 1.0 / float32(len(scene.Emitters))
		solidAngle := (lightCos * light.SurfaceArea()) / (dist * dist * lightPDF)
		lightColor := light.GetMaterial().(*RadiantMaterial).Color
		direct = lightColor.Times(solidAngle).Product(brdf).Times(lightFacing)
	}

	// indirect light sampling: random new ray
	randomDirection := si.object.GetMaterial().Sample(pt.random, si.normal)
	newRay := NewRay(si.Point, randomDirection)
	cos := si.normal.Dot(randomDirection)
	recursiveColor := pt.GetRayColor(newRay, scene, depth+1)
	//pdf := 1.0 / (2.0 * math.Pi)
	//TODO: albedo? just multiplying with standardalbedo leads to horrible results
	// probably because light intensity does not make sense yet
	indirect := recursiveColor.Times(cos * pdf).Product(brdf)
	return direct.Add(indirect)
}
