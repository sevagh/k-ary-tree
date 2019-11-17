package karytree

import (
	"testing"
)

func TestBitHacksNoAPI(t *testing.T) {
	a := NewNode("a")
	if a.n() != 0 {
		t.Errorf("n should be 0")
	}
	a.setN(5)
	if a.n() != 5 {
		t.Errorf("n should be 5")
	}
	a.setN(3)
	if a.n() != 3 {
		t.Errorf("n should be 3")
	}

	b := NewNode("b")
	a.setFirstChild(&b)

	if a.n() != 3 {
		t.Errorf("setting child shouldn't affect n")
	}

	a.setN(65535)
	if a.n() != 65535 {
		t.Errorf("setting child shouldn't affect n")
	}

	c := NewNode("c")
	a.setFirstChild(&c)

	if a.getFirstChild() != &c {
		t.Errorf("overwriting firstChild should work")
	}
}

func TestBitHacksAPI(t *testing.T) {
	a := NewNode("a")
	b := NewNode("b")
	c := NewNode("c")

	a.SetNthChild(5, &b)

	if a.getFirstChild() != &b {
		t.Errorf("didn't store pointer correctly!")
	}
	if b.n() != 5 {
		t.Errorf("didn't set n of b correctly")
	}

	a.SetNthChild(5, &c)

	if a.getFirstChild() != &c {
		t.Errorf("didn't overwrite pointer correctly!")
	}
	if c.n() != 5 {
		t.Errorf("didn't set n of c correctly")
	}
}
