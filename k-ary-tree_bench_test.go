package karytree_test

import (
	"testing"

	"github.com/sevagh/k-ary-tree"
)

func BenchmarkKaryTreeK2Sparse(b *testing.B) {
	prevTree := karyTreeKSparseHelper(2)

	b.ResetTimer()
	var tree karytree.Node

	for i := 0; i < b.N; i++ {
		tree = karyTreeKSparseHelper(2)

		// stop the timer, compare trees
		b.StopTimer()
		if !karytree.Equals(&tree, &prevTree) {
			b.Errorf("Benching Sparse K=2 small trees but I don't think they're identical...")
		}
		prevTree = tree
		b.StartTimer()
	}
}

func BenchmarkKaryTreeK2VerySparse(b *testing.B) {
	prevTree := karyTreeKVerySparseHelper(2)

	b.ResetTimer()
	var tree karytree.Node

	for i := 0; i < b.N; i++ {
		tree = karyTreeKVerySparseHelper(2)

		// stop the timer, compare trees
		b.StopTimer()
		if !karytree.Equals(&tree, &prevTree) {
			b.Errorf("Benching VerySparse K=2 small trees but I don't think they're identical...")
		}
		prevTree = tree
		b.StartTimer()
	}
}

func BenchmarkKaryTreeK2Complete(b *testing.B) {
	prevTree := karyTreeKCompleteHelper(2)

	b.ResetTimer()
	var tree karytree.Node

	for i := 0; i < b.N; i++ {
		tree = karyTreeKCompleteHelper(2)

		// stop the timer, compare trees
		b.StopTimer()
		if !karytree.Equals(&tree, &prevTree) {
			b.Errorf("Benching Complete K=2 small trees but I don't think they're identical...")
		}
		prevTree = tree
		b.StartTimer()
	}
}

func BenchmarkKaryTreeK8Sparse(b *testing.B) {
	prevTree := karyTreeKSparseHelper(8)

	b.ResetTimer()
	var tree karytree.Node

	for i := 0; i < b.N; i++ {
		tree = karyTreeKSparseHelper(8)

		// stop the timer, compare trees
		b.StopTimer()
		if !karytree.Equals(&tree, &prevTree) {
			b.Errorf("Benching Sparse K=8 small trees but I don't think they're identical...")
		}
		prevTree = tree
		b.StartTimer()
	}
}

func BenchmarkKaryTreeK8VerySparse(b *testing.B) {
	prevTree := karyTreeKVerySparseHelper(8)

	b.ResetTimer()
	var tree karytree.Node

	for i := 0; i < b.N; i++ {
		tree = karyTreeKVerySparseHelper(8)

		// stop the timer, compare trees
		b.StopTimer()
		if !karytree.Equals(&tree, &prevTree) {
			b.Errorf("Benching VerySparse K=8 small trees but I don't think they're identical...")
		}
		prevTree = tree
		b.StartTimer()
	}
}

func BenchmarkKaryTreeK8Complete(b *testing.B) {
	prevTree := karyTreeKCompleteHelper(8)

	b.ResetTimer()
	var tree karytree.Node

	for i := 0; i < b.N; i++ {
		tree = karyTreeKCompleteHelper(8)

		// stop the timer, compare trees
		b.StopTimer()
		if !karytree.Equals(&tree, &prevTree) {
			b.Errorf("Benching Complete K=8 small trees but I don't think they're identical...")
		}
		prevTree = tree
		b.StartTimer()
	}
}

func BenchmarkKaryTreeK32Sparse(b *testing.B) {
	prevTree := karyTreeKSparseHelper(32)

	b.ResetTimer()
	var tree karytree.Node

	for i := 0; i < b.N; i++ {
		tree = karyTreeKSparseHelper(32)

		// stop the timer, compare trees
		b.StopTimer()
		if !karytree.Equals(&tree, &prevTree) {
			b.Errorf("Benching Sparse K=32 small trees but I don't think they're identical...")
		}
		prevTree = tree
		b.StartTimer()
	}
}

func BenchmarkKaryTreeK32VerySparse(b *testing.B) {
	prevTree := karyTreeKVerySparseHelper(32)

	b.ResetTimer()
	var tree karytree.Node

	for i := 0; i < b.N; i++ {
		tree = karyTreeKVerySparseHelper(32)

		// stop the timer, compare trees
		b.StopTimer()
		if !karytree.Equals(&tree, &prevTree) {
			b.Errorf("Benching VerySparse K=32 small trees but I don't think they're identical...")
		}
		prevTree = tree
		b.StartTimer()
	}
}

func BenchmarkKaryTreeK32Complete(b *testing.B) {
	prevTree := karyTreeKCompleteHelper(32)

	b.ResetTimer()
	var tree karytree.Node

	for i := 0; i < b.N; i++ {
		tree = karyTreeKCompleteHelper(32)

		// stop the timer, compare trees
		b.StopTimer()
		if !karytree.Equals(&tree, &prevTree) {
			b.Errorf("Benching Complete K=32 small trees but I don't think they're identical...")
		}
		prevTree = tree
		b.StartTimer()
	}
}

func karyTreeKSparseHelper(K int) karytree.Node {
	var tree karytree.Node

	var key int

	tree = karytree.NewNode(key)
	key++

	var curr *karytree.Node
	curr = &tree

	for i := uint16(0); i < uint16(K); i++ {
		if i%2 == 0 {
			child := karytree.NewNode(key)
			key++

			// fill even children
			curr.SetNthChild(i, &child)
			for j := uint16(0); j < uint16(K); j++ {
				if j%2 != 0 {
					grandchild := karytree.NewNode(key)
					key++
					ith := curr.NthChild(i)

					// fill odd grandchildren
					ith.SetNthChild(j, &grandchild)
					for k := uint16(0); k < uint16(K); k++ {
						if k%2 == 0 {
							// fill even great grandchildren
							greatgrandchild := karytree.NewNode(key)
							key++
							jth := ith.NthChild(j)

							jth.SetNthChild(k, &greatgrandchild)
						}
					}
				}
			}
		}
	}

	return tree
}

func karyTreeKCompleteHelper(K int) karytree.Node {
	var tree karytree.Node

	var key int

	tree = karytree.NewNode(key)
	key++

	var curr *karytree.Node
	curr = &tree

	for i := uint16(0); i < uint16(K); i++ {
		child := karytree.NewNode(key)
		key++
		curr.SetNthChild(i, &child)
		for j := uint16(0); j < uint16(K); j++ {
			grandchild := karytree.NewNode(key)
			key++

			ith := curr.NthChild(i)
			ith.SetNthChild(j, &grandchild)
			for k := uint16(0); k < uint16(K); k++ {
				greatgrandchild := karytree.NewNode(key)
				key++

				jth := ith.NthChild(j)
				jth.SetNthChild(k, &greatgrandchild)
			}
		}
	}

	return tree
}

func karyTreeKVerySparseHelper(K int) karytree.Node {
	var tree karytree.Node

	var key int

	tree = karytree.NewNode(key)
	key++

	var curr *karytree.Node
	curr = &tree

	for i := uint16(0); i < uint16(K); i++ {
		if i == 0 {
			child := karytree.NewNode(key)
			key++

			// fill even children
			curr.SetNthChild(i, &child)
			for j := uint16(0); j < uint16(K); j++ {
				if j == 0 {
					grandchild := karytree.NewNode(key)
					key++
					ith := curr.NthChild(i)

					// fill odd grandchildren
					ith.SetNthChild(j, &grandchild)
					for k := uint16(0); k < uint16(K); k++ {
						if k == 0 {
							// fill even great grandchildren
							greatgrandchild := karytree.NewNode(key)
							key++
							jth := ith.NthChild(j)

							jth.SetNthChild(k, &greatgrandchild)
						}
					}
				}
			}
		}
	}

	return tree
}
