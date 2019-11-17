package karytree_test

import (
	"fmt"

	"github.com/sevagh/k-ary-tree"
)

func ExampleNode() {
	key := 0
	tree := karytree.NewNode(key)
	key++

	for i := uint16(0); i < uint16(16); i++ {
		newNode := karytree.NewNode(key)
		key++
		tree.SetNthChild(i, &newNode)
	}

	for node := range karytree.BFS(&tree, nil) {
		fmt.Printf("%d ", node.Key().(int))
	}

	// Output: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16
}
