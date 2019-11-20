package model

import (
	"math"
	"reflect"
	"testing"
)

func TestRecursiveBuild(t *testing.T) {
	for i, tt := range []struct {
		objects   []objectInfo
		splitFunc splitFunc
		wantTotal int
		wantOrder []int
		wantTree  bvhNode
	}{
		{
			objects: []objectInfo{
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
				firstOffset: 0,
				numObjects:  1,
				bounds:      NewAABB(Vector{-1, -1, -1}, Vector{1, 1, 1}),
			},
		},
		{
			objects: []objectInfo{
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
						firstOffset: 0,
						numObjects:  1,
						bounds:      NewAABB(Vector{-3, 0, -5}, Vector{-1, 2, -3}),
					},
					bvhLeaf{
						firstOffset: 1,
						numObjects:  1,
						bounds:      NewAABB(Vector{1, -1, -3}, Vector{3, 1, -1}),
					},
				},
				bounds:         NewAABB(Vector{-3, -1, -5}, Vector{3, 2, -1}),
				splitDimension: 0,
			},
		},
	} {
		gotTotal, prims := 0, 0
		gotOrder := make([]int, len(tt.objects))
		gotTree := recursiveBuildBVH(tt.objects, 0, len(tt.objects), &prims, &gotTotal, gotOrder, tt.splitFunc)
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
								firstOffset: 0,
								numObjects:  1,
								bounds:      NewAABB(Vector{1, 1, 1}, Vector{2, 2, 2}),
							}, bvhLeaf{
								firstOffset: 1,
								numObjects:  1,
								bounds:      NewAABB(Vector{3, 3, 3}, Vector{4, 4, 4}),
							},
						},
						bounds:         NewAABB(Vector{1, 1, 1}, Vector{4, 4, 4}),
						splitDimension: Y,
					},
					bvhInterior{
						children: [2]bvhNode{
							bvhLeaf{
								firstOffset: 2,
								numObjects:  1,
								bounds:      NewAABB(Vector{5, 5, 5}, Vector{6, 6, 6}),
							}, bvhLeaf{
								firstOffset: 3,
								numObjects:  1,
								bounds:      NewAABB(Vector{7, 7, 7}, Vector{8, 8, 8}),
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
				{bounds: NewAABB(Vector{1, 1, 1}, Vector{8, 8, 8}), offset: 4, numObjects: 0, axis: Z},
				{bounds: NewAABB(Vector{1, 1, 1}, Vector{4, 4, 4}), offset: 3, numObjects: 0, axis: Y},
				{bounds: NewAABB(Vector{1, 1, 1}, Vector{2, 2, 2}), offset: 0, numObjects: 1},
				{bounds: NewAABB(Vector{3, 3, 3}, Vector{4, 4, 4}), offset: 1, numObjects: 1},
				{bounds: NewAABB(Vector{5, 5, 5}, Vector{8, 8, 8}), offset: 6, numObjects: 0, axis: Z},
				{bounds: NewAABB(Vector{5, 5, 5}, Vector{6, 6, 6}), offset: 2, numObjects: 1},
				{bounds: NewAABB(Vector{7, 7, 7}, Vector{8, 8, 8}), offset: 3, numObjects: 1}},
		},
		{
			tree: bvhInterior{
				children: [2]bvhNode{
					bvhLeaf{
						firstOffset: 0,
						numObjects:  1,
						bounds:      NewAABB(Vector{-3, 0, -5}, Vector{-1, 2, -3}),
					},
					bvhLeaf{
						firstOffset: 1,
						numObjects:  1,
						bounds:      NewAABB(Vector{1, -1, -3}, Vector{3, 1, -1}),
					},
				},
				bounds:         NewAABB(Vector{-3, -1, -5}, Vector{3, 2, -1}),
				splitDimension: X,
			},
			numNodes: 3,
			want: []optimisedBVHNode{
				{bounds: NewAABB(Vector{-3, -1, -5}, Vector{3, 2, -1}), offset: 2, numObjects: 0, axis: X},
				{bounds: NewAABB(Vector{-3, 0, -5}, Vector{-1, 2, -3}), offset: 0, numObjects: 1},
				{bounds: NewAABB(Vector{1, -1, -3}, Vector{3, 1, -1}), offset: 1, numObjects: 1}},
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
		bvh             BVH
		ray             Ray
		want            SurfaceInteraction
		wantObjectIndex int
	}{
		{
			bvh: BVH{
				objects: []Object{
					Triangle{
						P0: Vector{-1, 0, 1},
						P1: Vector{1, 0, 1},
						P2: Vector{1, 1, 1},
					},
				},
				nodes: []optimisedBVHNode{
					{bounds: AABB{Pmin: Vector{X: -1, Y: 0, Z: 1}, Pmax: Vector{X: 1, Y: 1, Z: 1}}, offset: 0, numObjects: 1},
				},
			},
			ray: NewRay(Vector{0, 0, 0}, Vector{0, 0, 1}),
			want: SurfaceInteraction{
				distance: 1,
				// normal:   Vector{0, 0, -1},
			},
			wantObjectIndex: 0,
		},
		{
			bvh: BVH{
				objects: []Object{
					Sphere{
						Center: Vector{-2, 1, 4},
						Radius: 1.0,
					},
					Sphere{
						Center: Vector{2, 0, 2},
						Radius: 1.0,
					},
				},
				nodes: []optimisedBVHNode{
					{bounds: NewAABB(Vector{-3, -1, 5}, Vector{3, 2, 1}), offset: 2, numObjects: 0, axis: X},
					{bounds: NewAABB(Vector{-3, 0, 5}, Vector{-1, 2, 3}), offset: 0, numObjects: 1},
					{bounds: NewAABB(Vector{1, -1, 3}, Vector{3, 1, 1}), offset: 1, numObjects: 1}},
			},
			ray: NewRay(Vector{0, 0, 0}, Vector{1, 0, 1}),
			want: SurfaceInteraction{
				distance: 1.8284271247461907,
				// normal:   Vector{1.2928932188134528, 0, -1.2928932188134528},
			},
			wantObjectIndex: 1,
		},
	} {
		tt.want.object = tt.bvh.objects[tt.wantObjectIndex]
		got, ok := tt.bvh.ClosestIntersection(tt.ray, math.MaxFloat32)
		if !ok {
			t.Errorf("%d) got nil want %#v", i, tt.want)
			continue
		}
		if !reflect.DeepEqual(tt.want.object, got.object) {
			t.Errorf("%d) got %#v want %#v", i, got, tt.want)
		}
		if !compareFloat32(got.distance, tt.want.distance) {
			t.Errorf("%d) got %#v want %#v", i, got, tt.want)
		}
	}
}

func TestTriangle4BVH(t *testing.T) {
	for i, tt := range []struct {
		triangles []Triangle
		want      []bvh4Node
	}{
		{
			triangles: []Triangle{
				{
					P0: Vector{-1, -1, -1},
					P1: Vector{1, 1, 1},
					P2: Vector{-1, 1, 1},
				},
			},
			want: []bvh4Node{
				bvh4Leaf{
					firstOffset:  0,
					numTriangles: 1,
					p0x:          [4]float32{-1, 0, 0, 0},
					p0y:          [4]float32{-1, 0, 0, 0},
					p0z:          [4]float32{-1, 0, 0, 0},
					p1x:          [4]float32{1, 0, 0, 0},
					p1y:          [4]float32{1, 0, 0, 0},
					p1z:          [4]float32{1, 0, 0, 0},
					p2x:          [4]float32{-1, 0, 0, 0},
					p2y:          [4]float32{1, 0, 0, 0},
					p2z:          [4]float32{1, 0, 0, 0},
				},
			},
		},
		{
			triangles: []Triangle{
				{
					P0: Vector{-1, -1, -1},
					P1: Vector{1, 1, 1},
					P2: Vector{-1, 1, 1},
				},
				{
					P0: Vector{-2, -2, -2},
					P1: Vector{2, 2, 2},
					P2: Vector{-2, 2, 2},
				},
				{
					P0: Vector{-3, -3, -3},
					P1: Vector{3, 3, 3},
					P2: Vector{-3, 3, 3},
				},
				{
					P0: Vector{-3, 0, -5},
					P1: Vector{-1, 2, -3},
					P2: Vector{-2, 1, -4},
				},
				{
					P0: Vector{1, -1, -3},
					P1: Vector{3, 1, -1},
					P2: Vector{2, 0, -2},
				},
			},
			want: []bvh4Node{
				bvh4Interior{
					childOffsets: [4]int{1, -1, 2, -1},
					min4x:        [4]float32{-3, 0, -3, 0},
					min4y:        [4]float32{0, 0, -3, 0},
					min4z:        [4]float32{-5, 0, -3, 0},
					max4x:        [4]float32{-1, 0, 3, 0},
					max4y:        [4]float32{2, 0, 3, 0},
					max4z:        [4]float32{-3, 0, 3, 0},
				},
				bvh4Leaf{
					firstOffset:  0,
					numTriangles: 1,
					p0x:          [4]float32{-3, 0, 0, 0},
					p0y:          [4]float32{0, 0, 0, 0},
					p0z:          [4]float32{-5, 0, 0, 0},
					p1x:          [4]float32{-1, 0, 0, 0},
					p1y:          [4]float32{2, 0, 0, 0},
					p1z:          [4]float32{-3, 0, 0, 0},
					p2x:          [4]float32{-2, 0, 0, 0},
					p2y:          [4]float32{1, 0, 0, 0},
					p2z:          [4]float32{-4, 0, 0, 0},
				},
				bvh4Leaf{
					firstOffset:  1,
					numTriangles: 4,
					p0x:          [4]float32{-2, -3, -1, 1},
					p0y:          [4]float32{-2, -3, -1, -1},
					p0z:          [4]float32{-2, -3, -1, -3},
					p1x:          [4]float32{2, 3, 1, 3},
					p1y:          [4]float32{2, 3, 1, 1},
					p1z:          [4]float32{2, 3, 1, -1},
					p2x:          [4]float32{-2, -3, -1, 2},
					p2y:          [4]float32{2, 3, 1, 0},
					p2z:          [4]float32{2, 3, 1, -2},
				},
			},
		},
	} {
		got := NewTriangle4BVH(tt.triangles)
		if !reflect.DeepEqual(got.nodes, tt.want) {
			t.Errorf("%d) TREE: got %#v want %#v", i, got.nodes, tt.want)
		}
	}
}

func TestTriangle4BVHTraversal(t *testing.T) {
	for i, tt := range []struct {
		triangles       []Triangle
		ray             Ray
		want            SurfaceInteraction
		wantObjectIndex int
	}{
		{
			triangles: []Triangle{
				{
					P0: Vector{-1, 0, 1},
					P1: Vector{1, 0, 1},
					P2: Vector{1, 1, 1},
				},
			},
			ray: NewRay(Vector{0, 0, 0}, Vector{0, 0, 1}),
			want: SurfaceInteraction{
				distance: 1,
				// normal:   Vector{0, 0, -1},
			},
			wantObjectIndex: 0,
		},
		{
			triangles: []Triangle{
				{
					P0: Vector{-1, 0, 1},
					P1: Vector{1, 0, 1},
					P2: Vector{1, 1, 1},
				},
				{
					P0: Vector{-1, -1, -1},
					P1: Vector{1, 1, 1},
					P2: Vector{-1, 1, 1},
				},
				{
					P0: Vector{-2, -2, -2},
					P1: Vector{2, 2, 2},
					P2: Vector{-2, 2, 2},
				},
				{
					P0: Vector{-3, -3, -3},
					P1: Vector{3, 3, 3},
					P2: Vector{-3, 3, 3},
				},
				{
					P0: Vector{-3, 0, -5},
					P1: Vector{-1, 2, -3},
					P2: Vector{-2, 1, -4},
				},
				{
					P0: Vector{1, -1, -3},
					P1: Vector{3, 1, -1},
					P2: Vector{2, 0, -2},
				},
			},
			ray: NewRay(Vector{0, 0, 0}, Vector{0, 0, 1}),
			want: SurfaceInteraction{
				distance: 1,
				// normal:   Vector{0, 0, -1},
			},
			wantObjectIndex: 0,
		},
	} {
		bvh := NewTriangle4BVH(tt.triangles)
		tt.want.object = bvh.triangles[tt.wantObjectIndex]
		got, ok := bvh.ClosestIntersection(tt.ray, math.MaxFloat32)
		if !ok {
			t.Errorf("%d) got nil want %#v", i, tt.want)
			continue
		}
		/*
			if !reflect.DeepEqual(tt.want.object, got.object) {
				t.Errorf("%d) got %#v want %#v", i, got, tt.want)
			}
		*/
		if !compareFloat32(got.distance, tt.want.distance) {
			t.Errorf("%d) got %#v want %#v", i, got, tt.want)
		}
	}
}
