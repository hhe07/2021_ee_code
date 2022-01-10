package main

type AVLNode struct {
	Value  []rune
	Height int
	Weight int

	L *AVLNode
	R *AVLNode
	U *AVLNode
}

func (a *AVLNode) IsLeaf() bool {
	if a == nil {
		return false
	}
	return (a.L == nil) && (a.R == nil)
}

func (a *AVLNode) isLeaf() bool {
	return (a.L == nil) && (a.R == nil)
}

func (a *AVLNode) getHeight() int {
	if a != nil {
		return a.Height
	}
	return 0
}

func (a *AVLNode) updateHeight() {
	if a.isLeaf() {
		a.Height = 0
		return
	}
	a.Height = 1 + max(a.L.getHeight(), a.R.getHeight())
}

func (a *AVLNode) linkLeft(nl *AVLNode) *AVLNode {
	// normally, create new left link, return the current node
	if nl != nil {
		a.L = nl
		nl.U = a
		return a
	}
	// otherwise, compact
	return a.R
}

func (a *AVLNode) linkRight(nr *AVLNode) *AVLNode {
	// mirror of above
	if nr != nil {
		a.R = nr
		nr.U = a
		return a
	}
	return a.L
}

func (a *AVLNode) delinkLeft() *AVLNode {
	// delinks and returns old left
	tmp := a.L
	a.L.U = nil
	a.L = nil
	return tmp
}

func (a *AVLNode) delinkRight() *AVLNode {
	// delinks and returns old right
	tmp := a.R
	a.R.U = nil
	a.R = nil
	return tmp
}

func (a *AVLNode) sumTree() int {
	// returns the sum of leaves on this subtree
	if a == nil {
		return 0
	}
	if a.isLeaf() {
		return a.Weight
	}
	return a.Weight + a.R.sumTree()
}

func (a *AVLNode) updateWeight() {
	a.Weight = a.L.sumTree()
}

func (a *AVLNode) rotateLeft() *AVLNode {
	/*
		``a`` is the root of the subtree to be rotated left

		Returns the new root of the rotated subtree
	*/
	if a == nil {
		return a
	}
	z := a.R

	// check that rotation is actually required: disable during testing

	if (z == nil) || (z.L == nil) || (z.R == nil) {
		return a
	}

	inner := z.delinkLeft() // get inner child of right subtree
	// then, let x either be the newly rotated left a, or the compacted version of a if inner is nil
	x := a.linkRight(inner)
	// link x to z's left
	z.linkLeft(x)
	// clear z's U
	z.U = nil

	// update heights
	x.updateHeight()
	z.updateHeight()

	// update weights: since only z's weight is changing, this is
	// the only necessary change
	z.updateWeight()

	// return necessary value
	return z
}

func (a *AVLNode) rotateRight() *AVLNode {
	/*
		``a`` is the root of the subtree to be rotated right

		returns the new root of the rotated subtree
	*/
	if a == nil {
		return a
	}
	z := a.L

	// check that rotation is actually required: disable during testing
	if (z == nil) || (z.L == nil) || (z.R == nil) {
		return a
	}

	inner := z.delinkRight() // get inner child of left subtree
	// then, let x either be the newly rotated right a, or the compacted version of a if inner is nil
	x := a.linkLeft(inner)
	// link x to z's right
	z.linkRight(x)
	// clear z's U
	z.U = nil

	// update heights
	x.updateHeight()
	z.updateHeight()

	// update weights
	x.updateWeight()
	z.updateWeight()
	// return necessary value
	return z
}

func (a *AVLNode) balance() *AVLNode {
	/*
		``a`` is the root of the subtree to be balanced
		returns the new root of the balanced subtree
	*/
	if a == nil {
		return nil
	}
	balanceFactor := a.L.getHeight() - a.R.getHeight()
	if balanceFactor <= -2 {
		if a.R.L.getHeight() > a.R.R.getHeight() {
			a.R = a.R.rotateRight()
			// prolly gonna have a delightful number of ptr errors here
		}
		return a.rotateLeft()
	} else if balanceFactor >= 2 {
		if a.L.R.getHeight() > a.L.L.getHeight() {
			a.L = a.L.rotateLeft()
		}
		return a.rotateRight()
	}
	return a
}

func (a *AVLNode) ToRune() []rune {
	ret := make([]rune, 0)
	a.ApplyInorder(func(n *AVLNode) {
		ret = append(ret, n.Value...)
	})
	return ret
}

func (a *AVLNode) ApplyInorder(f func(n *AVLNode)) {
	if a == nil {
		return
	}
	a.L.ApplyInorder(f)
	if a.Value != nil {
		f(a)
	}
	a.R.ApplyInorder(f)
}

func (a *AVLNode) IndexNode(i int) (*AVLNode, int) {
	// finds and returns node that contains char at index i, returning the number of characters by which to advance in
	if a == nil {
		return nil, -1
	}
	if a.isLeaf() && (0 <= i) && (i < len(a.Value)) {
		return a, i
	}
	if a.Weight <= i {
		// go right
		return a.R.IndexNode(i - a.Weight)
	}
	if a.Weight > i {
		// go left
		return a.L.IndexNode(i)
	}
	return nil, -1
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
