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
	//"unsafe"
	//"fmt"
)

//var (
//	nBitmask uintptr = 0xFFFF000000000000
//	firstChildBitmask uintptr = 0x0000FFFFFFFFFFFF
//	nBitshift uint16 = 48
//)

type KeyType generic.Type

// A Node is a typical recursive tree node, and it represents a tree
// when it's traversed. The key is for data stored in the node.
type Node struct {
	key         KeyType
	n_			uint16
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
	//return uint16((k.firstChild & nBitmask) >> nBitshift)
	return k.n_
}

func (k *Node) setN(n uint16) {
	//k.firstChild = (k.firstChild & firstChildBitmask) | (uintptr(n) << nBitshift)
	k.n_ = n
}

func (k *Node) getFirstChild() *Node {
	//return (*Node)(unsafe.Pointer(k.firstChild & firstChildBitmask))
	return k.firstChild
}

func (k *Node) setFirstChild(child *Node) {
	//k.firstChild = (k.firstChild & nBitmask) | uintptr(unsafe.Pointer(child))
	k.firstChild = child
}

// SetNthChild sets the Nth child. If an existing node is replaced,
// that node is returned.
func (k *Node) SetNthChild(n uint16, other *Node) *Node {
	//use top 16 bits of firstChild pointer to store 'n'
	other.setN(n)

	if k.getFirstChild() == nil {
		k.setFirstChild(other)
		return nil
	}

	if k.getFirstChild().n() == n {
		// evict
		ret := k.getFirstChild()
		other.nextSibling = k.getFirstChild().nextSibling
		k.setFirstChild(other)
		return ret
	} else if k.getFirstChild().n() > n {
		// relink
		other.nextSibling = k.getFirstChild()
		k.setFirstChild(other)
		return nil
	}

	curr := k.getFirstChild()
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
	curr := k.getFirstChild()
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
