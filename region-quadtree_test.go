package karytree_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/flyingmutant/rapid"
	"github.com/google/gofuzz"
	"github.com/sevagh/k-ary-tree"
)

type quadtreeMachine struct {
	r     karytree.Node
	state []interface{}
	path  [][]int
}

func getQFuzzedKey() interface{} {
	f := fuzz.New()
	var ret interface{}

	switch n := rand.Intn(10); n {
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

func (m *quadtreeMachine) Init(t *rapid.T) {
	m.r = karytree.RegionQuadtree(getQFuzzedKey())
	t.Logf("Created binary-tree node\n")
}

func (m *quadtreeMachine) Get(t *rapid.T) {
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
			next = curr.NE()
		} else if p == 1 {
			next = curr.NW()
		} else if p == 2 {
			next = curr.SW()
		} else {
			next = curr.SE()
		}
		curr = next
	}

	if !reflect.DeepEqual(curr.Key(), currState) {
		t.Fatalf("got invalid value: %v vs expected %v", curr.Key(), currState)
	}

	m.state = m.state[1:]
	m.path = m.path[1:]
}

func (m *quadtreeMachine) Put(t *rapid.T) {
	path := rapid.SlicesOf(rapid.IntsRange(0, 3)).Draw(t, "ne-nw-sw-se").([]int)

	var curr *karytree.Node
	var lastFuzzedKey interface{}
	curr = &m.r
	lastFuzzedKey = curr.Key()
	for _, p := range path {
		var existingChild *karytree.Node
		if p == 0 {
			existingChild = curr.NE()
		} else if p == 1 {
			existingChild = curr.NW()
		} else if p == 2 {
			existingChild = curr.SW()
		} else {
			existingChild = curr.SE()
		}
		if existingChild != nil {
			// going through a path that already exists
			curr = existingChild
			lastFuzzedKey = curr.Key()
		} else {
			newFuzzedKey := getQFuzzedKey()
			newNode := karytree.RegionQuadtree(newFuzzedKey)
			if p == 0 {
				curr.SetNE(&newNode)
			} else if p == 1 {
				curr.SetNW(&newNode)
			} else if p == 2 {
				curr.SetSW(&newNode)
			} else {
				curr.SetSE(&newNode)
			}
			curr = &newNode
			lastFuzzedKey = newFuzzedKey
		}
	}

	m.state = append([]interface{}{lastFuzzedKey}, m.state...)
	m.path = append([][]int{path}, m.path...)
}

func TestQuadtreeTreePropertyFuzz(t *testing.T) {
	rapid.Check(t, rapid.StateMachine(&quadtreeMachine{}))
}
