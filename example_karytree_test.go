package karytree_test

import (
	"fmt"

	"github.com/sevagh/k-ary-tree"
)

func ExampleNode() {
	key := uint(0)
	tree := karytree.NewNode[uint](key)
	key++

	for i := uint(0); i < uint(16); i++ {
		newNode := karytree.NewNode(key)
		key++
		tree.SetNthChild(i, &newNode)
	}

	for node := range karytree.BFS(&tree, nil) {
		fmt.Printf("%d ", node.Key())
	}

	// Output: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16
}

func ExampleEquals() {
	a := karytree.NewNode("a")
	b := karytree.NewNode("b")
	c := karytree.NewNode("c")

	a_ := karytree.NewNode("a")
	b_ := karytree.NewNode("b")
	c_ := karytree.NewNode("c")

	a.SetNthChild(4, &b)
	// a -> firstChild -> (b, n: 4)

	a.SetNthChild(2, &c)
	// a -> firstChild -> (c, n: 2) -> nextSibling -> (b, n: 4)

	a_.SetNthChild(2, &c_)
	// a_ -> firstChild -> (c_, n: 2)

	a_.SetNthChild(4, &b_)
	// a_ -> firstChild -> (c_, n: 2) -> nextSibling -> (b_, n: 4)

	fmt.Println(karytree.Equals(&a, &a_))
	// Output: true
}
