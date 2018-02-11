package model

// Bounded Volume Hierarchy
type BVH struct {
	objects []Object
	nodes   []optimisedBVHNode
}

// A split function reorders the objects in the primitives list
// between [start, end) and returns a split index called mid
// additional parameters are the axis dimension to split on
// and the bounding box of centroids of all objects between start-end
type splitFunc func(primitives []primitiveInfo, start, end int, axis Dimension, bounds AABB) int

type bvhNode interface {
	getBounds() AABB
	getNumPrimitives() int
}

type bvhLeaf struct {
	firstPrimOffset int
	numPrimitives   int
	bounds          AABB
}

func newBVHLeaf(first, n int, bounds AABB) bvhLeaf {
	return bvhLeaf{
		firstPrimOffset: first,
		numPrimitives:   n,
		bounds:          bounds,
	}
}

func (n bvhLeaf) getBounds() AABB {
	return n.bounds
}

func (n bvhLeaf) getNumPrimitives() int {
	return n.numPrimitives
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

func (n bvhInterior) getNumPrimitives() int {
	return 0
}

type primitiveInfo struct {
	index    int
	bounds   AABB
	centroid Vector
}

func NewBVH(objects []Object, splitFunc splitFunc) BVH {
	primitives := make([]primitiveInfo, len(objects))
	for i, o := range objects {
		aabb := o.Bound()
		primitives[i] = primitiveInfo{
			index:    i,
			bounds:   aabb,
			centroid: aabb.Centroid(),
		}
	}
	primitiveOrder := make([]int, len(objects))
	total, numPrimitives := 0, 0
	root := recursiveBuildBVH(primitives, 0, len(primitives), &numPrimitives, &total, primitiveOrder, splitFunc)
	orderedObjects := make([]Object, len(objects))
	for i, p := range primitiveOrder {
		orderedObjects[i] = objects[p]
	}

	nodes := make([]optimisedBVHNode, total)
	offset := 0
	flattenBVHTree(root, nodes, &offset)

	return BVH{
		objects: orderedObjects,
		nodes:   nodes,
	}
}

// recursiveBuildBVH takes the list of primitives, a start and end index,
// the total so far of nodes created, the list of object indices to order
// and a function to split on, and returns the node that represents the
// range [start, end) and updated total
func recursiveBuildBVH(primitives []primitiveInfo, start, end int, prims, total *int, primitiveOrder []int, splitFunc splitFunc) bvhNode {
	*total++
	bounds := primitives[start].bounds
	for i := start + 1; i < end; i++ {
		bounds.AddAABB(primitives[i].bounds)
	}
	numPrimitives := end - start

	// Only one primitive remaining, return a leaf node
	if numPrimitives == 1 {
		offset := updateWithPrimitives(primitives, start, end, prims, primitiveOrder)
		return newBVHLeaf(offset, numPrimitives, bounds)
	}

	// Otherwise split in two sets and recurse
	centroidBounds := NewAABB(primitives[start].centroid, primitives[start+1].centroid)
	for i := start + 2; i < end; i++ {
		centroidBounds = centroidBounds.AddPoint(primitives[i].centroid)
	}
	dim := centroidBounds.MaximumExtent()

	// This case means all centroids are at the same position
	// further splitting would be ineffective
	if centroidBounds.Pmax.Get(dim) == centroidBounds.Pmin.Get(dim) {
		offset := updateWithPrimitives(primitives, start, end, prims, primitiveOrder)
		return newBVHLeaf(offset, numPrimitives, bounds)
	}

	mid := splitFunc(primitives, start, end, dim, centroidBounds)
	c1 := recursiveBuildBVH(primitives, start, mid, prims, total, primitiveOrder, splitFunc)
	c2 := recursiveBuildBVH(primitives, mid, end, prims, total, primitiveOrder, splitFunc)
	return newBVHInterior(dim, c1, c2)
}

func updateWithPrimitives(primitives []primitiveInfo, start, end int, prims *int, order []int) int {
	firstOffset := *prims
	for i := start; i < end; i++ {
		primNum := primitives[i].index
		order[*prims] = primNum
		*prims++
	}
	return firstOffset
}

// Find middle of bounding box along the axis
// order primitives by everything less than middle along the axis
// followed by everything greater than the middle along the axis.
// Return index of smallest node greater than middle
func SplitMiddle(primitives []primitiveInfo, start, end int, dim Dimension, centroidBounds AABB) int {
	axisMid := (centroidBounds.Pmin.Get(dim) + centroidBounds.Pmax.Get(dim)) / 2
	// partition with pivot = axisMid
	i, j := start-1, end
	for {
		i++
		for primitives[i].centroid.Get(dim) < axisMid {
			i++
		}
		j--
		for primitives[j].centroid.Get(dim) > axisMid {
			j--
		}
		if i >= j {
			break
		}
		primitives[i], primitives[j] = primitives[j], primitives[i]
	}
	return i
}

type optimisedBVHNode struct {
	bounds        AABB
	offset        int       // firstPrimitive for leaf, secondChild for interior
	numPrimitives int       // interiorNodes have 0
	axis          Dimension // only for interiorNodes
	// TODO: padding? not optimised for exact size yet
}

// flattenBVHTree takes a rootnode of the (sub)tree, the nodeslist to fill,
// and an offset, and returns the updated offset
func flattenBVHTree(node bvhNode, nodes []optimisedBVHNode, offset *int) int {
	optimisedNode := optimisedBVHNode{
		bounds:        node.getBounds(),
		numPrimitives: node.getNumPrimitives(),
	}
	myOffset := *offset
	*offset++
	switch n := node.(type) {
	case bvhLeaf:
		optimisedNode.offset = n.firstPrimOffset
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
func (bvh BVH) ClosestIntersection(ray Ray, maxDistance float64) *hit {
	var toVisitOffset, currentNodeIndex int
	d := maxDistance
	var objectHit Object
	nodesToVisit := make([]int, 64)
	for {
		node := bvh.nodes[currentNodeIndex]
		// TODO: disregard box intersection with distance of tMin
		// farther than current closest intersection distance d
		if node.bounds.Intersect(ray) {
			if node.numPrimitives > 0 {
				// this is a leaf node
				newD := d
				for i := 0; i < node.numPrimitives; i++ {
					o := bvh.objects[node.offset+i]
					if distance, ok := o.Intersect(ray); ok && distance < newD && distance > ERROR_MARGIN {
						newD = distance
						objectHit = o
					}
				}
				d = newD
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
	if d == maxDistance {
		return nil
	}
	return &hit{
		object: objectHit,
		point:  PointFromRay(ray, d),
	}
}
