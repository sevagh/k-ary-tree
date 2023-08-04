package karytree_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/flyingmutant/rapid"
	"github.com/google/gofuzz"
	"github.com/sevagh/k-ary-tree"
)

func TestBasicLinkedList(t *testing.T) {
	//k = 1 == we basically have a linked list

	a := karytree.NewNode[interface{}]("a")
	b := karytree.NewNode[interface{}]("b")
	c := karytree.NewNode[interface{}]("c")
	d := karytree.NewNode[interface{}]("d")

	a.SetNthChild(0, &b)
	b.SetNthChild(0, &c)
	c.SetNthChild(0, &d)

	//traverse the linkedlist

	aKey := a.Key().(string)
	if aKey != "a" {
		t.Errorf("expected \"a\", got %+v\n", aKey)
	}

	aNext := a.NthChild(0)
	if aNext != &b {
		t.Errorf("expected a's next node to be b, got: %+v", aNext)
	}
	bKey := aNext.Key().(string)
	if bKey != "b" {
		t.Errorf("expected \"b\", got %+v\n", bKey)
	}

	bNext := b.NthChild(0)
	if bNext != &c {
		t.Errorf("expected b's next node to be c, got: %+v", bNext)
	}
	cKey := bNext.Key().(string)
	if cKey != "c" {
		t.Errorf("expected \"c\", got %+v\n", cKey)
	}

	cNext := c.NthChild(0)
	if cNext != &d {
		t.Errorf("expected c's next node to be d, got: %+v", cNext)
	}
	dKey := cNext.Key().(string)
	if dKey != "d" {
		t.Errorf("expected \"d\", got %+v\n", cKey)
	}

	dNext := d.NthChild(0)
	if dNext != nil {
		t.Errorf("expected d's next node to be nil, got: %+v", dNext)
	}
}

func TestModifyKey(t *testing.T) {
	a := karytree.NewNode[interface{}]("a")
	if a.Key().(string) != "a" {
		t.Errorf("key was 'a', should not be %+v\n", a.Key().(string))
	}

	a.SetKey("b")
	if a.Key().(string) != "b" {
		t.Errorf("key was changed to 'b', should not be %+v\n", a.Key().(string))
	}
}

func TestSiblingTreeNLogic(t *testing.T) {
	a := karytree.NewNode[interface{}]("a")
	b := karytree.NewNode[interface{}]("b")
	c := karytree.NewNode[interface{}]("c")
	d := karytree.NewNode[interface{}]("d")
	e := karytree.NewNode[interface{}]("e")

	a.SetNthChild(32, &b)
	a.SetNthChild(5, &c)
	c.SetNthChild(0, &d)
	c.SetNthChild(1, &e)

	if a.NthChild(32).Key().(string) != "b" {
		t.Errorf("didn't set this child correctly")
	}
	if a.NthChild(5).Key().(string) != "c" {
		t.Errorf("didn't set this child correctly")
	}
	if c.NthChild(0).Key().(string) != "d" {
		t.Errorf("didn't set this child correctly")
	}
	if c.NthChild(1).Key().(string) != "e" {
		t.Errorf("didn't set this child correctly")
	}
}

func TestSetSameChildEvictsFormer(t *testing.T) {
	a := karytree.NewNode[interface{}]("a")
	b := karytree.NewNode[interface{}]("b")
	c := karytree.NewNode[interface{}]("c")
	d := karytree.NewNode[interface{}]("d")

	a.SetNthChild(4, &b)
	a.SetNthChild(5, &d)
	evicted := a.SetNthChild(4, &c)

	if a.NthChild(4).Key().(string) != "c" {
		t.Errorf("expected new child to be 'c', got '%+v'", a.NthChild(4).Key())
	}

	if evicted.Key().(string) != "b" {
		t.Errorf("expected evicted child to be 'b', got '%+v'", evicted.Key())
	}

	if a.NthChild(5) != &d {
		t.Errorf("evicted node's siblings weren't inherited")
	}

	e := karytree.NewNode[interface{}]("e")
	f := karytree.NewNode[interface{}]("f")

	a.SetNthChild(6, &e)
	evicted = a.SetNthChild(6, &f)
}

type karytreeMachine struct {
	r     karytree.Node[interface{}]
	path  [][]uint
	state []interface{}
}

func getKFuzzedKey() interface{} {
	f := fuzz.New().NilChance(0) // we can't use nils
	// my library uses nil interfaces as sentinels

	var ret interface{}
	var n int

	switch n = rand.Intn(9); n {
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

func (m *karytreeMachine) Init(t *rapid.T) {
	m.r = karytree.NewNode[interface{}](getKFuzzedKey())
	t.Logf("Created k-ary-tree root node\n")
}

func (m *karytreeMachine) Get(t *rapid.T) {
	if len(m.path) == 0 {
		t.Skip("tree probably empty")
	}

	currPath := m.path[0]
	currState := m.state[0]

	t.Logf("path is: %+v\n", currPath)

	var curr *karytree.Node[interface{}]
	curr = &m.r
	for _, p := range currPath {
		curr = curr.NthChild(p)
		t.Logf("descending through %dth child with key %+v\n", p, curr.Key())
	}

	if !reflect.DeepEqual(curr.Key(), currState) {
		t.Fatalf("got invalid value: %v vs expected %v", curr.Key(), currState)
	}

	m.state = m.state[1:]
	m.path = m.path[1:]
}

func (m *karytreeMachine) Put(t *rapid.T) {
	// can't set nth child > k for a k-ary tree
	path := rapid.SlicesOf(rapid.UintsRange(0, ^uint(0))).Draw(t, "nthChild").([]uint)

	var curr *karytree.Node[interface{}]
	var lastFuzzedKey interface{}
	curr = &m.r
	lastFuzzedKey = curr.Key()
	for _, p := range path {
		existingNthChild := curr.NthChild(p)
		if existingNthChild != nil {
			// going through a path that already exists
			curr = existingNthChild
			lastFuzzedKey = curr.Key()
		} else {
			newFuzzedKey := getKFuzzedKey()
			newNode := karytree.NewNode[interface{}](newFuzzedKey)
			curr.SetNthChild(p, &newNode)
			curr = &newNode
			lastFuzzedKey = newFuzzedKey
		}
	}

	m.state = append([]interface{}{lastFuzzedKey}, m.state...)
	m.path = append([][]uint{path}, m.path...)

	t.Logf("paths %+v\n", m.path)
	t.Logf("state %+v\n", m.state)
}

func TestKarytreePropertyFuzz(t *testing.T) {
	rapid.Check(t, rapid.StateMachine(&karytreeMachine{}))
}

func TestBFS(t *testing.T) {
	//k = 1 == we basically have a linked list

	a := karytree.NewNode[interface{}]("a")
	b := karytree.NewNode[interface{}]("b")
	c := karytree.NewNode[interface{}]("c")
	d := karytree.NewNode[interface{}]("d")

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

	a := karytree.NewNode[interface{}]("a")
	b := karytree.NewNode[interface{}]("b")
	c := karytree.NewNode[interface{}]("c")
	d := karytree.NewNode[interface{}]("d")

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
	if !karytree.Equals[interface{}](nil, nil) {
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

	rand := karytree.NewNode[interface{}]("hello world")
	tree2.SetNthChild(3, &rand)

	if karytree.Equals(&tree1, &tree2) {
		t.Errorf("tree1 and tree2 shouldn't be equal")
	}
}

func TestTreeInsertionSortedOrderEquals(t *testing.T) {
	a := karytree.NewNode[interface{}]("a")
	b := karytree.NewNode[interface{}]("b")
	c := karytree.NewNode[interface{}]("c")

	a_ := karytree.NewNode[interface{}]("a")
	b_ := karytree.NewNode[interface{}]("b")
	c_ := karytree.NewNode[interface{}]("c")

	a.SetNthChild(4, &b)
	a.SetNthChild(2, &c)

	a_.SetNthChild(2, &c_)
	a_.SetNthChild(4, &b_)

	if !karytree.Equals(&a, &a_) {
		t.Errorf("sibling list sorted order is not working right")
	}
}

func TestTreeInsertionSortedOrderNotEquals(t *testing.T) {
	a := karytree.NewNode[interface{}]("a")
	b := karytree.NewNode[interface{}]("b")
	c := karytree.NewNode[interface{}]("c")

	a_ := karytree.NewNode[interface{}]("a")
	b_ := karytree.NewNode[interface{}]("b")
	c_ := karytree.NewNode[interface{}]("c")

	a.SetNthChild(4, &b)
	a.SetNthChild(5, &c)

	a_.SetNthChild(4, &b_)
	a_.SetNthChild(9, &c_)

	if karytree.Equals(&a, &a_) {
		t.Errorf("sibling list sorted order is not working right")
	}
}

func constructTree(K int) karytree.Node[interface{}] {
	var key int
	tree := karytree.NewNode[interface{}](key)
	key++

	var curr *karytree.Node[interface{}]
	curr = &tree

	for k := uint(0); k < uint(K); k++ {
		child := karytree.NewNode[interface{}](key)
		key++
		curr.SetNthChild(k, &child)
		for l := uint(0); l < uint(K); l++ {
			grandchild := karytree.NewNode[interface{}](key)
			key++
			nth := curr.NthChild(k)
			nth.SetNthChild(l, &grandchild)
			for m := uint(0); m < uint(K); m++ {
				greatgrandchild := karytree.NewNode[interface{}](key)
				key++

				grandnth := nth.NthChild(l)
				grandnth.SetNthChild(m, &greatgrandchild)
			}
		}
	}

	return tree
}

func constructTreeSparse(K int) karytree.Node[interface{}] {
	var tree karytree.Node[interface{}]

	var key int

	tree = karytree.NewNode[interface{}](key)
	key++

	var curr *karytree.Node[interface{}]
	curr = &tree

	for i := uint(0); i < uint(K); i++ {
		if i%2 == 0 {
			child := karytree.NewNode[interface{}](key)
			key++

			// fill even children
			curr.SetNthChild(i, &child)
			for j := uint(0); j < uint(K); j++ {
				if j%2 != 0 {
					grandchild := karytree.NewNode[interface{}](key)
					key++
					ith := curr.NthChild(i)

					// fill odd grandchildren
					ith.SetNthChild(j, &grandchild)
					for k := uint(0); k < uint(K); k++ {
						if k%2 == 0 {
							// fill even great grandchildren
							greatgrandchild := karytree.NewNode[interface{}](key)
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
