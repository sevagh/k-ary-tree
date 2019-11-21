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
	n           uint
	firstChild  *Node
	nextSibling *Node
}

// NewNode creates a new node data key.
func NewNode(key interface{}) Node {
	n := Node{}
	n.key = key
	return n
}

// SetNthChild sets the Nth child. If an existing node is replaced,
// that node is returned.
func (k *Node) SetNthChild(n uint, other *Node) *Node {
	other.n = n

	if k.firstChild == nil {
		k.firstChild = other
		return nil
	}

	if k.firstChild.n == n {
		// evict
		ret := k.firstChild
		other.nextSibling = k.firstChild.nextSibling
		k.firstChild = other
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
func (k *Node) NthChild(n uint) *Node {
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

// BFS is a channel-based BFS for tree nodes.
// Channels are used similar to Python generators.
// Inspired by https://blog.carlmjohnson.net/post/on-using-go-channels-like-python-generators/
// Examples of how to use it can be seen in algorithms_test.go
func BFS(root *Node, quit <-chan struct{}) <-chan *Node {
	nChan := make(chan *Node)

	go func() {
		defer close(nChan)
		queue := []*Node{root}
		var curr *Node

		for len(queue) > 0 {
			curr, queue = queue[0], queue[1:]

			select {
			case <-quit:
				return
			case nChan <- curr:
			}

			next := curr.firstChild
			for next != nil {
				queue = append(queue, next)
				next = next.nextSibling
			}
		}
	}()

	return nChan
}

// Equals does a deep comparison of two tree nodes. The only special
// behavior is that two nils are considered "equal trees."
func Equals(a, b *Node) bool {
	if a == b {
		return true
	}

	if a == nil {
		return false
	}

	if b == nil {
		return false
	}

	if a.n != b.n || a.Key() != b.Key() {
		return false
	}

	nextA := a.firstChild
	nextB := b.firstChild

	if (nextA != nil && nextB != nil) && nextA.n != nextB.n {
		return false
	}

	for {
		if !Equals(nextA, nextB) {
			return false
		}

		if nextA == nil && nextB == nil {
			return true
		}

		nextA = nextA.nextSibling
		nextB = nextB.nextSibling
	}
}
