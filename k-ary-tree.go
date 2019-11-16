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
	key         interface{}
	firstChild  *Node
	nextSibling *Node
}

// New creates a new node with k = k and data key. []*Node children
// is an uninitialized slice.
//
// Don't use nil keys, these are used as a sentinel in the library.
func New(k int, key interface{}) Node {
	n := Node{}
	n.key = key
	return n
}

// K returns the k value of a tree node
//
// Deprecated: switched to a sibling-sibling binary tree representation
// where k is no longer stored.
func (k *Node) K() int {
	return 0
}

// SetNthChild sets the Nth child. If an existing node is replaced,
// that node is returned.
func (k *Node) SetNthChild(n int, other *Node) *Node {
	if n == 0 {
		ret := k.firstChild
		k.firstChild = other
		if ret != nil {
			k.firstChild.nextSibling = ret.nextSibling
		}
		return ret
	}

	if k.firstChild == nil {
		next := New(0, nil)
		k.firstChild = &next
	}
	curr := k.firstChild

	for nLocal := 1; nLocal != n; nLocal++ {
		if curr.nextSibling == nil {
			next := New(0, nil)
			curr.nextSibling = &next
		}
		curr = curr.nextSibling
	}

	ret := curr.nextSibling
	curr.nextSibling = other
	if ret != nil {
		curr.nextSibling.nextSibling = ret.nextSibling
	}
	return ret
}

// NthChild gets the Nth child.
func (k *Node) NthChild(n int) *Node {
	if n == 0 {
		if k.firstChild != nil && k.firstChild.key == nil {
			return nil // this is our own sentinel
		}
		return k.firstChild
	}

	if k.firstChild == nil {
		return nil
	}
	curr := k.firstChild

	for nLocal := 1; nLocal != n; nLocal++ {
		if curr.nextSibling == nil {
			return nil
		}
		curr = curr.nextSibling
	}

	ret := curr.nextSibling
	if ret != nil && ret.key == nil {
		return nil
	}
	return ret
}

// Key gets the data stored in a node
func (k *Node) Key() interface{} {
	return k.key
}

// SetKey modifies the data in a node
func (k *Node) SetKey(newKey interface{}) {
	k.key = newKey
}
