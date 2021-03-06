package model

import "testing"

func TestTriangleIntersect(t *testing.T) {
	for i, tt := range []struct {
		t         Triangle
		r         Ray
		want      float32
		wantTruth bool
	}{
		{
			t: Triangle{
				P0: Vector{-1, 0, -1},
				P1: Vector{1, 0, -1},
				P2: Vector{1, 1, -1},
			},
			r: Ray{
				Origin:    Vector{0, 0, 0},
				Direction: Vector{0, 0, -1},
			},
			want:      1.0,
			wantTruth: true,
		},
	} {
		got, found := tt.t.Intersect(tt.r)
		if !found && tt.wantTruth == false {
			continue
		}
		if (!found && tt.wantTruth == true) || (found && tt.wantTruth == false) {
			t.Errorf("%d) incorrect bool value; want %v", i, tt.wantTruth)
			continue
		}
		if got.distance != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestTriangleSurfaceNormal(t *testing.T) {
	for i, tt := range []struct {
		t    Triangle
		v    Vector
		want Vector
	}{
		{
			t: Triangle{
				P0: Vector{-1, 0, 0},
				P1: Vector{1, 0, 0},
				P2: Vector{1, 1, 0},
			},
			v:    Vector{1, 1, 1},
			want: Vector{0, 0, 1},
		},
	} {
		got := tt.t.SurfaceNormal(tt.v)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestTriangleSurfaceArea(t *testing.T) {
	for i, tt := range []struct {
		t    Triangle
		want float32
	}{
		{
			t: Triangle{
				P0: Vector{-1, 0, 0},
				P1: Vector{1, 0, 0},
				P2: Vector{1, 1, 0},
			},
			want: 1,
		},
		{
			t: Triangle{
				P0: Vector{3, 0, 2},
				P1: Vector{1, 0, 0},
				P2: Vector{3, 0, 8},
			},
			want: 6,
		},
	} {
		got := tt.t.SurfaceArea()
		if !compareFloat32(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}
