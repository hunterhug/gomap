package gomap

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestAvlBetterTree_Delete(t *testing.T) {
	tree := avlBetterTreeNode{}

	tree.balanceFactor++
	fmt.Println(tree)

	tree.balanceFactor += 1

	fmt.Println(tree)

	rand.Seed(1000000)
	for i := 0; i < 100; i++ {
		fmt.Println(rand.Int63n(1000000))
	}
}
