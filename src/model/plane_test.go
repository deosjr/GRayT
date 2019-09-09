package model

import "testing"

func TestPlaneIntersect(t *testing.T) {
	for i, tt := range []struct {
		p         Plane
		r         Ray
		want      float64
		wantTruth bool
	}{
		{
			p: Plane{
				Point:  Vector{0, 0, 0},
				Normal: Vector{1, 1, 1}.Normalize(),
			},
			r: Ray{
				Origin:    Vector{0, 0, 0},
				Direction: Vector{0, 0, 0},
			},
			want:      0.0,
			wantTruth: false,
		},
		{
			p: Plane{
				Point:  Vector{1, 1, 5},
				Normal: Vector{1, 1, 1}.Normalize(),
			},
			r: Ray{
				Origin:    Vector{0, 0, 0},
				Direction: Vector{6, 1, 4}.Normalize(),
			},
			want:      4.63279720226942,
			wantTruth: true,
		},
	} {
		got, found := tt.p.Intersect(tt.r)
		if !found && tt.wantTruth == false {
			continue
		}
		if (!found && tt.wantTruth == true) || (found && tt.wantTruth == false) {
			t.Errorf("%d) incorrect bool value; want %v", i, tt.wantTruth)
			continue
		}
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestPlaneSurfaceNormal(t *testing.T) {
	for i, tt := range []struct {
		p    Plane
		v    Vector
		want Vector
	}{
		{
			p: Plane{
				Point:  Vector{0, 0, 0},
				Normal: Vector{15, -3, 1.42}.Normalize(),
			},
			v:    Vector{1, 1, 1},
			want: Vector{15, -3, 1.42}.Normalize(),
		},
		{
			p: Plane{
				Point:  Vector{0, 0, 0},
				Normal: Vector{1, 1, 1}.Normalize(),
			},
			v:    Vector{15.3, 3.14, 6},
			want: Vector{1, 1, 1}.Normalize(),
		},
	} {
		got := tt.p.SurfaceNormal(tt.v)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}
