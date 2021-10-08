package main

import "errors"

type RopeContainer struct {
	Head *RopeNode
}

type RopeNode struct {
	Data     []rune
	Weight   int
	L        *RopeNode
	R        *RopeNode
	Upstream *RopeNode
}

// todo: make the one pointer to slice not that, it's pointless
// todo: defuculate errors

func createRope(r []rune) (*RopeContainer, error) {
	// todo: this needs to chunk things without breaking everything
	return nil, nil
}

func (r *RopeContainer) Balance() {
	// gonna conveniently ignore this and pretend it
	// doesnt exist
}

func (r *RopeContainer) Insert(pos int, s []rune) error { // we'll pretend everything is just one long string and do maths to figure out where things are
	// this is probably broken :(
	// 1 split and 2 concat
	nr, err := createRope(s)
	if err != nil {
		return errors.New("Could not create new rope")
	}
	ln, rn, err := r.Split(pos)
	if err != nil {
		return err
	}
	new1 := Concat(ln, nr.Head)
	r.Head = Concat(new1, rn)
	return nil
}

func (n *RopeNode) Index(i int) (rune, error) {

	rn, re, err := n.IndexNode(i)
	if err != nil {
		return -1, errors.New("bad indexnode output")
	}
	return rn.Data[re], nil
}

func (n *RopeNode) IndexNode(i int) (*RopeNode, int, error) {
	NotFoundError := errors.New("Position not found")
	if n.Weight <= i && (n.R != nil) {
		return n.R.IndexNode(i - n.Weight)
	}
	if n.L != nil {
		return n.L.IndexNode(i)
	}
	if i < len(n.Data) {
		return n, i, nil
	} else {
		return nil, -1, NotFoundError
	}

}

func (r *RopeContainer) Index(i int) (rune, error) {
	return r.Head.Index(i)
}

func (r *RopeContainer) IndexNode(i int) (*RopeNode, int, error) {
	return r.Head.IndexNode(i)
}

func Concat(S1, S2 *RopeNode) *RopeNode {
	ret := &RopeNode{
		Data:     nil,
		Weight:   S1.Weight,
		L:        S1,
		R:        S2,
		Upstream: nil,
	}
	S1.Upstream = ret
	S2.Upstream = ret
	return ret
}

func (r *RopeContainer) Concat(S1, S2 *RopeNode) {
	newNode := Concat(S1, S2)
	r.Head = newNode
	r.Balance()
}

func (r *RopeContainer) Split(i int) (*RopeNode, *RopeNode, error) {
	// TODO: TEST TEST TEST!
	sn, idx, err := r.IndexNode(i)
	if err != nil {
		return nil, nil, errors.New("bad index")
	}
	out := make([]*RopeNode, 0)
	// special case: i is in middle of a string
	if idx > 0 {
		nn := &RopeNode{
			Data:     sn.Data[i:], // this is prolly broken
			Weight:   sn.Weight - i,
			L:        nil,
			R:        nil,
			Upstream: nil,
		}
		sn.Weight = i
		sn.Data = sn.Data[:i]

		np := &RopeNode{
			Data:     nil,
			Weight:   nn.Weight,
			L:        sn,
			R:        nn,
			Upstream: sn.Upstream,
		}
		sn.Upstream = np
		nn.Upstream = np

		out = np.Split(sn, out) // will this break?
	}
	out = sn.Upstream.Split(sn, out)
	for len(out) > 1 {
		ne := Concat(out[0], out[1])
		out = append([]*RopeNode{ne}, out[2:]...) // todo: is this inclusive?
	}
	return r.Head, out[0], nil

}

func (r *RopeNode) Split(Origin *RopeNode, Components []*RopeNode) []*RopeNode {
	// TODO: TEST TEST TEST!
	/*
		should be run on upstream of the bottommost node that contains the index you want to split from. splits off that node.
		split by keeping track if you're on right or left of parent. if you're on right, split section off. if you're on left, continue up.
		Components is organised by which are visited first. hopefully nothing breaks.
	*/
	// base case: reached top
	if r.Upstream == nil {
		Components = append(Components, r.R)
		r.R.Upstream = nil
		r.R = nil
		return Components
	}
	// special case: coming from right of current node
	if Origin == r.R {
		Components = append(Components, r.R)
		r.R.Upstream = nil
		r.R = nil
		return r.Upstream.Split(r, Components)
	}
	// default case: going up
	return r.Upstream.Split(r, Components)

}

func (r *RopeContainer) Delete(i, j int) ([]rune, error) {
	// two split and one concat
	// first split: leaves string 0...i and i...n
	ln1, rn1, err := r.Split(i)
	if err != nil {
		return nil, errors.New("bad split in delete")
	}
	tmp := RopeContainer{Head: rn1}
	ln2, rn2, err := tmp.Split(j - i)
	if err != nil {
		return nil, errors.New("bad split in delete")
	}
	r.Head = Concat(ln1, rn2)
	return ln2.Data, nil
}

func (r *RopeContainer) Report(i, j int) ([]rune, error) {
	ret := make([]rune, 0)
	rns, startIdx, err := r.ReportNodes(i, j)
	if err != nil {
		return nil, err
	}
	rnsd := *rns
	ret = append(ret, rnsd[0].Data[startIdx:]...)
	rnsd = rnsd[1:]
	for _, node := range rnsd {
		if len(node.Data) > (j-i)-len(ret) {
			ret = append(ret, node.Data[:(j-i-len(ret))]...)
		} else {
			ret = append(ret, node.Data...)
		}
	}
	return ret, nil
}

func (r *RopeContainer) ReportNodes(i, j int) (*[]*RopeNode, int, error) {
	ret := make([]*RopeNode, 0)
	visitLength := 0
	visitStack := make([]*RopeNode, 0)
	sn, startIdx, err := r.IndexNode(i)
	if err != nil {
		return nil, -1, errors.New("Error finding starting position i")
	}
	if sn.Weight < j {
		return nil, -1, errors.New("Bad weight, tree's probably screwed up")
	}
	ret = append(ret, sn)
	node := sn
	for node != nil || !(len(visitStack) == 0) { // todo: ask edward/think about how this might be made more efficient with "upstream"
		if node != nil {
			visitStack = append(visitStack, node)
			node = node.L
		} else {
			node, visitStack = visitStack[len(visitStack)-1], visitStack[:len(visitStack)-1]
			ret = append(ret, node)
			visitLength += len(node.Data)
			node = node.R
		}
		if visitLength >= j {
			break
		}
	}
	return &ret, startIdx, nil
}
