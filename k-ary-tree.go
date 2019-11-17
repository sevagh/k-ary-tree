/*
Package karytree implements a recursive k-ary tree data structure.

The children of the tree is a slice of Node pointers.

The slice of children is originally uninitialized, but is filled
with k nils when any one child is set.

The tree also stores the value of k. The caller is responsible
for in-range indexing (e.g. don't set a 5th child of a k=3 node).
*/
package karytree

import (
	"github.com/cheekybits/genny/generic"
	"unsafe"
)

type KeyType generic.Type

// A Node is a typical recursive tree node, and it represents a tree
// when it's traversed. The key is for data stored in the node.
type Node struct {
	key         KeyType
	parent      uintptr
	firstChild  *Node
	nextSibling *Node
}

// NewNode creates a new node data key.
func NewNode(key KeyType) Node {
	n := Node{}
	n.key = key
	return n
}

func (k *Node) n() uint16 {
	return uint16((k.parent & 0xFFFF000000000000) >> 48)
}

// Parent gets a pointer to the parent of a Node.
func (k *Node) Parent() *Node {
	return (*Node)(unsafe.Pointer(k.parent & 0x0000FFFFFFFFFFFF))
}

func (k *Node) setParent(n uint16, child *Node) {
	k.parent = uintptr(unsafe.Pointer(child)) | (uintptr(n) << 48)
}

// SetNthChild sets the Nth child. If an existing node is replaced,
// that node is returned.
func (k *Node) SetNthChild(n uint16, other *Node) *Node {
	//use top 16 bits of pointer to store 'n'
	if other.Parent() != nil {
		panic("same node can't have different parents")
	}
	other.setParent(n, k)

	if k.firstChild == nil {
		k.firstChild = other
		return nil
	}

	if k.firstChild.n() == n {
		// evict
		ret := k.firstChild
		other.nextSibling = k.firstChild.nextSibling
		k.firstChild = other
		ret.setParent(0, nil)
		return ret
	} else if k.firstChild.n() > n {
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
		if curr.nextSibling.n() == n {
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
			ret.setParent(0, nil)
			return ret
		} else if curr.nextSibling.n() > n {
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
func (k *Node) NthChild(n uint16) *Node {
	curr := k.firstChild
	for curr != nil {
		if curr.n() == n {
			// exact match
			return curr
		} else if curr.n() > n {
			// overshoot, nth child doesn't exist
			return nil
		}
		// keep traversing the sibling linkedlist
		curr = curr.nextSibling
	}

	return nil
}

// Key gets the data stored in a node
func (k *Node) Key() KeyType {
	return k.key
}

// SetKey modifies the data in a node
func (k *Node) SetKey(newKey KeyType) {
	k.key = newKey
}
