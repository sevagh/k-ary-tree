package karytree_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/flyingmutant/rapid"
	"github.com/google/gofuzz"
	"github.com/sevagh/k-ary-tree"
)

type binarytreeMachine struct {
	r     karytree.Node
	state []interface{}
	path  [][]int
}

func getBFuzzedKey() interface{} {
	f := fuzz.New().NilChance(0) // we can't use nils
	// my library uses nil interfaces as sentinels
	var ret interface{}

	switch n := rand.Intn(9); n {
	case 0:
		var key string
		f.Fuzz(&key)
		ret = key
	case 1:
		var key int
		f.Fuzz(&key)
		ret = key
	case 2:
		var key []string
		f.Fuzz(&key)
		ret = key
	case 3:
		var key []byte
		f.Fuzz(&key)
		ret = key
	case 4:
		var key map[string]string
		f.Fuzz(&key)
		ret = key
	case 5:
		var key bool
		f.Fuzz(&key)
		ret = key
	case 6:
		var key uintptr
		f.Fuzz(&key)
		ret = key
	case 7:
		var key []int64
		f.Fuzz(&key)
		ret = key
	case 8:
		var key map[int]uint32
		f.Fuzz(&key)
		ret = key
	}

	return ret
}

func (m *binarytreeMachine) Init(t *rapid.T) {
	m.r = karytree.Binary(getBFuzzedKey())
	t.Logf("Created binary-tree node\n")
}

func (m *binarytreeMachine) Get(t *rapid.T) {
	if len(m.path) == 0 {
		t.Skip("tree probably empty")
	}

	currPath := m.path[0]
	currState := m.state[0]

	t.Logf("path is: %+v\n", currPath)

	var curr *karytree.Node
	var next *karytree.Node
	curr = &m.r
	for _, p := range currPath {
		if p == 0 {
			next = curr.Left()
		} else {
			next = curr.Right()
		}
		curr = next
	}

	if !reflect.DeepEqual(curr.Key(), currState) {
		t.Fatalf("got invalid value: %v vs expected %v", curr.Key(), currState)
	}

	m.state = m.state[1:]
	m.path = m.path[1:]
}

func (m *binarytreeMachine) Put(t *rapid.T) {
	path := rapid.SlicesOf(rapid.IntsRange(0, 1)).Draw(t, "left-or-right").([]int)

	var curr *karytree.Node
	var lastFuzzedKey interface{}
	curr = &m.r
	lastFuzzedKey = curr.Key()
	for _, p := range path {
		var existingChild *karytree.Node
		if p == 0 {
			existingChild = curr.Left()
		} else {
			existingChild = curr.Right()
		}
		if existingChild != nil {
			// going through a path that already exists
			curr = existingChild
			lastFuzzedKey = curr.Key()
		} else {
			newFuzzedKey := getBFuzzedKey()
			newNode := karytree.Binary(newFuzzedKey)
			if p == 0 {
				curr.SetLeft(&newNode)
			} else {
				curr.SetRight(&newNode)
			}
			curr = &newNode
			lastFuzzedKey = newFuzzedKey
		}
	}

	m.state = append([]interface{}{lastFuzzedKey}, m.state...)
	m.path = append([][]int{path}, m.path...)
}

func TestBinaryTreePropertyFuzz(t *testing.T) {
	rapid.Check(t, rapid.StateMachine(&binarytreeMachine{}))
}

func TestTreeEquals(t *testing.T) {
	if !karytree.Equals(nil, nil) {
		t.Errorf("nil node should be nil node")
	}

	a := karytree.Binary("a")
	if karytree.Equals(&a, nil) || karytree.Equals(nil, &a) {
		t.Errorf("a is not nil")
	}

	b := karytree.Binary("b")
	if karytree.Equals(&a, &b) {
		t.Errorf("a is not b")
	}

	aPrime := karytree.Binary("a")

	if !karytree.Equals(&a, &aPrime) {
		t.Errorf("contents of a and aPrime are the same")
	}

	c := karytree.Binary("c")
	cc := karytree.Binary("cc")
	c.SetLeft(&cc)
	d := karytree.Binary("c")
	dd := karytree.Binary("dd")
	d.SetRight(&dd)

	if karytree.Equals(&c, &d) {
		t.Errorf("unequal trees problem")
	}

	if karytree.Equals(&c, &a) {
		t.Errorf("unequal trees problem")
	}
}

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
