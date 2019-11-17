package karytree_test

import (
	"testing"

	"github.com/sevagh/k-ary-tree"
)

func TestBasicLinkedList(t *testing.T) {
	//k = 1 == we basically have a linked list

	a := karytree.NewNode("a")
	b := karytree.NewNode("b")
	c := karytree.NewNode("c")
	d := karytree.NewNode("d")

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
	a := karytree.NewNode("a")
	if a.Key().(string) != "a" {
		t.Errorf("key was 'a', should not be %+v\n", a.Key().(string))
	}

	a.SetKey("b")
	if a.Key().(string) != "b" {
		t.Errorf("key was changed to 'b', should not be %+v\n", a.Key().(string))
	}
}

func TestSiblingTreeNLogic(t *testing.T) {
	a := karytree.NewNode("a")
	b := karytree.NewNode("b")
	c := karytree.NewNode("c")
	d := karytree.NewNode("d")
	e := karytree.NewNode("e")

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
