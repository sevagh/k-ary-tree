# k-ary-tree

### SmallK optimizations

Modified to not allocate the slice of children unless it's demanded:

```go
func NewNodeSmallK(k int, key interface{}) NodeSmallK {
...
-       n.children = make([]*NodeSmallK, k)
        return n
 }

func (k *NodeSmallK) SetNthChild(n int, other *NodeSmallK) error {
...
+       if k.children == nil {
+               k.children = make([]*NodeSmallK, k.k)
+       }
+

func (k *NodeSmallK) NthChildSmallK(n int) (*NodeSmallK, error) {
...
+       if k.children == nil {
+               return nil, nil
+       }
```

Results for K=2:

```
sevagh:zoo $ benchcmp old.txt new.txt
benchmark                              old ns/op     new ns/op     delta
BenchmarkKaryTreeSmallK2Sparse-8       591           588           -0.51%
BenchmarkKaryTreeSmallK2Complete-8     1559          1305          -16.29%

benchmark                              old allocs     new allocs     delta
BenchmarkKaryTreeSmallK2Sparse-8       11             10             -9.09%
BenchmarkKaryTreeSmallK2Complete-8     44             36             -18.18%

benchmark                              old bytes     new bytes     delta
BenchmarkKaryTreeSmallK2Sparse-8       352           336           -4.55%
BenchmarkKaryTreeSmallK2Complete-8     1312          1184          -9.76%
```

Bigger gains in the complete tree - makes sense, likely from the removal of pointlessly allocated slices for the empty children of the leaf nodes.

### SmallK vs. LargeK

- Complete: tree is complete at every height for h=3
- Sparse: tree has ~50% populated children at each height (alternating odd and even children to fill) for h=3
- VerySparse: imbalanced tree that's a linkedlist at the 0th successor for h=3 - 4 nodes total

|  | Small | Large |
| - | ----- | ----- |
| Children | <code>children []*NodeSmallK</code> | <code>children map[int]*NodeLargeK</code> | 
| Bench K2VerySparse | 588 ns/op             336 B/op         10 allocs/op | 954 ns/op             848 B/op         14 allocs/op |
| Bench K2Sparse   | 569 ns/op 336 B/op 10 allocs/op | 935 ns/op 848 B/op 14 allocs/op |
| Bench K2Complete | 1380 ns/op 1184 B/op 36 allocs/op | 2708 ns/op 2560 B/op 51 allocs/op |
| Bench K8VerySparse | 622 ns/op             480 B/op         10 allocs/op | 963 ns/op             848 B/op         14 allocs/op |
| Bench K8Sparse   | 5850 ns/op 7456 B/op 190 allocs/op | 11959 ns/op 11856 B/op 275 allocs/op |
| Bench K8Complete | 32936 ns/op 46784 B/op 1242 allocs/op | 82513 ns/op 71344 B/op 1827 allocs/op |
| Bench K32VerySparse | 861 ns/op            1056 B/op         10 allocs/op |  1009 ns/op             848 B/op         14 allocs/op |
| Bench K32Sparse | 283698 ns/op          384457 B/op       9010 allocs/op |  895696 ns/op          741177 B/op      14005 allocs/op |
| Bench K32Complete | 2112040 ns/op         2706015 B/op      68706 allocs/op | 7559845 ns/op         5968475 B/op     106864 allocs/op |

The results indicate to me that using a map of children has no tangible benefits at any reasonable value of K, especially since the Small optimization mentioned above.
### Simulating Python generators with channels in Go

Inspired by [this](https://blog.carlmjohnson.net/post/on-using-go-channels-like-python-generators/).

In the iterative traversals, I use channels to behave like Python generators. Let's look at the simplest - the BFS:

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

			child := curr.Left()
			if child != nil {
				queue = append(queue, child)
			}

			child = curr.Right()
			if child != nil {
				queue = append(queue, child)
			}
		}
	}()

	return nChan
}
```

If there were a `yield` keyword in Go, this is what the above code might look like:

```go
func BFS(root *Node}) *Node {
    queue := []*Node{root}
    var curr *Node

    for len(queue) > 0 {
        curr, queue = queue[0], queue[1:]

        yield curr

        child := curr.Left()
        if child != nil {
            queue = append(queue, child)
        }

        child = curr.Right()
        if child != nil {
            queue = append(queue, child)
        }
    }
	return
}
```

The channel-based goroutine tricks are there to "suspend and resume" the execution of the function like a Python generator.

### The visit trick for recursion

I reached for the above solution because I couldn't think of how I would create a recursive tree walk (that returned nodes) until I came across the visitor function idea [here](https://flaviocopes.com/golang-data-structure-binary-search-tree/).

```go
func InorderRecursive(root *Node, f func(*Node)) {
        inorder(root, f)
}

func inorder(root *Node, f func(*Node)) {
        if root != nil {
                inorder(root.Left(), f)
                f(root)
                inorder(root.Right(), f)
        }
}
```

Driver code:

```go
nodes := []*binarytree.Node{}
visit := func(node *binarytree.Node) {
        nodes = append(nodes, node)
}

binarytree.InorderRecursive(&a, visit)
```

Apparently this has a name: [the visitor pattern](https://en.wikipedia.org/wiki/Visitor_pattern). Cool.
