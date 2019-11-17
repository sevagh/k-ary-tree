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
type KeyType generic.Type

type Node struct {
	key         KeyType
	firstChild  uintptr
	nextSibling *Node
}
```

The code works as if key were an `interface{}`, but by using genny you have the option of generating copies of [k-ary-tree.go](./k-ary-tree.go) with specific KeyTypes for higher performance. E.g. I use a generated copy with `KeyType = uint32` in my [quadtree-compression](https://github.com/sevagh/quadtree-compression) project with good results. There is a convenience wrapper for a binary tree to implement and demonstrate traversals.

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

I use the top 16 bits of the firstChild uintptr to store the `n` value of a child. On amd64 architectures, apparently only 48 bits of 64-bit pointers are used for addressing (I didn't do any thorough research but tests are passing). Two notes:

* This is highly unrecommended and I use the unsafe package to do it, but I need to spice up my life somehow
* This should limit the k of our k-ary-tree to 65536 - who could ever need more than that?

### Traversals

I simulate Python generators with channels in Go, inspired by [this](https://blog.carlmjohnson.net/post/on-using-go-channels-like-python-generators/). Let's look at the simplest - the BFS:

```go
func BFS(root *Node, quit <-chan struct{}) <-chan *Node {
	nChan := make(chan *Node)

	go func() {
		defer close(nChan)
		queue := []*Node{root}
		var curr *Node

		for len(queue) > 0 {
			curr, queue = queue[0], queue[1:]

			select {
			case <-quit:
				return
			case nChan <- curr:
			}

			next := curr.firstChild
			for next != nil {
				queue = append(queue, next)
				next = next.nextSibling
			}
		}
	}()

	return nChan
}
```

The channel-based goroutine tricks are there to "suspend and resume" the execution of the function like a Python generator. The recursive traversals are implemented using the visitor pattern.
