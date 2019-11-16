package karytree_test

import (
	"testing"

	"github.com/sevagh/k-ary-tree"
)

func TestBasicLinkedList(t *testing.T) {
	//k = 1 == we basically have a linked list

	a := karytree.New(1, "a")
	b := karytree.New(1, "b")
	c := karytree.New(1, "c")
	d := karytree.New(1, "d")

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
	a := karytree.New(1, "a")
	if a.Key().(string) != "a" {
		t.Errorf("key was 'a', should not be %+v\n", a.Key().(string))
	}

	a.SetKey("b")
	if a.Key().(string) != "b" {
		t.Errorf("key was changed to 'b', should not be %+v\n", a.Key().(string))
	}
}
