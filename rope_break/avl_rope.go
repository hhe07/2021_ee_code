package main

import (
	"errors"
	"fmt"
)

const (
	Segsize = 10
)

func JoinRight(L, R *AVLNode) *AVLNode {
	// Joins L and R together with right bias
	ll, lr := L.L, L.R
	if lr.getHeight() <= (R.getHeight() + 1) {
		inter := mkNode(nil, lr, R)
		if inter.getHeight() <= (ll.getHeight() + 1) {
			return mkNode(nil, ll, inter)
		}
		return mkNode(nil, ll, inter.rotateRight()).rotateLeft()
	}
	inter := JoinRight(lr, R)
	if inter.getHeight() <= (ll.getHeight() + 1) {
		return mkNode(nil, ll, inter)
	}
	return mkNode(nil, ll, inter).rotateLeft()
}

func JoinLeft(L, R *AVLNode) *AVLNode {
	// Mirror of above
	rl, rr := R.L, R.R
	if rl.getHeight() <= (L.getHeight() + 1) {
		inter := mkNode(nil, L, rl)
		if inter.getHeight() <= (rr.getHeight() + 1) {
			return mkNode(nil, inter, rr)
		}
		return mkNode(nil, inter.rotateLeft(), rr).rotateRight()
	}
	inter := JoinLeft(L, rl)
	if inter.getHeight() <= (rr.getHeight() + 1) {
		return mkNode(nil, inter, rr)
	}
	return mkNode(nil, inter, rr).rotateRight()
}

func Join(L, R *AVLNode) *AVLNode {
	if L.getHeight() > (R.getHeight() + 1) {
		return JoinRight(L, R)
	} else if R.getHeight() > (L.getHeight() + 1) {
		return JoinLeft(L, R)
	}
	return mkNode(nil, L, R)
}

func Split(ar *AVLNode, weight int) (*AVLNode, *AVLNode) {
	if ar == nil {
		return nil, nil
	}
	if ar.isLeaf() {

		return mkLeaf(ar.Value[:weight]), mkLeaf(ar.Value[weight:])
	}
	L, R := ar.L, ar.R
	if ar.Weight == weight {
		return L, R
	}
	if weight < ar.Weight {
		nL, nR := Split(L, weight)
		return nL, Join(nR, R)
	}
	// otherwise, weight > ar.Weight
	nL, nR := Split(R, weight-ar.Weight)
	return Join(L, nL), nR
}

func TraverseRange(start *AVLNode, startPos, len int) []rune {
	/*
		Returns ``len`` characters starting from startPos, inclusive
	*/
	ret := start.Value[startPos:]
	// if the starting point isn't a leaf, don't bother
	if !(start.isLeaf()) {
		return nil
	}
	return ret

}

func mkNode(val []rune, L, R *AVLNode) *AVLNode {
	ret := &AVLNode{
		Value: val,
	}
	ret = ret.linkLeft(L)
	ret = ret.linkRight(R)

	ret.updateHeight()
	ret.updateWeight()
	return ret
}

func mkLeaf(val []rune) *AVLNode {
	return &AVLNode{
		Value:  val,
		Height: 0,
		Weight: len(val),
	}
}

type AVLRope struct {
	Head *AVLNode
}

func mkRope() *AVLRope {
	return &AVLRope{}
}

func (ar *AVLRope) Report(content []rune) ([]rune, error) {
	return ar.Head.ToRune(), nil
}

func (ar *AVLRope) IndexNode(i int) (*AVLNode, int, error) {
	node, num := ar.Head.IndexNode(i)
	return node, num, nil
}

func (ar *AVLRope) ReportRange(i, j int) ([]rune, error) {
	ret := make([]rune, 0)
	// get start point
	from, ct, err := ar.IndexNode(i)
	if err != nil {
		return nil, errors.New("bad initial index")
	}
	// get len to be reported
	reportLen := j - i + 1 // remove this +1 to make exclusive
	if len(from.Value)-ct < reportLen {
		ret = append(ret, from.Value[ct:]...)
		reportLen -= len(ret)
	} else {
		ret = append(ret, from.Value[ct:ct+reportLen]...)
		return ret, nil
	}

	// then, while the length to be read has yet to be reached:
	curr := from.U
	for reportLen > 0 {
		// if came from right of curr, continue ascending
		if from == curr.R {
			from = curr
			curr = curr.U
		} else if from == curr.L {
			// if came from left of curr, go down to right
			from = curr
			curr = curr.R
		} else if from == curr.U {
			// otherwise, if coming from higher, go down to left
			from = curr
			curr = curr.L
		}

		if curr.isLeaf() {
			// insert as much data as possible
			if len(curr.Value) < reportLen {
				ret = append(ret, curr.Value...)
				reportLen -= len(curr.Value)
			} else {
				ret = append(ret, curr.Value[:reportLen]...)
				return ret, nil
			}
			// begin backtrack
			from = curr
			curr = curr.U
		}

	}
	return ret, nil

}

func (ar *AVLRope) Insert(i int, content []rune) error {
	new := mkRope()
	new.LoadFromRune(content)
	l, r := Split(ar.Head, i)
	l = Join(l, new.Head)
	ar.Head = Join(l, r)
	return nil
}

func (ar *AVLRope) Append(content []rune) error {
	// since this isn't actually tested, doesn't need to be efficient
	new := mkRope()
	new.LoadFromRune(content)
	ar.Head = Join(ar.Head, new.Head)
	ar.Head.updateWeight()
	return nil
}

/*
func (ar *AVLRope) Split(i int) (StorageType, error) {
	l, r := Split(ar.Head, i)
	ar.Head = l
	return r, nil
}
*/
func (ar *AVLRope) DeleteRange(i, j int) ([]rune, error) {
	l1, r := Split(ar.Head, j)
	l2, m := Split(l1, i)
	ar.Head = Join(l2, r)
	return m.ToRune(), nil
}

func (ar *AVLRope) LoadFromRune(contents []rune) {
	// step 1: create the necessary leaves

	leaves := make([]*AVLNode, (len(contents) / Segsize))
	for i := range leaves {

		data := make([]rune, Segsize)
		for r := 0; r < Segsize; r++ {
			data[r] = contents[(i*Segsize)+r]
		}
		leaves[i] = mkLeaf(data)
	}
	// step 2: recursively concat until something resembling balance has been achieved
	// TODO: check whether current or previous ends up having more allocs: in case slices are being weird
	for len(leaves) > 1 {
		half := len(leaves) / 2
		if len(leaves)%2 == 0 {
			for r := 0; r < half; r++ {
				leaves[r] = mkNode(nil, leaves[2*r], leaves[(2*r)+1])
			}
		} else {
			for r := 0; r < half; r++ {
				leaves[r] = mkNode(nil, leaves[2*r], leaves[(2*r)+1])
			}
			leaves[half-1] = mkNode(nil, leaves[half-1], leaves[len(leaves)-1])
		}
		leaves = leaves[:half]
	}
	if len(leaves) != 1 {
		panic("bad construction")
	}
	ar.Head = leaves[0]
}

func (ar *AVLRope) Load(contents []byte) {
	/*
		TODO: make sure to update this in the Gap load method as well
		TODO: currently these tests include the file opening time, which is Very Bad.
		TODO: also make sure that actual memory reallocs are happening because slicing a slice
		TODO: doesn't create any new memory
	*/

	// step 1: create the necessary leaves

	leaves := make([]*AVLNode, (len(contents) / Segsize))
	for i := range leaves {

		data := make([]rune, Segsize)
		for r := 0; r < Segsize; r++ {
			data[r] = rune(contents[(i*Segsize)+r])
		}
		leaves[i] = mkLeaf(data)
	}
	// step 2: recursively concat until something resembling balance has been achieved
	// TODO: check whether current or previous ends up having more allocs: in case slices are being weird
	//var a, b *AVLNode
	for len(leaves) > 1 {
		/*
			fmt.Println("NEW CYCLE")
			for _, c := range leaves {
				PrintInline(c)
			}
			fmt.Println()

			a, b, leaves = leaves[0], leaves[1], leaves[2:]
			if b.getHeight() != a.getHeight() {
				fmt.Println("yes")
				// case where odd number of leaves in queue, so this must be concatted to the end
				leaves[len(leaves)-1] = mkNode(nil, leaves[len(leaves)-1], a)
				leaves = append([]*AVLNode{b}, leaves...)
				//a = b
				//b, leaves = leaves[0], leaves[1:]
			} else {
				leaves = append(leaves, mkNode(nil, a, b))
			}

			for _, c := range leaves {
				PrintInline(c)
			}
			fmt.Println()
		*/

		half := len(leaves) / 2
		fmt.Println()
		if len(leaves)%2 == 0 {
			for r := 0; r < half; r++ {
				leaves[r] = mkNode(nil, leaves[2*r], leaves[(2*r)+1])
			}
		} else {
			for r := 0; r < half; r++ {
				leaves[r] = mkNode(nil, leaves[2*r], leaves[(2*r)+1])
			}
			leaves[half-1] = mkNode(nil, leaves[half-1], leaves[len(leaves)-1])
		}
		leaves = leaves[:half]
	}
	if len(leaves) != 1 {
		panic("bad construction")
	}
	ar.Head = leaves[0]
}

func (ar *AVLRope) ToString() string {
	ret := ""
	ar.Head.ApplyInorder(func(n *AVLNode) { ret += (string(n.Value) + ",") })
	return ret
}
