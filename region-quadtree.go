package karytree

const ( // region quadrants
	ne = iota
	nw = iota
	sw = iota
	se = iota
)

// RegionQuadtree creates a region quadtree karytree.Node
func RegionQuadtree(key interface{}) Node {
	return NewNode(4, key)
}

// SetNW sets the north-west quadrant
func (k *Node) SetNW(other *Node) {
	k.SetNthChild(nw, other)
}

// SetNE sets the north-east quadrant
func (k *Node) SetNE(other *Node) {
	k.SetNthChild(ne, other)
}

// SetSE sets the south-east child.
func (k *Node) SetSE(other *Node) {
	k.SetNthChild(se, other)
}

// SetSW sets the south-west child.
func (k *Node) SetSW(other *Node) {
	k.SetNthChild(sw, other)
}

// NW gets the nw child
func (k *Node) NW() *Node {
	return k.NthChild(nw)
}

// NE gets the ne child
func (k *Node) NE() *Node {
	return k.NthChild(ne)
}

// SE gets the se child
func (k *Node) SE() *Node {
	return k.NthChild(se)
}

// SW gets the sw child
func (k *Node) SW() *Node {
	return k.NthChild(sw)
}
