package main

import (
	"errors"
)

// basic doubly and singly linked lists

type DoubleLink struct {
	Content rune
	Next    *DoubleLink
	Prior   *DoubleLink
}

func (dl *DoubleLink) TForward() *DoubleLink {
	return dl.Next
}

func (dl *DoubleLink) TReverse() *DoubleLink {
	return dl.Prior
}

func (dl *DoubleLink) Link(d *DoubleLink) {
	// Links the new node to the right 
}

type DoublyLinkedList struct {
	Head   *DoubleLink
	End    *DoubleLink
	Length int
}

func (dll *DoublyLinkedList) getLen() int {
	return dll.Length
}

func (dll *DoublyLinkedList) TraverseTo(i int) (*DoubleLink, error) {
	if i > dll.Length {
		return nil, errors.New("Index exceeds list length")
	}
	if i == 0 {
		return dll.Head, nil
	} else if i == dll.Length-1 {
		return dll.End, nil
	}
	var ret *DoubleLink
	if i < dll.Length/2 {
		ret = dll.Head
		for ci := 0; ci < i; ci++ {
			ret = ret.TForward()
		}
	} else if i >= dll.Length/2 {
		ret = dll.End
		for ci := dll.Length - 1; ci > i; ci-- {
			ret = ret.TReverse()
		}
	}

	return ret, nil
}

func (dll *DoublyLinkedList) TraverseRangeF(i, j int) ([]*DoubleLink, error) {
	if j < i {
		return nil, errors.New("Bad indices")
	}
	initial, err := dll.TraverseTo(i)
	ret := []*DoubleLink{initial}
	if err != nil {
		return nil, err
	}
	for i := i; i < j; i++ {
		initial = initial.TForward()
		ret = append(ret, initial)
	}
	return ret, nil
}

func (dll *DoublyLinkedList) TraverseRangeR(i, j int) ([]*DoubleLink, error) {
	if i < j {
		return nil, errors.New("Bad indices")
	}
	initial, err := dll.TraverseTo(j)
	ret := []*DoubleLink{initial}
	if err != nil {
		return nil, err
	}
	for j := j; j > i; j-- {
		initial = initial.TReverse()
		ret = append(ret, initial)
	}
	return ret, nil
}

func (dll *DoublyLinkedList) Insert(i int, r rune) error {
	// inserts after index i (replaces current index i with new node with value r)
	if i > dll.Length || i < 0 {
		return errors.New("bad index")
	}
	nl := DoubleLink{
		Content: r,
		Prior:   nil,
		Next:    nil,
	}
	if dll.Length == 0 {
		dll.Head = &nl
		dll.End = &nl
		dll.Length++
		return nil
	}

	if i == 0 {
		// case: start
		nl.Next = dll.Head
		dll.Head.Prior = &nl
		dll.Head = &nl
	} else if i == dll.Length {
		// case: end
		nl.Prior = dll.End
		dll.End.Next = &nl
		dll.End = &nl
	} else {
		p, _ := dll.TraverseTo(i - 1)
		a := p.Next
		p.Next = &nl
		a.Prior = &nl
		nl.Next = a
		nl.Prior = p
	}

	dll.Length++
	return nil
}

func (dll *DoublyLinkedList) Remove(i int) (rune, error) {
	removed, err := dll.TraverseTo(i)
	if err != nil {
		return -1, errors.New("Bad traverse during removal")
	}
	removed.Prior.Next = removed.Next
	removed.Next.Prior = removed.Prior
	dll.Length--
	return removed.Content, nil

}

func (dll *DoublyLinkedList) Split(i int) (*DoublyLinkedList, error) {
	// splits the linked list into segments 0...i, i+1 ... end. returns the portion i+1 ... end.
	splitEndNode, err := dll.TraverseTo(i)
	if err != nil {
		return nil, errors.New("error in traversing to index in split")
	}
	splitStartNode := splitEndNode.Next

	splitEndNode.Next = nil
	splitStartNode.Prior = nil
	right := &DoublyLinkedList{
		Head:   splitStartNode,
		End:    dll.End,
		Length: dll.Length - i - 1,
	}
	dll.End = splitEndNode
	dll.Length = i + 1

	return right, nil

}

/*
func (dll *DoublyLinkedList) Swap(i int, j int) error {
	first, err := dll.TraverseTo(i)
	if err != nil {
		return errors.New("bad traverse in swap")
	}
	second, err := dll.TraverseTo(j)
	if err != nil {
		return errors.New("bad traverse in swap")
	}

	// TODO: handle special cases of start/end

	firstNext, firstPrior := first.Next, first.Prior


	first.Prior, first.Next = second.Prior, second.Next

	second.Prior, second.Next = firstPrior, firstNext
	if (i == 0){
		dll.Head = second
	} else if (j == 0){
		dll.Head = first
	} else if (i == dll.Length - 1){
		dll.End = second
	} else if (j == dll.Length - 1){
		dll.End = first
	}

	return nil
}
*/
func (dll *DoublyLinkedList) ToString() string {
	ret := ""
	fn := dll.Head
	for fn != nil {
		ret += string(fn.Content)
		fn = fn.Next
	}
	return ret
}
