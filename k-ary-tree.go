/*
Package karytree implements a recursive k-ary tree data structure.

The children of the tree is a slice of Node pointers.

The slice of children is originally uninitialized, but is filled
with k nils when any one child is set.

The tree also stores the value of k. The caller is responsible
for in-range indexing (e.g. don't set a 5th child of a k=3 node).
*/
package karytree

// A Node is a typical recursive tree node, and it represents a tree
// when it's traversed. The key is for data stored in the node.
type Node struct {
	key      interface{}
	children []*Node
	k        int
}

// New creates a new node with k = k and data key. []*Node children
// is an uninitialized slice.
func New(k int, key interface{}) Node {
	n := Node{}

	n.k = k
	n.key = key
	return n
}

// K returns the k value of a tree node
func (k *Node) K() int {
	return k.k
}

// SetNthChild sets the Nth child. If the []*Node children slice is
// uninitialized (i.e. first time setting a child), it's first initialized
// to [nil]*k.
func (k *Node) SetNthChild(n int, other *Node) {
	if k.children == nil {
		k.children = make([]*Node, k.k)
	}

	k.children[n] = other
}

// NthChild gets the Nth child.
func (k *Node) NthChild(n int) *Node {
	if k.children == nil {
		return nil
	}

	return k.children[n]
}

// Key gets the data stored in a node
func (k *Node) Key() interface{} {
	return k.key
}

// SetKey modifies the data in a node
func (k *Node) SetKey(newKey interface{}) {
	k.key = newKey
}
