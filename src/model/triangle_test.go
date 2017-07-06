package model

import "testing"

func TestTriangleIntersect(t *testing.T) {
	for i, tt := range []struct {
		t         Triangle
		r         Ray
		want      float64
		wantTruth bool
	}{
		{
			t: Triangle{
				P0: Vector{1, 0, 0},
				P1: Vector{0, 1, 0},
				P2: Vector{0, 0, 1},
			},
			r: Ray{
				Origin:    Vector{0, 0, 0},
				Direction: Vector{0, 0, -1},
			},
			want:      0.0,
			wantTruth: false,
		},
	} {
		got, gotTruth := tt.t.Intersect(tt.r)
		if gotTruth != tt.wantTruth {
			t.Errorf("%d) incorrect bool value; got %v want %v", i, gotTruth, tt.wantTruth)
			continue
		}
		if got != tt.want {
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
				P0: Vector{1, 0, 0},
				P1: Vector{0, 1, 0},
				P2: Vector{0, 0, 1},
			},
			v:    Vector{1, 1, 1},
			want: Vector{0.5773502691896258, 0.5773502691896258, 0.5773502691896258},
		},
	} {
		got := tt.t.SurfaceNormal(tt.v)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}
