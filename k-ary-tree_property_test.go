package karytree_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/flyingmutant/rapid"
	"github.com/google/gofuzz"
	"github.com/sevagh/k-ary-tree"
)

type karytreeMachine struct {
	r     karytree.Node
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
	m.r = karytree.NewNode(getKFuzzedKey())
	t.Logf("Created k-ary-tree root node\n")
}

func (m *karytreeMachine) Get(t *rapid.T) {
	if len(m.path) == 0 {
		t.Skip("tree probably empty")
	}

	currPath := m.path[0]
	currState := m.state[0]

	t.Logf("path is: %+v\n", currPath)

	var curr *karytree.Node
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

	var curr *karytree.Node
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
			newNode := karytree.NewNode(newFuzzedKey)
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
