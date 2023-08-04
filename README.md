# k-ary-tree

A child-sibling k-ary tree. [This blog post](https://blog.mozilla.org/nnethercote/2012/03/07/n-ary-trees-in-c/) has a neat visualization of the firstChild/nextSibling tree structure:
```
       ---
      | a |
       ---
     /
    / firstChild
   /
 ---              ---      ---
| b | ---------- | c | -- | d |
 ---   nextSib    ---      ---
  |                         |
  |                         | firstChild
  | firstChild              |
  |                        ---
 ---       ---            | g |
| e | --- | f |            ---
 ---       ---
```

Advantages include simple traversals and tree operations, e.g. in a BFS:
```go
next := curr.firstChild
for next != nil {
    queue = append(queue, next)
    next = next.nextSibling
}
``` 

A karytree.Node is defined as:

```go
type Node[T comparable] struct {
	key         T
	n           uint
	firstChild  *Node
	nextSibling *Node
}
```

### Sibling list operations

Think of the children list as a linkedlist:

```
firstChild -> nextSibling -> nextSibling -> nil
```

The linked list is sorted on insertions since tree comparisons will depend on the sortedness:

```
import (
	"fmt"

	"github.com/sevagh/k-ary-tree"
)

func main() {
    a := karytree.NewNode("a")
	b := karytree.NewNode("b")
	c := karytree.NewNode("c")

	a_ := karytree.NewNode("a")
	b_ := karytree.NewNode("b")
	c_ := karytree.NewNode("c")

	a.SetNthChild(4, &b)
	// a -> firstChild -> (b, n: 4)

	a.SetNthChild(2, &c)
	// a -> firstChild -> (c, n: 2) -> nextSibling -> (b, n: 4)

	a_.SetNthChild(2, &c_)
	// a_ -> firstChild -> (c_, n: 2)

	a_.SetNthChild(4, &b_)
	// a_ -> firstChild -> (c_, n: 2) -> nextSibling -> (b_, n: 4)

	fmt.Println(karytree.Equals(&a, &a_))
	// Output: true
}
```

The field `n` determines which child a node is. It's a `uint` which gives us plenty of headroom.
