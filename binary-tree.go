package karytree

const (
	left  = iota
	right = iota
)

// Binary creates a binary karytree.Node
func Binary(key KeyType) Node {
	return NewNode(key)
}

// SetLeft sets the left child.
func (k *Node) SetLeft(other *Node) {
	k.SetNthChild(left, other)
}

// SetRight sets the left child.
func (k *Node) SetRight(other *Node) {
	k.SetNthChild(right, other)
}

// Left gets the left child
func (k *Node) Left() *Node {
	return k.NthChild(left)
}

// Right gets the right child
func (k *Node) Right() *Node {
	return k.NthChild(right)
}
