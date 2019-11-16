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
	parent      *Node
	firstChild  *Node
	nextSibling *Node
	n           int
}

// New creates a new node with k = k and data key. []*Node children
// is an uninitialized slice.
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
	if other.parent != nil {
		panic("node already has a parent")
	}

	other.n = n
	other.parent = k

	if k.firstChild == nil {
		k.firstChild = other
		return nil
	}

	if n == 0 {
		if k.firstChild.n > n {
			// relink
			other.nextSibling = k.firstChild
			k.firstChild = other
			return nil
		} else if k.firstChild.n == n {
			// evict
			other.nextSibling = k.firstChild.nextSibling
			ret := k.firstChild
			ret.parent = nil //unparent it
			k.firstChild = other
			return ret
		}
	}

	if k.firstChild.n == n {
		// evict
		ret := k.firstChild
		other.nextSibling = k.firstChild.nextSibling
		k.firstChild = other
		ret.parent = nil //unparent it
		return ret
	} else if k.firstChild.n > n {
		// relink
		other.nextSibling = k.firstChild
		k.firstChild = other
		return nil
	}

	curr := k.firstChild
	for {
		if curr.nextSibling == nil {
			curr.nextSibling = other
			return nil
		}
		if curr.nextSibling.n == n {
			/* evict the existing nth child

			 *       other
			 * curr -> nextSibling -> ...
			 *
			 * curr -> other -> ..., return nextSibling
			 */
			other.nextSibling = curr.nextSibling.nextSibling
			curr.nextSibling.nextSibling = nil // wipe the rest of the links from the evicted node
			ret := curr.nextSibling
			curr.nextSibling = other
			ret.parent = nil //unparent the evicted node
			return ret
		} else if curr.nextSibling.n > n {
			/* relink
			 *       other
			 * curr -> nextSibling
			 *
			 * curr -> other -> nextSibling
			 */
			other.nextSibling = curr.nextSibling
			curr.nextSibling = other
			return nil
		}
		// keep traversing the sibling linkedlist
		curr = curr.nextSibling
	}
}

// NthChild gets the Nth child.
func (k *Node) NthChild(n int) *Node {
	curr := k.firstChild
	for curr != nil {
		if curr.n == n {
			// exact match
			return curr
		} else if curr.n > n {
			// overshoot, nth child doesn't exist
			return nil
		}
		// keep traversing the sibling linkedlist
		curr = curr.nextSibling
	}

	return nil
}

// Key gets the data stored in a node
func (k *Node) Key() interface{} {
	return k.key
}

// SetKey modifies the data in a node
func (k *Node) SetKey(newKey interface{}) {
	k.key = newKey
}
