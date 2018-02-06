package model

// NOTE:
// multiple places in this code where I pass slices
// which could be optimised by using pointers?
// need to investigate further
// (mostly in recursiveBuildBVH and flattenBVHTree)

// Bounded Volume Hierarchy
type BVH struct {
	objects []Object
	nodes   []optimisedBVHNode
}

// A split function reorders the objects in the primitives list
// between [start, end) and returns a split index called mid
type splitFunc func(primitives []primitiveInfo, start, end int) ([]primitiveInfo, int)

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
	root, primitiveOrder, total := recursiveBuildBVH(primitives, 0, len(primitives), 0, []int{}, splitFunc)
	orderedObjects := make([]Object, len(objects))
	for i, p := range primitiveOrder {
		orderedObjects[i] = objects[p]
	}

	nodes := make([]optimisedBVHNode, total)
	nodes, _ = flattenBVHTree(root, nodes, 0)

	return BVH{
		objects: orderedObjects,
		nodes:   nodes,
	}
}

// recursiveBuildBVH takes the list of primitives, a start and end index,
// the total so far of nodes created, the list of object indices to order
// and a function to split on, and returns the node that represents the
// range [start, end), updated order list and updated total
func recursiveBuildBVH(primitives []primitiveInfo, start, end, total int, primitiveOrder []int, splitFunc splitFunc) (bvhNode, []int, int) {
	total++
	bounds := primitives[start].bounds
	for i := start + 1; i < end; i++ {
		bounds.AddAABB(primitives[i].bounds)
	}
	numPrimitives := end - start

	// Only one primitive remaining, return a leaf node
	if numPrimitives == 1 {
		for i := start; i < end; i++ {
			primNum := primitives[i].index
			primitiveOrder = append(primitiveOrder, primNum)
		}
		return newBVHLeaf(len(primitiveOrder), numPrimitives, bounds), primitiveOrder, total
	}

	// Otherwise split in two sets and recurse
	centroidBounds := NewAABB(primitives[start].centroid, primitives[start+1].centroid)
	for i := start + 2; i < end; i++ {
		centroidBounds.AddPoint(primitives[i].centroid)
	}
	dim := centroidBounds.MaximumExtent()

	// This case means all centroids are at the same position
	// further splitting would be ineffective
	if centroidBounds.Pmax.Get(dim) == centroidBounds.Pmin.Get(dim) {
		for i := start; i < end; i++ {
			primNum := primitives[i].index
			primitiveOrder = append(primitiveOrder, primNum)
		}
		return newBVHLeaf(len(primitiveOrder), numPrimitives, bounds), primitiveOrder, total
	}

	primitives, mid := splitFunc(primitives, start, end)
	c1, primitiveOrder, total := recursiveBuildBVH(primitives, start, mid, total, primitiveOrder, splitFunc)
	c2, primitiveOrder, total := recursiveBuildBVH(primitives, mid, end, total, primitiveOrder, splitFunc)
	return newBVHInterior(dim, c1, c2), primitiveOrder, total
}

func SplitTODO(primitives []primitiveInfo, start, end int) ([]primitiveInfo, int) {
	mid := (start + end) / 2
	return primitives, mid
}

type optimisedBVHNode struct {
	bounds        AABB
	offset        int       // firstPrimitive for leaf, secondChild for interior
	numPrimitives int       // interiorNodes have 0
	axis          Dimension // only for interiorNodes
	// TODO: padding? not optimised for exact size yet
}

// flattenBVHTree takes a rootnode of the (sub)tree, the nodeslist to fill,
// and an offset, and returns the (more) filled nodeslist and updated offset
func flattenBVHTree(node bvhNode, nodes []optimisedBVHNode, offset int) ([]optimisedBVHNode, int) {
	offset++
	var optimisedNode optimisedBVHNode
	switch n := node.(type) {
	case bvhLeaf:
		optimisedNode.offset = n.firstPrimOffset
		optimisedNode.numPrimitives = n.getNumPrimitives()
	case bvhInterior:
		nodes, offset = flattenBVHTree(n.children[0], nodes, offset)
		nodes, offset = flattenBVHTree(n.children[1], nodes, offset)
		optimisedNode.axis = n.splitDimension
		optimisedNode.offset = offset
		optimisedNode.numPrimitives = 0
	}
	return nodes, offset
}

func (bvh BVH) ClosestIntersection(ray Ray) *hit {
	// TODO:
	// - ray-AABB intersection
	// - actual traversal
	return nil
}
