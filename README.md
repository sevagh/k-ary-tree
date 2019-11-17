# k-ary-tree

A sibling-sibling k-ary tree. [This blog post](https://blog.mozilla.org/nnethercote/2012/03/07/n-ary-trees-in-c/) has a neat visualization of the firstChild/nextSibling tree structure:

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

A karytree.Node is defined as:

```go
type Node struct {
	key         interface{}
	parent      uintptr
	firstChild  *Node
	nextSibling *Node
}
```

There are convenience wrappers for binary trees and region quadtrees (as many trees are specialized forms of k-ary-trees).

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

I use the top 16 bits of the parent uintptr to store the `n` value of a child. On amd64 architectures, apparently only 48 bits of 64-bit pointers are used for addressing (I didn't do any thorough research but tests are passing). Two notes:

* This is highly unrecommended and I use the unsafe package to do it, but I need to spice up my life somehow
* This should limit the k of our k-ary-tree to 65536 - who could ever need more than that?

The `key` stores any data the user wants. Nodes have a parent pointer and the library panics on re-parenting nodes with parents. It's a very "manual" library - you're expected to manage nodes, create links, etc. yourself.
