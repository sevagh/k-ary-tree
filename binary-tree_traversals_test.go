package karytree_test

import (
	"testing"

	"github.com/sevagh/k-ary-tree"
)

func TestBFSBinary(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	d := karytree.Binary("d")

	a.SetLeft(&b)
	b.SetRight(&c)
	c.SetRight(&d)

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

func TestBFSBinaryEarlyQuit(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	d := karytree.Binary("d")

	a.SetLeft(&b)
	b.SetLeft(&c)
	c.SetRight(&d)

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

func TestInorderIterative1(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	e := karytree.Binary("e")
	f := karytree.Binary("f")
	g := karytree.Binary("g")
	h := karytree.Binary("h")

	a.SetLeft(&b)
	a.SetRight(&c)

	b.SetLeft(&e)
	b.SetRight(&f)

	c.SetLeft(&g)
	c.SetRight(&h)

	/*
		       a
			 /   \
			b     c
		   / \   / \
		  e   f g   h
	*/

	ctr := 0
	for node := range karytree.InorderIterative(&a, nil) {
		nodeKey := node.Key().(string)
		if ctr == 0 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		} else if ctr == 1 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if ctr == 2 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if ctr == 3 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if ctr == 4 && nodeKey != "g" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if ctr == 5 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if ctr == 6 && nodeKey != "h" {
			t.Errorf("expected node key 'h', got '%s'", nodeKey)
		} else if ctr >= 7 {
			t.Errorf("the inorder traversal should've ended")
		}
		ctr++
	}
}

func TestInorderIterative2(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	d := karytree.Binary("d")
	e := karytree.Binary("e")
	f := karytree.Binary("f")

	a.SetLeft(&b)
	b.SetRight(&c)
	b.SetLeft(&f)
	c.SetLeft(&d)
	c.SetRight(&e)

	/*
		      a
			 /
			b
		   / \
		  f   c
		 	 / \
			d   e
	*/

	ctr := 0
	for node := range karytree.InorderIterative(&a, nil) {
		nodeKey := node.Key().(string)
		if ctr == 0 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if ctr == 1 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if ctr == 2 && nodeKey != "d" {
			t.Errorf("expected node key 'd', got '%s'", nodeKey)
		} else if ctr == 3 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if ctr == 4 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		} else if ctr == 5 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if ctr >= 6 {
			t.Errorf("the inorder traversal should've ended")
		}
		ctr++
	}
}

func TestInorderIterativeEarlyQuit(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	d := karytree.Binary("d")

	a.SetLeft(&b)
	b.SetLeft(&c)
	c.SetRight(&d)

	//traverse the linkedlist
	quit := make(chan struct{})

	ctr := 0
	for node := range karytree.InorderIterative(&a, quit) {
		nodeKey := node.Key().(string)
		if ctr == 0 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if ctr >= 1 {
			t.Errorf("expected early quit, still getting values on inorder chan")
		}

		t.Logf("early quit inorder")
		quit <- struct{}{}
		ctr++
	}
}

func TestPreorderIterative1(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	e := karytree.Binary("e")
	f := karytree.Binary("f")
	g := karytree.Binary("g")
	h := karytree.Binary("h")

	a.SetLeft(&b)
	a.SetRight(&c)

	b.SetLeft(&e)
	b.SetRight(&f)

	c.SetLeft(&g)
	c.SetRight(&h)

	/*
		       a
			 /   \
			b     c
		   / \   / \
		  e   f g   h
	*/

	ctr := 0
	for node := range karytree.PreorderIterative(&a, nil) {
		nodeKey := node.Key().(string)
		if ctr == 0 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if ctr == 1 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if ctr == 2 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		} else if ctr == 3 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if ctr == 4 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if ctr == 5 && nodeKey != "g" {
			t.Errorf("expected node key 'g', got '%s'", nodeKey)
		} else if ctr == 6 && nodeKey != "h" {
			t.Errorf("expected node key 'h', got '%s'", nodeKey)
		} else if ctr >= 7 {
			t.Errorf("the preorder traversal should've ended")
		}
		ctr++
	}
}

func TestPreorderIterative2(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	d := karytree.Binary("d")
	e := karytree.Binary("e")
	f := karytree.Binary("f")

	a.SetLeft(&b)
	b.SetRight(&c)
	b.SetLeft(&f)
	c.SetLeft(&d)
	c.SetRight(&e)

	/*
		      a
			 /
			b
		   / \
		  f   c
		 	 / \
			d   e
	*/

	ctr := 0
	for node := range karytree.PreorderIterative(&a, nil) {
		nodeKey := node.Key().(string)
		if ctr == 0 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if ctr == 1 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if ctr == 2 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if ctr == 3 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if ctr == 4 && nodeKey != "d" {
			t.Errorf("expected node key 'd', got '%s'", nodeKey)
		} else if ctr == 5 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		} else if ctr >= 6 {
			t.Errorf("the preorder traversal should've ended")
		}
		ctr++
	}
}

func TestPreorderIterativeEarlyQuit(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	d := karytree.Binary("d")

	a.SetLeft(&b)
	b.SetLeft(&c)
	c.SetRight(&d)

	//traverse the linkedlist
	quit := make(chan struct{})

	ctr := 0
	for node := range karytree.PreorderIterative(&a, quit) {
		nodeKey := node.Key().(string)
		if ctr == 0 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if ctr >= 1 {
			t.Errorf("expected early quit, still getting values on preorder chan")
		}

		t.Logf("early quit preorder")
		quit <- struct{}{}
		ctr++
	}
}

func TestPostorderIterative1(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	e := karytree.Binary("e")
	f := karytree.Binary("f")
	g := karytree.Binary("g")
	h := karytree.Binary("h")

	a.SetLeft(&b)
	a.SetRight(&c)

	b.SetLeft(&e)
	b.SetRight(&f)

	c.SetLeft(&g)
	c.SetRight(&h)

	/*
		       a
			 /   \
			b     c
		   / \   / \
		  e   f g   h
	*/

	ctr := 0
	for node := range karytree.PostorderIterative(&a, nil) {
		nodeKey := node.Key().(string)
		t.Logf("got node: %+v\n", nodeKey)
		if ctr == 0 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		} else if ctr == 1 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if ctr == 2 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if ctr == 3 && nodeKey != "g" {
			t.Errorf("expected node key 'g', got '%s'", nodeKey)
		} else if ctr == 4 && nodeKey != "h" {
			t.Errorf("expected node key 'h', got '%s'", nodeKey)
		} else if ctr == 5 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if ctr == 6 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if ctr >= 7 {
			t.Errorf("the postorder traversal should've ended")
		}
		ctr++
	}

	if ctr != 7 {
		t.Errorf("didn't even get any nodes")
	}
}

func TestPostorderIterative2(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	d := karytree.Binary("d")
	e := karytree.Binary("e")
	f := karytree.Binary("f")

	a.SetLeft(&b)
	b.SetRight(&c)
	b.SetLeft(&f)
	c.SetLeft(&d)
	c.SetRight(&e)

	/*
		      a
			 /
			b
		   / \
		  f   c
		 	 / \
			d   e
	*/

	ctr := 0
	for node := range karytree.PostorderIterative(&a, nil) {
		nodeKey := node.Key().(string)
		if ctr == 0 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if ctr == 1 && nodeKey != "d" {
			t.Errorf("expected node key 'd', got '%s'", nodeKey)
		} else if ctr == 2 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		} else if ctr == 3 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if ctr == 4 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if ctr == 5 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if ctr >= 6 {
			t.Errorf("the postorder traversal should've ended")
		}
		ctr++
	}

	if ctr != 6 {
		t.Errorf("didn't even get any nodes")
	}
}

func TestPostorderIterativeEarlyQuit(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	d := karytree.Binary("d")

	a.SetLeft(&b)
	b.SetLeft(&c)
	c.SetRight(&d)

	//traverse the linkedlist
	quit := make(chan struct{})

	ctr := 0
	for node := range karytree.PostorderIterative(&a, quit) {
		nodeKey := node.Key().(string)
		if ctr == 0 && nodeKey != "d" {
			t.Errorf("expected node key 'd', got '%s'", nodeKey)
		} else if ctr >= 1 {
			t.Errorf("expected early quit, still getting values on postorder chan")
		}

		t.Logf("early quit postorder")
		quit <- struct{}{}
		ctr++
	}
}

func TestPostorderRecursive1(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	e := karytree.Binary("e")
	f := karytree.Binary("f")
	g := karytree.Binary("g")
	h := karytree.Binary("h")

	a.SetLeft(&b)
	a.SetRight(&c)

	b.SetLeft(&e)
	b.SetRight(&f)

	c.SetLeft(&g)
	c.SetRight(&h)

	/*
		       a
			 /   \
			b     c
		   / \   / \
		  e   f g   h
	*/

	nodes := []*karytree.Node{}
	visit := func(node *karytree.Node) {
		nodes = append(nodes, node)
	}

	karytree.PostorderRecursive(&a, visit)
	i := 0
	var node *karytree.Node

	for i, node = range nodes {
		nodeKey := node.Key().(string)
		if i == 0 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		} else if i == 1 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if i == 2 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if i == 3 && nodeKey != "g" {
			t.Errorf("expected node key 'g', got '%s'", nodeKey)
		} else if i == 4 && nodeKey != "h" {
			t.Errorf("expected node key 'h', got '%s'", nodeKey)
		} else if i == 5 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if i == 6 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		}
	}

	if i != 6 {
		t.Errorf("got wrong amount of nodes: %+v, expected 6", i)
	}
}

func TestPostorderRecursive2(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	d := karytree.Binary("d")
	e := karytree.Binary("e")
	f := karytree.Binary("f")

	a.SetLeft(&b)
	b.SetRight(&c)
	b.SetLeft(&f)
	c.SetLeft(&d)
	c.SetRight(&e)

	/*
		      a
			 /
			b
		   / \
		  f   c
		 	 / \
			d   e
	*/
	nodes := []*karytree.Node{}
	visit := func(node *karytree.Node) {
		nodes = append(nodes, node)
	}

	karytree.PostorderRecursive(&a, visit)
	i := 0
	var node *karytree.Node
	for i, node = range nodes {
		nodeKey := node.Key().(string)
		if i == 0 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if i == 1 && nodeKey != "d" {
			t.Errorf("expected node key 'd', got '%s'", nodeKey)
		} else if i == 2 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		} else if i == 3 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if i == 4 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if i == 5 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if i >= 6 {
			t.Errorf("the postorder traversal should've ended")
		}
		i++
	}

	if i != 6 {
		t.Errorf("got wrong amount of nodes: %+v, expected 6", i)
	}
}

func TestPreorderRecursive1(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	e := karytree.Binary("e")
	f := karytree.Binary("f")
	g := karytree.Binary("g")
	h := karytree.Binary("h")

	a.SetLeft(&b)
	a.SetRight(&c)

	b.SetLeft(&e)
	b.SetRight(&f)

	c.SetLeft(&g)
	c.SetRight(&h)

	/*
		       a
			 /   \
			b     c
		   / \   / \
		  e   f g   h
	*/

	nodes := []*karytree.Node{}
	visit := func(node *karytree.Node) {
		nodes = append(nodes, node)
	}

	karytree.PreorderRecursive(&a, visit)
	i := 0
	var node *karytree.Node
	for i, node = range nodes {
		nodeKey := node.Key().(string)
		if i == 0 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if i == 1 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if i == 2 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		} else if i == 3 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if i == 4 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if i == 5 && nodeKey != "g" {
			t.Errorf("expected node key 'g', got '%s'", nodeKey)
		} else if i == 6 && nodeKey != "h" {
			t.Errorf("expected node key 'h', got '%s'", nodeKey)
		}
	}

	if i != 6 {
		t.Errorf("got wrong amount of nodes: %+v, expected 6", i)
	}
}

func TestPreorderRecursive2(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	d := karytree.Binary("d")
	e := karytree.Binary("e")
	f := karytree.Binary("f")

	a.SetLeft(&b)
	b.SetRight(&c)
	b.SetLeft(&f)
	c.SetLeft(&d)
	c.SetRight(&e)

	/*
		      a
			 /
			b
		   / \
		  f   c
		 	 / \
			d   e
	*/

	nodes := []*karytree.Node{}
	visit := func(node *karytree.Node) {
		nodes = append(nodes, node)
	}

	karytree.PreorderRecursive(&a, visit)
	i := 0
	var node *karytree.Node
	for i, node = range nodes {
		nodeKey := node.Key().(string)
		if i == 0 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if i == 1 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if i == 2 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if i == 3 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if i == 4 && nodeKey != "d" {
			t.Errorf("expected node key 'd', got '%s'", nodeKey)
		} else if i == 5 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		}
	}

	if i != 5 {
		t.Errorf("got wrong amount of nodes: %+v, expected 5", i)
	}
}

func TestInorderRecursive1(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	e := karytree.Binary("e")
	f := karytree.Binary("f")
	g := karytree.Binary("g")
	h := karytree.Binary("h")

	a.SetLeft(&b)
	a.SetRight(&c)

	b.SetLeft(&e)
	b.SetRight(&f)

	c.SetLeft(&g)
	c.SetRight(&h)

	/*
		       a
			 /   \
			b     c
		   / \   / \
		  e   f g   h
	*/

	nodes := []*karytree.Node{}
	visit := func(node *karytree.Node) {
		nodes = append(nodes, node)
	}

	karytree.InorderRecursive(&a, visit)
	i := 0
	var node *karytree.Node
	for i, node = range nodes {
		nodeKey := node.Key().(string)
		if i == 0 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		} else if i == 1 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if i == 2 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if i == 3 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		} else if i == 4 && nodeKey != "g" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if i == 5 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if i == 6 && nodeKey != "h" {
			t.Errorf("expected node key 'h', got '%s'", nodeKey)
		} else if i >= 7 {
			t.Errorf("the inorder traversal should've ended")
		}
	}

	if i != 6 {
		t.Errorf("got wrong amount of nodes: %+v, expected 6", i)
	}
}

func TestInorderRecursive2(t *testing.T) {
	a := karytree.Binary("a")
	b := karytree.Binary("b")
	c := karytree.Binary("c")
	d := karytree.Binary("d")
	e := karytree.Binary("e")
	f := karytree.Binary("f")

	a.SetLeft(&b)
	b.SetRight(&c)
	b.SetLeft(&f)
	c.SetLeft(&d)
	c.SetRight(&e)

	/*
		      a
			 /
			b
		   / \
		  f   c
		 	 / \
			d   e
	*/

	nodes := []*karytree.Node{}
	visit := func(node *karytree.Node) {
		nodes = append(nodes, node)
	}

	karytree.InorderRecursive(&a, visit)
	i := 0
	var node *karytree.Node
	for i, node = range nodes {
		nodeKey := node.Key().(string)
		if i == 0 && nodeKey != "f" {
			t.Errorf("expected node key 'f', got '%s'", nodeKey)
		} else if i == 1 && nodeKey != "b" {
			t.Errorf("expected node key 'b', got '%s'", nodeKey)
		} else if i == 2 && nodeKey != "d" {
			t.Errorf("expected node key 'd', got '%s'", nodeKey)
		} else if i == 3 && nodeKey != "c" {
			t.Errorf("expected node key 'c', got '%s'", nodeKey)
		} else if i == 4 && nodeKey != "e" {
			t.Errorf("expected node key 'e', got '%s'", nodeKey)
		} else if i == 5 && nodeKey != "a" {
			t.Errorf("expected node key 'a', got '%s'", nodeKey)
		}
	}

	if i != 5 {
		t.Errorf("got wrong amount of nodes: %+v, expected 5", i)
	}
}
