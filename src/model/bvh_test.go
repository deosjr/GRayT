package model

import (
	"reflect"
	"testing"
)

func TestRecursiveBuild(t *testing.T) {
	for i, tt := range []struct {
		primitives []primitiveInfo
		splitFunc  splitFunc
		wantTotal  int
		wantOrder  []int
		wantTree   bvhNode
	}{
		{
			primitives: []primitiveInfo{
				{
					index:    0,
					bounds:   NewAABB(Vector{-1, -1, -1}, Vector{1, 1, 1}),
					centroid: Vector{0, 0, 0},
				},
			},
			splitFunc: SplitMiddle,
			wantTotal: 1,
			wantOrder: []int{0},
			wantTree: bvhLeaf{
				firstPrimOffset: 0,
				numPrimitives:   1,
				bounds:          NewAABB(Vector{-1, -1, -1}, Vector{1, 1, 1}),
			},
		},
		{
			primitives: []primitiveInfo{
				{
					index:    0,
					bounds:   NewAABB(Vector{-3, 0, -5}, Vector{-1, 2, -3}),
					centroid: Vector{-2, 1, -4},
				},
				{
					index:    1,
					bounds:   NewAABB(Vector{1, -1, -3}, Vector{3, 1, -1}),
					centroid: Vector{2, 0, -2},
				},
			},
			splitFunc: SplitMiddle,
			wantTotal: 3,
			wantOrder: []int{0, 1},
			wantTree: bvhInterior{
				children: [2]bvhNode{
					bvhLeaf{
						firstPrimOffset: 0,
						numPrimitives:   1,
						bounds:          NewAABB(Vector{-3, 0, -5}, Vector{-1, 2, -3}),
					},
					bvhLeaf{
						firstPrimOffset: 1,
						numPrimitives:   1,
						bounds:          NewAABB(Vector{1, -1, -3}, Vector{3, 1, -1}),
					},
				},
				bounds:         NewAABB(Vector{-3, -1, -5}, Vector{3, 2, -1}),
				splitDimension: 0,
			},
		},
	} {
		gotTotal, prims := 0, 0
		gotOrder := make([]int, len(tt.primitives))
		gotTree := recursiveBuildBVH(tt.primitives, 0, len(tt.primitives), &prims, &gotTotal, gotOrder, tt.splitFunc)
		if gotTotal != tt.wantTotal {
			t.Errorf("%d) TOTAL: got %#v want %#v", i, gotTotal, tt.wantTotal)
		}
		if !reflect.DeepEqual(gotOrder, tt.wantOrder) {
			t.Errorf("%d) ORDER: got %#v want %#v", i, gotOrder, tt.wantOrder)
		}
		if !reflect.DeepEqual(gotTree, tt.wantTree) {
			t.Errorf("%d) TREE: got %#v want %#v", i, gotTree, tt.wantTree)
		}
	}
}

func TestFlattenBVHTree(t *testing.T) {
	for i, tt := range []struct {
		tree     bvhNode
		numNodes int
		want     []optimisedBVHNode
	}{
		{
			tree:     bvhLeaf{},
			numNodes: 1,
			want:     []optimisedBVHNode{{}},
		},
		{
			tree: bvhInterior{
				children: [2]bvhNode{
					bvhInterior{
						children: [2]bvhNode{
							bvhLeaf{
								firstPrimOffset: 0,
								numPrimitives:   1,
								bounds:          NewAABB(Vector{1, 1, 1}, Vector{2, 2, 2}),
							}, bvhLeaf{
								firstPrimOffset: 1,
								numPrimitives:   1,
								bounds:          NewAABB(Vector{3, 3, 3}, Vector{4, 4, 4}),
							},
						},
						bounds:         NewAABB(Vector{1, 1, 1}, Vector{4, 4, 4}),
						splitDimension: Y,
					},
					bvhInterior{
						children: [2]bvhNode{
							bvhLeaf{
								firstPrimOffset: 2,
								numPrimitives:   1,
								bounds:          NewAABB(Vector{5, 5, 5}, Vector{6, 6, 6}),
							}, bvhLeaf{
								firstPrimOffset: 3,
								numPrimitives:   1,
								bounds:          NewAABB(Vector{7, 7, 7}, Vector{8, 8, 8}),
							},
						},
						bounds:         NewAABB(Vector{5, 5, 5}, Vector{8, 8, 8}),
						splitDimension: Z,
					},
				},
				bounds:         NewAABB(Vector{1, 1, 1}, Vector{8, 8, 8}),
				splitDimension: Z,
			},
			numNodes: 7,
			want: []optimisedBVHNode{
				{bounds: NewAABB(Vector{1, 1, 1}, Vector{8, 8, 8}), offset: 4, numPrimitives: 0, axis: Z},
				{bounds: NewAABB(Vector{1, 1, 1}, Vector{4, 4, 4}), offset: 3, numPrimitives: 0, axis: Y},
				{bounds: NewAABB(Vector{1, 1, 1}, Vector{2, 2, 2}), offset: 0, numPrimitives: 1},
				{bounds: NewAABB(Vector{3, 3, 3}, Vector{4, 4, 4}), offset: 1, numPrimitives: 1},
				{bounds: NewAABB(Vector{5, 5, 5}, Vector{8, 8, 8}), offset: 6, numPrimitives: 0, axis: Z},
				{bounds: NewAABB(Vector{5, 5, 5}, Vector{6, 6, 6}), offset: 2, numPrimitives: 1},
				{bounds: NewAABB(Vector{7, 7, 7}, Vector{8, 8, 8}), offset: 3, numPrimitives: 1}},
		},
		{
			tree: bvhInterior{
				children: [2]bvhNode{
					bvhLeaf{
						firstPrimOffset: 0,
						numPrimitives:   1,
						bounds:          NewAABB(Vector{-3, 0, -5}, Vector{-1, 2, -3}),
					},
					bvhLeaf{
						firstPrimOffset: 1,
						numPrimitives:   1,
						bounds:          NewAABB(Vector{1, -1, -3}, Vector{3, 1, -1}),
					},
				},
				bounds:         NewAABB(Vector{-3, -1, -5}, Vector{3, 2, -1}),
				splitDimension: X,
			},
			numNodes: 3,
			want: []optimisedBVHNode{
				{bounds: NewAABB(Vector{-3, -1, -5}, Vector{3, 2, -1}), offset: 2, numPrimitives: 0, axis: X},
				{bounds: NewAABB(Vector{-3, 0, -5}, Vector{-1, 2, -3}), offset: 0, numPrimitives: 1},
				{bounds: NewAABB(Vector{1, -1, -3}, Vector{3, 1, -1}), offset: 1, numPrimitives: 1}},
		},
	} {
		got := make([]optimisedBVHNode, tt.numNodes)
		offset := 0
		flattenBVHTree(tt.tree, got, &offset)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %#v want %#v", i, got, tt.want)
		}
	}
}

func TestBVHTraversal(t *testing.T) {
	for i, tt := range []struct {
		bvh  BVH
		ray  Ray
		want hit
	}{
		{
			bvh: BVH{
				objects: []Object{
					Triangle{
						P0: Vector{-1, 0, -1},
						P1: Vector{1, 0, -1},
						P2: Vector{1, 1, -1},
					},
				},
				nodes: []optimisedBVHNode{
					{bounds: AABB{Pmin: Vector{X: -1, Y: 0, Z: -1}, Pmax: Vector{X: 1, Y: 1, Z: -1}}, offset: 0, numPrimitives: 1},
				},
			},
			ray: NewRay(Vector{0, 0, 0}, Vector{0, 0, -1}),
			want: hit{
				object: Triangle{
					P0: Vector{-1, 0, -1},
					P1: Vector{1, 0, -1},
					P2: Vector{1, 1, -1},
				},
				point: Vector{0, 0, -1},
			},
		},
		{
			bvh: BVH{
				objects: []Object{
					Sphere{
						Center: Vector{-2, 1, -4},
						Radius: 1.0,
					},
					Sphere{
						Center: Vector{2, 0, -2},
						Radius: 1.0,
					},
				},
				nodes: []optimisedBVHNode{
					{bounds: NewAABB(Vector{-3, -1, -5}, Vector{3, 2, -1}), offset: 2, numPrimitives: 0, axis: X},
					{bounds: NewAABB(Vector{-3, 0, -5}, Vector{-1, 2, -3}), offset: 0, numPrimitives: 1},
					{bounds: NewAABB(Vector{1, -1, -3}, Vector{3, 1, -1}), offset: 1, numPrimitives: 1}},
			},
			ray: NewRay(Vector{0, 0, 0}, Vector{1, 0, -1}),
			want: hit{
				object: Sphere{
					Center: Vector{2, 0, -2},
					Radius: 1.0,
				},
				point: Vector{1.2928932188134528, 0, -1.2928932188134528},
			},
		},
	} {
		hit := tt.bvh.ClosestIntersection(tt.ray)
		if hit == nil {
			t.Errorf("%d) got nil want %#v", i, tt.want)
			continue
		}
		got := *hit
		if got != tt.want {
			t.Errorf("%d) got %#v want %#v", i, got, tt.want)
		}
	}
}