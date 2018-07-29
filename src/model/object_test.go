package model

import "testing"

func TestObjectIntersect(t *testing.T) {
	for i, tt := range []struct {
		co        Object
		r         Ray
		want      float64
		wantTruth bool
	}{
		{
			co: NewComplexObject([]Object{
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
	} {
		hit := tt.co.Intersect(tt.r)
		if hit == nil && tt.wantTruth == false {
			continue
		}
		if (hit == nil && tt.wantTruth == true) || (hit != nil && tt.wantTruth == false) {
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
		co   Object
		want AABB
	}{
		{
			co: NewComplexObject([]Object{
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
			want: NewAABB(Vector{-1, 0, -1}, Vector{1, 1, 1}),
		},
	} {
		got := tt.co.Bound()
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}
