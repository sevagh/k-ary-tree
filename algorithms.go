package karytree

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

			next := curr.getFirstChild()
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

	if a.n() != b.n() || a.Key() != b.Key() {
		return false
	}

	nextA := a.getFirstChild()
	nextB := b.getFirstChild()

	if (nextA != nil && nextB != nil) && nextA.n() != nextB.n() {
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
