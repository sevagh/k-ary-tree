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
