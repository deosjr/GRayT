package model

import (
	"math"
	"testing"
)

func TestObjectIntersect(t *testing.T) {
	for i, tt := range []struct {
		o         Object
		r         Ray
		want      float64
		wantTruth bool
	}{
		{
			o: NewComplexObject([]Object{
				Triangle{
					P0: Vector{-1, 0, -1},
					P1: Vector{1, 0, -1},
					P2: Vector{1, 1, -1},
				},
				Triangle{
					P0: Vector{1, 0, 1},
					P1: Vector{1, 0, -1},
					P2: Vector{1, 1, -1},
				}}),
			r: Ray{
				Origin:    Vector{0, 0, 0},
				Direction: Vector{0, 0, -1},
			},
			want:      1.0,
			wantTruth: true,
		},
		{
			o: NewSharedObject(
				Sphere{
					Center: Vector{0, 0, 0},
					Radius: 0.5,
				}, Translate(Vector{0, 0, -1})),
			r: Ray{
				Origin:    Vector{0, 0, 0},
				Direction: Vector{0, 0, -1},
			},
			want:      0.5,
			wantTruth: true,
		},
		{
			o: NewSharedObject(
				Sphere{
					Center: Vector{0, 0, 0},
					Radius: 0.5,
				}, Translate(Vector{0, 0, 2}).Mul(RotateY(math.Pi/2))),
			r: Ray{
				Origin:    Vector{0, 0, 0},
				Direction: Vector{0, 0, 1},
			},
			want:      1.5,
			wantTruth: true,
		},
	} {
		hit, found := tt.o.Intersect(tt.r)
		if !found && tt.wantTruth == false {
			continue
		}
		if (!found && tt.wantTruth == true) || (found && tt.wantTruth == false) {
			t.Errorf("%d) incorrect bool value; want %v", i, tt.wantTruth)
			continue
		}
		got := hit.distance
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestObjectBound(t *testing.T) {
	for i, tt := range []struct {
		o    Object
		t    Transform
		want AABB
	}{
		{
			o: NewComplexObject([]Object{
				Triangle{
					P0: Vector{-1, 0, -1},
					P1: Vector{1, 0, -1},
					P2: Vector{1, 1, -1},
				},
				Triangle{
					P0: Vector{1, 0, 1},
					P1: Vector{1, 0, -1},
					P2: Vector{1, 1, -1},
				}}),
			t:    identity,
			want: NewAABB(Vector{-1, 0, -1}, Vector{1, 1, 1}),
		},
		{
			o: NewSharedObject(
				Sphere{
					Center: Vector{0, 0, 0},
					Radius: 1,
				}, Translate(Vector{0, 0, 0})),
			t:    identity,
			want: NewAABB(Vector{-1, -1, -1}, Vector{1, 1, 1}),
		},
		{
			o: NewSharedObject(
				Sphere{
					Center: Vector{0, 0, 0},
					Radius: 1,
				}, Translate(Vector{0, 0, 1})),
			t:    identity,
			want: NewAABB(Vector{-1, -1, 0}, Vector{1, 1, 2}),
		},
		{
			o: NewSharedObject(
				Sphere{
					Center: Vector{0, 0, 0},
					Radius: 1,
				}, RotateY(math.Pi/2)),
			t:    identity,
			want: NewAABB(Vector{-1, -1, -1}, Vector{1, 1, 1}),
		},
		{
			o: NewSharedObject(
				Sphere{
					Center: Vector{0, 0, 0},
					Radius: 1,
				}, Translate(Vector{2, 2, 2}).Mul(RotateY(math.Pi/2))),
			t:    identity,
			want: NewAABB(Vector{1, 1, 1}, Vector{3, 3, 3}),
		},
	} {
		got := tt.o.Bound(tt.t)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}
