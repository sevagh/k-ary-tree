package karytree

const (
	left  = iota
	right = iota
)

// Binary creates a binary karytree.Node
func Binary[T comparable](key T) Node[T] {
	return NewNode(key)
}

// SetLeft sets the left child.
func (k *Node[T]) SetLeft(other *Node[T]) {
	k.SetNthChild(left, other)
}

// SetRight sets the left child.
func (k *Node[T]) SetRight(other *Node[T]) {
	k.SetNthChild(right, other)
}

// Left gets the left child
func (k *Node[T]) Left() *Node[T] {
	return k.NthChild(left)
}

// Right gets the right child
func (k *Node[T]) Right() *Node[T] {
	return k.NthChild(right)
}

// InorderIterative is a channel-based iterative implementation of an preorder traversal.
func InorderIterative[T comparable](root *Node[T], quit <-chan struct{}) <-chan *Node[T] {
	nChan := make(chan *Node[T])

	go func() {
		defer close(nChan)
		stack := []*Node[T]{}
		var curr *Node[T]
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
func PreorderIterative[T comparable](root *Node[T], quit <-chan struct{}) <-chan *Node[T] {
	nChan := make(chan *Node[T])

	go func() {
		defer close(nChan)
		stack := []*Node[T]{}
		var curr *Node[T]
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
func PostorderIterative[T comparable](root *Node[T], quit <-chan struct{}) <-chan *Node[T] {
	nChan := make(chan *Node[T])

	go func() {
		defer close(nChan)
		stack1 := []*Node[T]{}
		stack2 := []*Node[T]{}

		var curr *Node[T]

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
func InorderRecursive[T comparable](root *Node[T], f func(*Node[T])) {
	inorder(root, f)
}

func inorder[T comparable](root *Node[T], f func(*Node[T])) {
	if root != nil {
		inorder(root.Left(), f)
		f(root)
		inorder(root.Right(), f)
	}
}

// PreorderRecursive is a recursive inorder traversal with visitors
func PreorderRecursive[T comparable](root *Node[T], f func(*Node[T])) {
	preorder(root, f)
}

func preorder[T comparable](root *Node[T], f func(*Node[T])) {
	if root != nil {
		f(root)
		preorder(root.Left(), f)
		preorder(root.Right(), f)
	}
}

// PostorderRecursive is a recursive inorder traversal with visitors
func PostorderRecursive[T comparable](root *Node[T], f func(*Node[T])) {
	postorder(root, f)
}

func postorder[T comparable](root *Node[T], f func(*Node[T])) {
	if root != nil {
		postorder(root.Left(), f)
		postorder(root.Right(), f)
		f(root)
	}
}
