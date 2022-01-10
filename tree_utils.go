package main

import (
	"fmt"
	"sort"
)

func PrintInline(head *AVLNode) {
	/*
		Given the ``head`` of an AVL Rope, print out the non-nil items in order.
	*/
	r := make([]string, 0)
	head.ApplyInorder(func(n *AVLNode) { r = append(r, (string)(n.Value)) })
	for _, e := range r {
		fmt.Printf("%s,", e)

	}
	fmt.Println()
}

func LeavesOnly(head *AVLNode) []*AVLNode {
	/*
		Given the ``head`` of an AVL Rope, verify that only leaves store data. Returns offenders.
	*/
	r := true
	ret := make([]*AVLNode, 0)
	head.ApplyInorder(func(n *AVLNode) {
		if !n.IsLeaf() && n.Value != nil {
			r = false
			ret = append(ret, n)
		}
	})
	fmt.Printf("Leaves only? %t\n", r)
	return ret
}

func IntAbs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func CheckAVL(head *AVLNode) []*AVLNode {
	/*
		Given the ``head`` of an AVL Rope, verify that it conforms to AVL specifications. Returns offenders.
	*/
	r := true
	ret := make([]*AVLNode, 0)
	head.ApplyInorder(func(n *AVLNode) {
		if IntAbs(n.L.getHeight()-n.R.getHeight()) > 1 {
			r = false
			ret = append(ret, n)
		}
	})
	fmt.Printf("AVL? %t", r)
	return ret
}

func TreeStats(head *AVLNode) {
	// Print out tree stats
	leafct := 0
	nodect := 0
	heightcts := make([]int, 0)
	head.ApplyInorder(func(n *AVLNode) {
		if n.IsLeaf() {

			leafct++
			tmp := n.U
			ht := 0
			for tmp != nil {
				ht++
				tmp = tmp.U
			}
			heightcts = append(heightcts, ht)
		}
		nodect++
	})
	sort.Ints(heightcts)
	fmt.Printf("Leaf Ct: %d \n Node Ct: %d \n Max Height: %d\n", leafct, nodect, heightcts[len(heightcts)-1])
}

/*
func CheckHeight(head *AVLNode) []*AVLNode {

	//	Given the ``head`` of an AVL Rope, verify that all heights are accurate (with leaves being height 0). Returns offenders.

	r := true
	ret := make([]*AVLNode, 0)
	head.ApplyInorder(func(n *AVLNode) {
		if (n == n.U.L){

		} else if n == (n.U.R){

		}
		// todo: leaf case
	})
}

func CheckWeight(head)

*/
