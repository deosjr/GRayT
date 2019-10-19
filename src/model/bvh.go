package model

import (
	"math"
	"sort"

	"github.com/deosjr/GRayT/src/simd"
)

// Bounded Volume Hierarchy
type BVH struct {
	objects []Object
	nodes   []optimisedBVHNode
}

func (bvh *BVH) GetObjects() []Object {
	return bvh.objects
}

type bvhNode interface {
	getBounds() AABB
	getNumObjects() int
}

type bvhLeaf struct {
	firstOffset int
	numObjects  int
	bounds      AABB
}

func newBVHLeaf(first, n int, bounds AABB) bvhLeaf {
	return bvhLeaf{
		firstOffset: first,
		numObjects:  n,
		bounds:      bounds,
	}
}

func (n bvhLeaf) getBounds() AABB {
	return n.bounds
}

func (n bvhLeaf) getNumObjects() int {
	return n.numObjects
}

type bvhInterior struct {
	children       [2]bvhNode
	bounds         AABB
	splitDimension Dimension
}

func newBVHInterior(dim Dimension, c1, c2 bvhNode) bvhInterior {
	return bvhInterior{
		children:       [2]bvhNode{c1, c2},
		bounds:         c1.getBounds().AddAABB(c2.getBounds()),
		splitDimension: dim,
	}
}

func (n bvhInterior) getBounds() AABB {
	return n.bounds
}

func (n bvhInterior) getNumObjects() int {
	return 0
}

type objectInfo struct {
	index    int
	bounds   AABB
	centroid Vector
}

func NewBVH(objects []Object, splitFunc splitFunc) *BVH {
	objectInfos := make([]objectInfo, len(objects))
	for i, o := range objects {
		aabb := o.Bound(identity)
		objectInfos[i] = objectInfo{
			index:    i,
			bounds:   aabb,
			centroid: aabb.Centroid(),
		}
	}
	objectOrder := make([]int, len(objects))
	total, numObjects := 0, 0
	root := recursiveBuildBVH(objectInfos, 0, len(objectInfos), &numObjects, &total, objectOrder, splitFunc)
	orderedObjects := make([]Object, len(objects))
	for i, p := range objectOrder {
		orderedObjects[i] = objects[p]
	}

	nodes := make([]optimisedBVHNode, total)
	offset := 0
	flattenBVHTree(root, nodes, &offset)

	return &BVH{
		objects: orderedObjects,
		nodes:   nodes,
	}
}

// recursiveBuildBVH takes the list of objectInfo, a start and end index,
// the total so far of nodes created, the list of object indices to order
// and a function to split on, and returns the node that represents the
// range [start, end) and updated total
func recursiveBuildBVH(objectInfos []objectInfo, start, end int, objs, total *int, objectOrder []int, splitFunc splitFunc) bvhNode {
	*total++
	bounds := objectInfos[start].bounds
	for i := start + 1; i < end; i++ {
		bounds = bounds.AddAABB(objectInfos[i].bounds)
	}
	numObjects := end - start

	// Only one object remaining, return a leaf node
	if numObjects == 1 {
		offset := updateWithObjects(objectInfos, start, end, objs, objectOrder)
		return newBVHLeaf(offset, numObjects, bounds)
	}

	// Otherwise split in two sets and recurse
	centroidBounds := NewAABB(objectInfos[start].centroid, objectInfos[start+1].centroid)
	for i := start + 2; i < end; i++ {
		centroidBounds = centroidBounds.AddPoint(objectInfos[i].centroid)
	}
	dim := centroidBounds.MaximumExtent()

	// This case means all centroids are at the same position
	// further splitting would be ineffective
	if centroidBounds.Pmax.Get(dim) == centroidBounds.Pmin.Get(dim) {
		offset := updateWithObjects(objectInfos, start, end, objs, objectOrder)
		return newBVHLeaf(offset, numObjects, bounds)
	}

	splitIndex, createLeaf := splitFunc(objectInfos, start, end, dim, bounds, centroidBounds)
	if createLeaf {
		offset := updateWithObjects(objectInfos, start, end, objs, objectOrder)
		return newBVHLeaf(offset, numObjects, bounds)
	}
	c1 := recursiveBuildBVH(objectInfos, start, splitIndex, objs, total, objectOrder, splitFunc)
	c2 := recursiveBuildBVH(objectInfos, splitIndex, end, objs, total, objectOrder, splitFunc)
	return newBVHInterior(dim, c1, c2)
}

func updateWithObjects(objectInfos []objectInfo, start, end int, objs *int, order []int) int {
	firstOffset := *objs
	for i := start; i < end; i++ {
		objNum := objectInfos[i].index
		order[*objs] = objNum
		*objs++
	}
	return firstOffset
}

// A split function reorders the objects in the objects list
// between [start, end) and returns a split index called mid
// it also returns a bool param indicating whether to just create a leaf node instead
// additional parameters are the axis dimension to split on
// and the bounding box of centroids of all objects between start-end
type splitFunc func(objects []objectInfo, start, end int, axis Dimension, bounds, centroidBounds AABB) (index int, createLeaf bool)

// Find middle of bounding box along the axis
// order objectInfos by everything less than middle along the axis
// followed by everything greater than the middle along the axis.
// Return index of smallest node greater than middle
func SplitMiddle(objectInfos []objectInfo, start, end int, dim Dimension, bounds, centroidBounds AABB) (int, bool) {
	axisMid := (centroidBounds.Pmin.Get(dim) + centroidBounds.Pmax.Get(dim)) / 2
	// partition with pivot = axisMid
	i, j := start-1, end
	for {
		i++
		for objectInfos[i].centroid.Get(dim) < axisMid {
			i++
		}
		j--
		for objectInfos[j].centroid.Get(dim) > axisMid {
			j--
		}
		if i >= j {
			break
		}
		objectInfos[i], objectInfos[j] = objectInfos[j], objectInfos[i]
	}
	return i, false
}

type byCentroidDim struct {
	dim Dimension
	oi  []objectInfo
}

func (s byCentroidDim) Len() int {
	return len(s.oi)
}
func (s byCentroidDim) Swap(i, j int) {
	s.oi[i], s.oi[j] = s.oi[j], s.oi[i]
}
func (s byCentroidDim) Less(i, j int) bool {
	return s.oi[i].centroid.Get(s.dim) < s.oi[j].centroid.Get(s.dim)
}

// TODO: pbrt uses C++ std::nth_element, which does a partial sort in O(n)
// for now we will just sort the subarray
func SplitEqualCounts(objectInfos []objectInfo, start, end int, dim Dimension, bounds, centroidBounds AABB) (int, bool) {
	sort.Sort(byCentroidDim{dim: dim, oi: objectInfos[start:end]})
	mid := (start + end) / 2
	return mid, false
}

type bucketInfo struct {
	count  int
	bounds AABB
}

const nBuckets int = 12

func Split4SurfaceAreaHeuristic(objectInfos []objectInfo, start, end int, dim Dimension, bounds, centroidBounds AABB) (int, bool) {
	return splitSurfaceAreaHeuristic(objectInfos, start, end, dim, bounds, centroidBounds, true)
}

func SplitSurfaceAreaHeuristic(objectInfos []objectInfo, start, end int, dim Dimension, bounds, centroidBounds AABB) (int, bool) {
	return splitSurfaceAreaHeuristic(objectInfos, start, end, dim, bounds, centroidBounds, false)
}

// force4 boolean switches between the mbvh optimized version and default version of SAH split
func splitSurfaceAreaHeuristic(objectInfos []objectInfo, start, end int, dim Dimension, bounds, centroidBounds AABB, force4 bool) (int, bool) {
	nPrimitives := end - start
	if nPrimitives <= 4 {
		if force4 {
			return 0, true
		}
		return SplitEqualCounts(objectInfos, start, end, dim, bounds, centroidBounds)
	}
	// initialize buckets for SAH partition
	buckets := make([]bucketInfo, nBuckets)
	bucketMapping := make([]int, nPrimitives)
	for i := start; i < end; i++ {
		b := int(float32(nBuckets) * centroidBounds.Offset(objectInfos[i].centroid).Get(dim))
		if b == nBuckets {
			b = nBuckets - 1
		}
		bucketMapping[i-start] = b
		bucket := buckets[b]
		if bucket.count == 0 {
			bucket.bounds = objectInfos[i].bounds
		} else {
			bucket.bounds = bucket.bounds.AddAABB(objectInfos[i].bounds)
		}
		bucket.count = bucket.count + 1
		buckets[b] = bucket
	}
	// compute costs for splitting after each bucket
	// doesn't consider last bucket since splitting on it achieves nothing
	// estimated intersection cost = 1 and traversal cost = 1/8 or 0.125
	cost := make([]float32, nBuckets-1)
	for i := 0; i < nBuckets-1; i++ {
		b0 := buckets[0].bounds
		count0 := buckets[0].count
		for j := 1; j <= i; j++ {
			b0 = b0.AddAABB(buckets[j].bounds)
			count0 += buckets[j].count
		}
		b1 := buckets[i+1].bounds
		count1 := buckets[i+1].count
		for j := i + 2; j < nBuckets; j++ {
			b1 = b1.AddAABB(buckets[j].bounds)
			count1 += buckets[j].count
		}
		if count0%4 != 0 && count1%4 != 0 {
			cost[i] = math.MaxFloat32
		} else {
			cost[i] = 0.125 + (float32(count0)*b0.SurfaceArea()+float32(count1)*b1.SurfaceArea())/bounds.SurfaceArea()
		}
	}
	// find bucket to split at that minimizes SAH metric
	minCost := cost[0]
	minCostSplitBucket := 0
	for i := 1; i < nBuckets-1; i++ {
		if cost[i] < minCost {
			minCost = cost[i]
			minCostSplitBucket = i
		}
	}
	// either create leaf or split primitives at selected SAH bucket
	if !force4 {
		leafCost := float32(nPrimitives)
		if minCost >= leafCost {
			return 0, true
		}
	}
	// partition on bucket <= minCostSplitBucket
	i, j := start-1, end
	for {
		i++
		for bucketMapping[i-start] <= minCostSplitBucket {
			i++
		}
		j--
		for bucketMapping[j-start] > minCostSplitBucket {
			j--
		}
		if i >= j {
			break
		}
		objectInfos[i], objectInfos[j] = objectInfos[j], objectInfos[i]
		bucketMapping[i-start], bucketMapping[j-start] = bucketMapping[j-start], bucketMapping[i-start]
	}
	return i, false
}

type optimisedBVHNode struct {
	bounds     AABB
	offset     int       // firstOffset for leaf, secondChild for interior
	numObjects int       // interiorNodes have 0
	axis       Dimension // only for interiorNodes
	// TODO: padding? not optimised for exact size yet
}

// flattenBVHTree takes a rootnode of the (sub)tree, the nodeslist to fill,
// and an offset, and returns the updated offset
func flattenBVHTree(node bvhNode, nodes []optimisedBVHNode, offset *int) int {
	optimisedNode := optimisedBVHNode{
		bounds:     node.getBounds(),
		numObjects: node.getNumObjects(),
	}
	myOffset := *offset
	*offset++
	switch n := node.(type) {
	case bvhLeaf:
		optimisedNode.offset = n.firstOffset
	case bvhInterior:
		flattenBVHTree(n.children[0], nodes, offset)
		secondChildOffset := flattenBVHTree(n.children[1], nodes, offset)
		optimisedNode.axis = n.splitDimension
		optimisedNode.offset = secondChildOffset
	}
	nodes[myOffset] = optimisedNode
	return myOffset
}

// Actual traversal of the BVH
func (bvh *BVH) ClosestIntersection(ray Ray, maxDistance float32) (*SurfaceInteraction, bool) {
	var toVisitOffset, currentNodeIndex int
	var found bool
	var surfaceInteraction *SurfaceInteraction
	distance := maxDistance
	nodesToVisit := make([]int, 64)
	for {
		node := bvh.nodes[currentNodeIndex]
		if tMin, hit := node.bounds.Intersect(ray); hit && tMin < maxDistance {
			if node.numObjects > 0 {
				// this is a leaf node
				for i := 0; i < node.numObjects; i++ {
					o := bvh.objects[node.offset+i]
					if si, ok := o.Intersect(ray); ok && si.distance < distance && si.distance > ERROR_MARGIN {
						distance = si.distance
						found = true
						surfaceInteraction = si
					}
				}
				if toVisitOffset == 0 {
					break
				}
				toVisitOffset--
				currentNodeIndex = nodesToVisit[toVisitOffset]
			} else {
				// this is an interior node
				// TODO: optimisation regarding direction of ray and child visit order
				nodesToVisit[toVisitOffset] = currentNodeIndex + 1
				toVisitOffset++
				currentNodeIndex = node.offset

			}
		} else {
			if toVisitOffset == 0 {
				break
			}
			toVisitOffset--
			currentNodeIndex = nodesToVisit[toVisitOffset]
		}
	}
	if !found {
		return nil, false
	}
	return surfaceInteraction, true
}

// Bounded Volume Hierarchy for triangles only!
type TriangleBVH struct {
	triangles []Triangle
	nodes     []optimisedBVHNode
}

func (bvh *TriangleBVH) GetObjects() []Object {
	objects := make([]Object, len(bvh.triangles))
	for i, t := range bvh.triangles {
		objects[i] = t
	}
	return objects
}

func NewTriangleBVH(triangles []Triangle, splitFunc splitFunc) *TriangleBVH {
	objectInfos := make([]objectInfo, len(triangles))
	for i, o := range triangles {
		aabb := o.Bound(identity)
		objectInfos[i] = objectInfo{
			index:    i,
			bounds:   aabb,
			centroid: aabb.Centroid(),
		}
	}
	triangleOrder := make([]int, len(triangles))
	total, numTriangles := 0, 0
	root := recursiveBuildBVH(objectInfos, 0, len(objectInfos), &numTriangles, &total, triangleOrder, splitFunc)
	orderedTriangles := make([]Triangle, len(triangles))
	for i, p := range triangleOrder {
		orderedTriangles[i] = triangles[p]
	}

	nodes := make([]optimisedBVHNode, total)
	offset := 0
	flattenBVHTree(root, nodes, &offset)

	return &TriangleBVH{
		triangles: orderedTriangles,
		nodes:     nodes,
	}
}

func (bvh *TriangleBVH) ClosestIntersection(ray Ray, maxDistance float32) (*SurfaceInteraction, bool) {
	var toVisitOffset, currentNodeIndex int
	var found bool
	var triangle Triangle
	distance := maxDistance
	nodesToVisit := make([]int, 64)
	for {
		node := bvh.nodes[currentNodeIndex]
		if tMin, hit := node.bounds.Intersect(ray); hit && tMin < maxDistance {
			if node.numObjects > 0 {
				// this is a leaf node
				for i := 0; i < node.numObjects; i++ {
					t := bvh.triangles[node.offset+i]
					if d, ok := t.IntersectOptimized(ray); ok && d < distance && d > ERROR_MARGIN {
						distance = d
						found = true
						triangle = t
					}
				}
				if toVisitOffset == 0 {
					break
				}
				toVisitOffset--
				currentNodeIndex = nodesToVisit[toVisitOffset]
			} else {
				// this is an interior node
				// TODO: optimisation regarding direction of ray and child visit order
				nodesToVisit[toVisitOffset] = currentNodeIndex + 1
				toVisitOffset++
				currentNodeIndex = node.offset

			}
		} else {
			if toVisitOffset == 0 {
				break
			}
			toVisitOffset--
			currentNodeIndex = nodesToVisit[toVisitOffset]
		}
	}
	if !found {
		return nil, false
	}
	normal := triangle.SurfaceNormal(PointFromRay(ray, distance))
	return NewSurfaceInteraction(triangle, distance, normal, ray), true
}

// Bounded Volume Hierarchy for triangles only, simd optimized
type Triangle4BVH struct {
	triangles []Triangle
	nodes     []bvh4Node
}

func (bvh *Triangle4BVH) GetObjects() []Object {
	objects := make([]Object, len(bvh.triangles))
	for i, t := range bvh.triangles {
		objects[i] = t
	}
	return objects
}

func NewTriangle4BVH(triangles []Triangle) *Triangle4BVH {
	objectInfos := make([]objectInfo, len(triangles))
	for i, o := range triangles {
		aabb := o.Bound(identity)
		objectInfos[i] = objectInfo{
			index:    i,
			bounds:   aabb,
			centroid: aabb.Centroid(),
		}
	}
	triangleOrder := make([]int, len(triangles))
	total, numTriangles := 0, 0
	binaryRoot := recursiveBuildBVH(objectInfos, 0, len(objectInfos), &numTriangles, &total, triangleOrder, Split4SurfaceAreaHeuristic)
	orderedTriangles := make([]Triangle, len(triangles))
	for i, p := range triangleOrder {
		orderedTriangles[i] = triangles[p]
	}

	nodes := make([]bvh4Node, total)
	offset := 0
	convertTo4mbvhInterior(binaryRoot.(bvhInterior), orderedTriangles, nodes, &offset)

	return &Triangle4BVH{
		triangles: orderedTriangles,
		nodes:     nodes,
	}
}

type bvh4Node interface {
	isLeaf() bool
}

type bvh4Interior struct {
	// NOTE: offset -1 means empty node
	childOffsets [4]int
	// boundsData in simd
	min4x [4]float32
	min4y [4]float32
	min4z [4]float32
	max4x [4]float32
	max4y [4]float32
	max4z [4]float32
	// TODO: traversalOrder
}

func (bvh4Interior) isLeaf() bool {
	return false
}

type bvh4Leaf struct {
	firstOffset  int
	numTriangles int
	// triangleData in simd
	p0x [4]float32
	p0y [4]float32
	p0z [4]float32
	p1x [4]float32
	p1y [4]float32
	p1z [4]float32
	p2x [4]float32
	p2y [4]float32
	p2z [4]float32
}

func (bvh4Leaf) isLeaf() bool {
	return false
}

// converts bvh to 4-mbvh according to Multi Bounding Volume Hierarchies - Ernst, Greiner 2008
// also flattens the tree into a node list
func convertTo4mbvhInterior(node bvhInterior, triangles []Triangle, nodes []bvh4Node, offset *int) int {
	myOffset := *offset
	*offset++
	nc1 := node.children[0]
	o1, b1, o2, b2 := convertTo4mbvhIntermediate(nc1, triangles, nodes, offset)
	nc2 := node.children[1]
	o3, b3, o4, b4 := convertTo4mbvhIntermediate(nc2, triangles, nodes, offset)
	childOffsets := [4]int{o1, o2, o3, o4}
	min4x, min4y, min4z, max4x, max4y, max4z := boundsToSimd(b1, b2, b3, b4)
	nodes[myOffset] = bvh4Interior{
		childOffsets: childOffsets,
		min4x:        min4x,
		min4y:        min4y,
		min4z:        min4z,
		max4x:        max4x,
		max4y:        max4y,
		max4z:        max4z,
	}
	return myOffset
}

func convertTo4mbvhIntermediate(node bvhNode, triangles []Triangle, nodes []bvh4Node, offset *int) (int, AABB, int, AABB) {
	switch n := node.(type) {
	case bvhLeaf:
		o := convertTo4mbvhLeaf(n, triangles, nodes, offset)
		return o, n.getBounds(), -1, AABB{}
	case bvhInterior:
		var o1, o2 int
		nc1 := n.children[0]
		b1 := nc1.getBounds()
		switch n1 := nc1.(type) {
		case bvhLeaf:
			o1 = convertTo4mbvhLeaf(n1, triangles, nodes, offset)
		case bvhInterior:
			o1 = convertTo4mbvhInterior(n1, triangles, nodes, offset)
		}
		nc2 := n.children[1]
		b2 := nc2.getBounds()
		switch n2 := nc2.(type) {
		case bvhLeaf:
			o2 = convertTo4mbvhLeaf(n2, triangles, nodes, offset)
		case bvhInterior:
			o2 = convertTo4mbvhInterior(n2, triangles, nodes, offset)
		}
		return o1, b1, o2, b2
	}
	panic("convert to 4-mbvh: assumptions failed!")
	return 0, AABB{}, 0, AABB{}
}

func convertTo4mbvhLeaf(node bvhLeaf, triangles []Triangle, nodes []bvh4Node, offset *int) int {
	trs := [4]Triangle{}
	for i := 0; i < node.numObjects; i++ {
		trs[i] = triangles[node.firstOffset+i]
	}
	p0x, p0y, p0z, p1x, p1y, p1z, p2x, p2y, p2z := trianglesToSimd(trs[0], trs[1], trs[2], trs[3])

	myOffset := *offset
	nodes[myOffset] = bvh4Leaf{
		firstOffset:  node.firstOffset,
		numTriangles: node.numObjects,
		p0x:          p0x,
		p0y:          p0y,
		p0z:          p0z,
		p1x:          p1x,
		p1y:          p1y,
		p1z:          p1z,
		p2x:          p2x,
		p2y:          p2y,
		p2z:          p2z,
	}
	*offset++
	return myOffset
}

func (bvh *Triangle4BVH) ClosestIntersection(ray Ray, maxDistance float32) (*SurfaceInteraction, bool) {
	rox, roy, roz, rdx, rdy, rdz := ray.toSimd()
	var toVisitOffset, currentNodeIndex int
	var found bool
	var triangle Triangle
	distance := maxDistance
	nodesToVisit := make([]int, 64)
	for {
		n := bvh.nodes[currentNodeIndex].(bvh4Interior)
		boxTs := simd.Box4Intersect(rox, roy, roz, rdx, rdy, rdz, n.min4x, n.min4y, n.min4z, n.max4x, n.max4y, n.max4z)
		// TODO: traversal order
		for i := 0; i < 4; i++ {
			boxT := boxTs[i]
			if boxT == 0 || boxT > distance {
				continue
			}
			offset := n.childOffsets[i]
			// empty child node
			if offset == -1 {
				continue
			}
			child := bvh.nodes[offset]
			switch c := child.(type) {
			case bvh4Leaf:
				ts := simd.Triangle4Intersect(c.p0x, c.p0y, c.p0z, c.p1x, c.p1y, c.p1z, c.p2x, c.p2y, c.p2z, rox, roy, roz, rdx, rdy, rdz)
				for j := 0; j < c.numTriangles; j++ {
					t := ts[j]
					if t == 0 || t > distance || t <= ERROR_MARGIN {
						continue
					}
					triangle = bvh.triangles[c.firstOffset+j]
					found = true
					distance = t
				}
			case bvh4Interior:
				nodesToVisit[toVisitOffset] = offset
				toVisitOffset++
			}
		}
		if toVisitOffset == 0 {
			break
		}
		toVisitOffset--
		currentNodeIndex = nodesToVisit[toVisitOffset]
	}
	if !found {
		return nil, false
	}
	normal := triangle.SurfaceNormal(PointFromRay(ray, distance))
	return NewSurfaceInteraction(triangle, distance, normal, ray), true
}
