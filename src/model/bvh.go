package model

// Bounded Volume Hierarchy
type BVH struct {
	objects []Object
	nodes   []optimisedBVHNode
}

func (bvh *BVH) GetObjects() []Object {
	return bvh.objects
}

// A split function reorders the objects in the objects list
// between [start, end) and returns a split index called mid
// additional parameters are the axis dimension to split on
// and the bounding box of centroids of all objects between start-end
type splitFunc func(objects []objectInfo, start, end int, axis Dimension, bounds AABB) int

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

// max recursion depth, see recursiveBuildBVH function
var maxDepth = 10

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
	root := recursiveBuildBVH(objectInfos, 0, len(objectInfos), &numObjects, &total, objectOrder, splitFunc, maxDepth)
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
// TODO: max recursion depth added as a fix to this function causing stack overflow
// this probably happens when objects overlap, not sure why though.
// underlying cause should be found and fixed instead
func recursiveBuildBVH(objectInfos []objectInfo, start, end int, objs, total *int, objectOrder []int, splitFunc splitFunc, depth int) bvhNode {
	*total++
	bounds := objectInfos[start].bounds
	for i := start + 1; i < end; i++ {
		bounds = bounds.AddAABB(objectInfos[i].bounds)
	}
	numObjects := end - start

	// Only one object remaining, return a leaf node or recursion depth reached
	if numObjects == 1 || depth == 0 {
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

	mid := splitFunc(objectInfos, start, end, dim, centroidBounds)
	c1 := recursiveBuildBVH(objectInfos, start, mid, objs, total, objectOrder, splitFunc, depth-1)
	c2 := recursiveBuildBVH(objectInfos, mid, end, objs, total, objectOrder, splitFunc, depth-1)
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

// Find middle of bounding box along the axis
// order objectInfos by everything less than middle along the axis
// followed by everything greater than the middle along the axis.
// Return index of smallest node greater than middle
func SplitMiddle(objectInfos []objectInfo, start, end int, dim Dimension, centroidBounds AABB) int {
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
	return i
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
func (bvh *BVH) ClosestIntersection(ray Ray, maxDistance float64) (*SurfaceInteraction, bool) {
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
// 4-ary tree
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
	root := recursiveBuildBVH(objectInfos, 0, len(objectInfos), &numTriangles, &total, triangleOrder, splitFunc, maxDepth)
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

func (bvh *TriangleBVH) ClosestIntersection(ray Ray, maxDistance float64) (*SurfaceInteraction, bool) {
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
