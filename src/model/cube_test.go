package model

import "testing"

func TestNewAABB(t *testing.T) {
	for i, tt := range []struct {
		p1   Vector
		p2   Vector
		want AABB
	}{
		{
			p1: Vector{0, 0, 0},
			p2: Vector{1, 1, 1},
			want: AABB{
				Pmin: Vector{0, 0, 0},
				Pmax: Vector{1, 1, 1},
			},
		},
		{
			p1: Vector{235, 234, 678},
			p2: Vector{567, 123, 567},
			want: AABB{
				Pmin: Vector{235, 123, 567},
				Pmax: Vector{567, 234, 678},
			},
		},
	} {
		got := NewAABB(tt.p1, tt.p2)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestAABBAddPoint(t *testing.T) {
	for i, tt := range []struct {
		aabb AABB
		p    Vector
		want AABB
	}{
		{
			aabb: AABB{
				Pmin: Vector{0, 0, 0},
				Pmax: Vector{1, 1, 1},
			},
			p: Vector{2, 2, 2},
			want: AABB{
				Pmin: Vector{0, 0, 0},
				Pmax: Vector{2, 2, 2},
			},
		},
		{
			aabb: AABB{
				Pmin: Vector{235, 123, 567},
				Pmax: Vector{567, 234, 678},
			},
			p: Vector{2, 999, 589},
			want: AABB{
				Pmin: Vector{2, 123, 567},
				Pmax: Vector{567, 999, 678},
			},
		},
	} {
		got := tt.aabb.AddPoint(tt.p)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestAABBAddBox(t *testing.T) {
	for i, tt := range []struct {
		b1   AABB
		b2   AABB
		want AABB
	}{
		{
			b1: AABB{
				Pmin: Vector{0, 0, 0},
				Pmax: Vector{1, 1, 1},
			},
			b2: AABB{
				Pmin: Vector{2, 2, 2},
				Pmax: Vector{3, 3, 3},
			},
			want: AABB{
				Pmin: Vector{0, 0, 0},
				Pmax: Vector{3, 3, 3},
			},
		},
		{
			b1: AABB{
				Pmin: Vector{235, 123, 567},
				Pmax: Vector{567, 234, 678},
			},
			b2: AABB{
				Pmin: Vector{333, 155, 444},
				Pmax: Vector{999, 999, 999},
			},
			want: AABB{
				Pmin: Vector{235, 123, 444},
				Pmax: Vector{999, 999, 999},
			},
		},
	} {
		got := tt.b1.AddAABB(tt.b2)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestAABBRayIntersect(t *testing.T) {
	for i, tt := range []struct {
		aabb AABB
		ray  Ray
		want bool
	}{
		{
			aabb: AABB{
				Pmin: Vector{2, 2, 2},
				Pmax: Vector{3, 3, 3},
			},
			ray:  NewRay(Vector{0, 0, 0}, Vector{1, 1, 1}),
			want: true,
		},
		{
			aabb: AABB{
				Pmin: Vector{2, 2, 2},
				Pmax: Vector{3, 3, 3},
			},
			ray:  NewRay(Vector{0, 0, 0}, Vector{1, 0, 0}),
			want: false,
		},
	} {
		got := tt.aabb.Intersect(tt.ray)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}
