package karytree_test

import (
	"testing"

	"github.com/sevagh/k-ary-tree"
)

func TestBFS(t *testing.T) {
	//k = 1 == we basically have a linked list

	a := karytree.New(1, "a")
	b := karytree.New(1, "b")
	c := karytree.New(1, "c")
	d := karytree.New(1, "d")

	a.SetNthChild(0, &b)
	b.SetNthChild(0, &c)
	c.SetNthChild(0, &d)

	//traverse the linkedlist

	bfsCtr := 0
	for node := range karytree.BFS(&a, nil) {
		nodeKey := node.Key().(string)
		if bfsCtr == 0 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if bfsCtr == 1 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if bfsCtr == 2 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if bfsCtr == 3 && nodeKey != "d" {
			t.Errorf("expected node key 'd', got '%s'", nodeKey)
		} else if bfsCtr >= 4 {
			t.Errorf("the bfs traversal should've ended")
		}
		bfsCtr++
	}
}

func TestBFSEarlyQuit(t *testing.T) {
	//k = 1 == we basically have a linked list

	a := karytree.New(1, "a")
	b := karytree.New(1, "b")
	c := karytree.New(1, "c")
	d := karytree.New(1, "d")

	a.SetNthChild(0, &b)
	b.SetNthChild(0, &c)
	c.SetNthChild(0, &d)

	//traverse the linkedlist
	quit := make(chan struct{})

	bfsCtr := 0
	for node := range karytree.BFS(&a, quit) {
		nodeKey := node.Key().(string)
		if bfsCtr == 0 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if bfsCtr >= 1 {
			t.Errorf("expected early quit, still getting values on bfs chan")
		}

		t.Logf("early quit bfs")
		quit <- struct{}{}
		bfsCtr++
	}
}

func TestCompareTrees(t *testing.T) {
	tree1 := constructTree(8)
	tree2 := constructTree(8)

	if !karytree.Equals(&tree1, &tree2) {
		t.Errorf("expected identical trees to be equal")
	}
}

func TestCompareSparseTrees(t *testing.T) {
	tree1 := constructTreeSparse(8)
	tree2 := constructTreeSparse(8)

	if !karytree.Equals(&tree1, &tree2) {
		t.Errorf("expected identical trees to be equal")
	}
}

func TestEmptyTreesEqual(t *testing.T) {
	if !karytree.Equals(nil, nil) {
		t.Errorf("two nil nodes are be equal")
	}
}

func TestOneNilEqual(t *testing.T) {
	tree1 := constructTreeSparse(8)

	if karytree.Equals(nil, &tree1) || karytree.Equals(&tree1, nil) {
		t.Errorf("nil and real trees can't be equal")
	}
}

func TestNotEqualTrees(t *testing.T) {
	tree1 := constructTree(8)
	tree2 := constructTree(8)

	x := tree2.NthChild(3)
	rand := karytree.New(x.K(), "hello world")
	tree2.SetNthChild(3, &rand)

	if karytree.Equals(&tree1, &tree2) {
		t.Errorf("tree1 and tree2 shouldn't be equal")
	}
}

func constructTree(K int) karytree.Node {
	var key int
	tree := karytree.New(K, key)
	key++

	var curr *karytree.Node
	curr = &tree

	for k := 0; k < K; k++ {
		child := karytree.New(K, key)
		key++
		curr.SetNthChild(k, &child)
		for l := 0; l < K; l++ {
			grandchild := karytree.New(K, key)
			key++
			nth := curr.NthChild(k)
			nth.SetNthChild(l, &grandchild)
			for m := 0; m < K; m++ {
				greatgrandchild := karytree.New(K, key)
				key++

				grandnth := nth.NthChild(l)
				grandnth.SetNthChild(m, &greatgrandchild)
			}
		}
	}

	return tree
}

func constructTreeSparse(K int) karytree.Node {
	var tree karytree.Node

	var key int

	tree = karytree.New(K, key)
	key++

	var curr *karytree.Node
	curr = &tree

	for i := 0; i < K; i++ {
		if i%2 == 0 {
			child := karytree.New(K, key)
			key++

			// fill even children
			curr.SetNthChild(i, &child)
			for j := 0; j < K; j++ {
				if j%2 != 0 {
					grandchild := karytree.New(K, key)
					key++
					ith := curr.NthChild(i)

					// fill odd grandchildren
					ith.SetNthChild(j, &grandchild)
					for k := 0; k < K; k++ {
						if k%2 == 0 {
							// fill even great grandchildren
							greatgrandchild := karytree.New(K, key)
							key++
							jth := ith.NthChild(j)

							jth.SetNthChild(k, &greatgrandchild)
						}
					}
				}
			}
		}
	}

	return tree
}
