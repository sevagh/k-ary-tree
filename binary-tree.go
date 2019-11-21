package karytree

const (
	left  = iota
	right = iota
)

// Binary creates a binary karytree.Node
func Binary(key interface{}) Node {
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

// InorderIterative is a channel-based iterative implementation of an preorder traversal.
func InorderIterative(root *Node, quit <-chan struct{}) <-chan *Node {
	nChan := make(chan *Node)

	go func() {
		defer close(nChan)
		stack := []*Node{}
		var curr *Node
		curr = root

		for {
			for curr != nil {
				stack = append(stack, curr)
				curr = curr.Left()
			}

			if len(stack) == 0 {
				break
			}

			stack, curr = stack[:len(stack)-1], stack[len(stack)-1]

			select {
			case <-quit:
				return
			case nChan <- curr:
			}

			curr = curr.Right()
		}
	}()

	return nChan
}

// PreorderIterative is a channel-based iterative implementation of an preorder traversal.
func PreorderIterative(root *Node, quit <-chan struct{}) <-chan *Node {
	nChan := make(chan *Node)

	go func() {
		defer close(nChan)
		stack := []*Node{}
		var curr *Node
		curr = root

		for {
			select {
			case <-quit:
				return
			case nChan <- curr:
			}

			left := curr.Left()
			if left != nil {
				right := curr.Right()
				if right != nil {
					stack = append(stack, right)
				}

				curr = left
				continue
			}

			if len(stack) == 0 {
				break
			}

			stack, curr = stack[:len(stack)-1], stack[len(stack)-1]
			continue
		}
	}()

	return nChan
}

// PostorderIterative is a channel-based iterative implementation of an preorder traversal.
func PostorderIterative(root *Node, quit <-chan struct{}) <-chan *Node {
	nChan := make(chan *Node)

	go func() {
		defer close(nChan)
		stack1 := []*Node{}
		stack2 := []*Node{}

		var curr *Node

		stack1 = append(stack1, root)

		for len(stack1) != 0 {
			stack1, curr = stack1[:len(stack1)-1], stack1[len(stack1)-1]
			stack2 = append(stack2, curr)

			left := curr.Left()
			if left != nil {
				stack1 = append(stack1, left)
			}
			right := curr.Right()
			if right != nil {
				stack1 = append(stack1, right)
			}
		}
		for len(stack2) != 0 {
			stack2, curr = stack2[:len(stack2)-1], stack2[len(stack2)-1]

			select {
			case <-quit:
				return
			case nChan <- curr:
			}
		}
	}()

	return nChan
}

// InorderRecursive is a recursive inorder traversal with visitors
func InorderRecursive(root *Node, f func(*Node)) {
	inorder(root, f)
}

func inorder(root *Node, f func(*Node)) {
	if root != nil {
		inorder(root.Left(), f)
		f(root)
		inorder(root.Right(), f)
	}
}

// PreorderRecursive is a recursive inorder traversal with visitors
func PreorderRecursive(root *Node, f func(*Node)) {
	preorder(root, f)
}

func preorder(root *Node, f func(*Node)) {
	if root != nil {
		f(root)
		preorder(root.Left(), f)
		preorder(root.Right(), f)
	}
}

// PostorderRecursive is a recursive inorder traversal with visitors
func PostorderRecursive(root *Node, f func(*Node)) {
	postorder(root, f)
}

func postorder(root *Node, f func(*Node)) {
	if root != nil {
		postorder(root.Left(), f)
		postorder(root.Right(), f)
		f(root)
	}
}
