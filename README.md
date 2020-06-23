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
type Node struct {
	key         interface{}
	n           uint
	firstChild  *Node
	nextSibling *Node
}
```

### Generics

Months ago, I experimented with genny to generate a copy of k-ary-tree with `uint32` keys in an experimental branch of [quadtree-compression](https://github.com/sevagh/quadtree-compression/tree/k-ary-tree-experiment) project with OK (but not great) results. The sibling linked-list implementation incurs some CPU costs.

On 2020/06/22, I experimented with the [new generics proposal](https://go.googlesource.com/proposal/+/master/design/go2draft-generics-overview.md) ([overview](https://blog.golang.org/generics-next-step), [go2go README](https://go.googlesource.com/go/+/refs/heads/dev.go2go/README.go2go.md)). The resulting code is [here](./k-ary-tree_generics.go2). I again attempted to use it in a real project, [quadtree-compression](https://github.com/sevagh/quadtree-compression/tree/generics-experiment), but go2go takes too long (30+ min with no output). Presumably the `go2go` tool isn't ready to compile code with significant dependencies yet.

### Sibling list operations

Think of the children list as a linkedlist:

```
firstChild -> nextSibling -> nextSibling -> nil
```

The linked list is sorted on insertions since tree comparisons will depend on the sortedness:

```
a := NewNode("a")
b := NewNode("b")
c := NewNode("c")

a_ := NewNode("a")
b_ := NewNode("b")
c_ := NewNode("c")

a.SetNthChild(4, &b)
// a -> firstChild -> (b, n: 4)

a.SetNthChild(2, &c)
// a -> firstChild -> (c, n: 2) -> nextSibling -> (b, n: 4)

a_.SetNthChild(2, &c_)
// a_ -> firstChild -> (c_, n: 2)

a_.SetNthChild(4, &b_)
// a_ -> firstChild -> (c_, n: 2) -> nextSibling -> (b_, n: 4)

fmt.Println(kary.Equals(&a, &a_))
// True
```

The field `n` determines which child a node is. It's a `uint` which gives us plenty of headroom.
