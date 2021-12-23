package main

import (
	"io"
)

var (
	// SplitLength is the threshold above which slices will be split into
	// separate nodes.
	SplitLength = 4096 * 4
	// JoinLength is the threshold below which nodes will be merged into
	// slices.
	JoinLength = SplitLength / 2
	// RebalanceRatio is the threshold used to trigger a rebuild during a
	// rebalance operation.
	RebalanceRatio = 1.2
)

type Rope struct {
	Root *Node
}

func (r *Rope) Report() ([]rune, error) {
	return r.Root.Value(), nil
}

func (r *Rope) ReportRange(i, j int) ([]rune, error) {
	return r.Root.Slice(i, j), nil
}

func (r *Rope) ReportCharacter(i int) (rune, error) {
	return r.Root.At(i), nil
}

func (r *Rope) Insert(i int, content []rune) error {
	r.Root.Insert(i, content)
	return nil
}

func (r *Rope) Append(content []rune) error {
	r.Root.Insert(r.Root.length, content) // todo: fix
	return nil
}

func (r *Rope) Split(i int) (StorageType, error) {
	x, y := r.Root.SplitAt(i - 1)
	r.Root = x
	return &Rope{Root: y}, nil
}

func (r *Rope) DeleteRange(i, j int) ([]rune, error) {
	res := r.Root.Remove(i, j, make([]rune, 0)) // todo: fix
	return res, nil
}

func (r *Rope) Concat(content []rune) error {
	r.Append(content)
	return nil
}

func (r *Rope) ToString() string {
	res, err := r.Report()
	if err != nil {
		panic(err)
	}
	ret := ""
	for _, e := range res {
		ret += string(e)
	}
	return ret
}

/*
func (r *Rope) Save(f *FileWrapper) error {

}
*/

func (r *Rope) Load(f *FileWrapper, ct int) error {
	content := f.Open(ct)
	in := make([]rune, 0)
	for _, c := range content {
		in = append(in, rune(c))
	}
	r.Root = New(in)
	return nil
}

type nodeType byte

const (
	tLeaf nodeType = iota
	tNode
)

// A Node in the rope structure. If the kind is tLeaf, only the value and
// length are valid, and if the kind is tNode, only length, left, right are
// valid.
type Node struct {
	kind        nodeType
	value       []rune
	length      int
	left, right *Node
}

// New returns a new rope node from the given byte slice. The underlying
// data is not copied so the user should ensure that it is okay to insert and
// delete from the input slice.
func New(b []rune) *Node {
	n := &Node{
		kind:   tLeaf,
		value:  b[0:len(b):len(b)],
		length: len(b),
	}
	n.adjust()
	return n
}

// Len returns the number of elements stored in the rope.
func (n *Node) Len() int {
	return n.length
}

func (n *Node) adjust() {
	switch n.kind {
	case tLeaf:
		if n.length > SplitLength {
			divide := n.length / 2
			n.left = New(n.value[:divide])
			n.right = New(n.value[divide:])
			n.value = nil
			n.kind = tNode
			n.length = n.left.length + n.right.length
		}
	case tNode:
		if n.length < JoinLength {
			n.value = n.Value()
			n.left = nil
			n.right = nil
			n.kind = tLeaf
			n.length = len(n.value)
		}
	}
}

// Value returns the elements of this node concatenated into a slice. May
// return the underyling slice without copying, so do not modify the returned
// slice.
func (n *Node) Value() []rune {
	switch n.kind {
	case tLeaf:
		return n.value
	case tNode:
		return concat(n.left.Value(), n.right.Value())
	}
	panic("unreachable")
}

// Remove deletes the range [start:end) (exclusive bound) from the rope.
func (n *Node) Remove(start, end int, res []rune) []rune {
	switch n.kind {
	case tLeaf:
		// slice tricks delete
		res = append(res, n.value[start:end]...)
		n.value = append(n.value[:start], n.value[end:]...)
		n.length = len(n.value)
		return res
	case tNode:
		leftLength := n.left.length
		leftStart := min(start, leftLength)
		leftEnd := min(end, leftLength)
		rightLength := n.right.length
		rightStart := max(0, min(start-leftLength, rightLength))
		rightEnd := max(0, min(end-leftLength, rightLength))
		if leftStart < leftLength {
			return n.left.Remove(leftStart, leftEnd, res)
		}
		if rightEnd > 0 {
			return n.right.Remove(rightStart, rightEnd, res)
		}
		n.length = n.left.length + n.right.length
	}
	n.adjust()
	return res
}

// Insert inserts the given value at pos.
func (n *Node) Insert(pos int, value []rune) {
	switch n.kind {
	case tLeaf:
		// slice tricks insert
		n.value = insert(n.value, pos, value)
		n.length = len(n.value)
	case tNode:
		leftLength := n.left.length
		if pos < leftLength {
			n.left.Insert(pos, value)
		} else {
			n.right.Insert(pos-leftLength, value)
		}
		n.length = n.left.length + n.right.length
	}
	n.adjust()
}

// Slice returns the range of the rope from [start:end). The returned slice
// is not copied.
func (n *Node) Slice(start, end int) []rune {
	if start >= end {
		return []rune{}
	}

	switch n.kind {
	case tLeaf:
		return n.value[start:end]
	case tNode:
		leftLength := n.left.length
		leftStart := min(start, leftLength)
		leftEnd := min(end, leftLength)
		rightLength := n.right.length
		rightStart := max(0, min(start-leftLength, rightLength))
		rightEnd := max(0, min(end-leftLength, rightLength))

		if leftStart != leftEnd {
			if rightStart != rightEnd {
				return concat(n.left.Slice(leftStart, leftEnd), n.right.Slice(rightStart, rightEnd))
			} else {
				return n.left.Slice(leftStart, leftEnd)
			}
		} else {
			if rightStart != rightEnd {
				return n.right.Slice(rightStart, rightEnd)
			} else {
				return []rune{}
			}
		}
	}
	panic("unreachable")
}

// At returns the element at the given position.
func (n *Node) At(pos int) rune {
	s := n.Slice(pos, pos+1)
	return s[0]
}

// SplitAt splits the node at the given index and returns two new ropes
// corresponding to the left and right portions of the split.
func (n *Node) SplitAt(i int) (*Node, *Node) {
	switch n.kind {
	case tLeaf:
		return New(n.value[:i]), New(n.value[i:])
	case tNode:
		m := n.left.length
		if i == m {
			return n.left, n.right
		} else if i < m {
			l, r := n.left.SplitAt(i)
			return l, join(r, n.right)
		}
		l, r := n.right.SplitAt(i - m)
		return join(n.left, l), r
	}
	panic("unreachable")
}

func join(l, r *Node) *Node {
	n := &Node{
		left:   l,
		right:  r,
		length: l.length + r.length,
		kind:   tNode,
	}
	n.adjust()
	return n
}

// Join merges all the given ropes together into one rope.
func Join(a, b *Node, more ...*Node) *Node {
	s := join(a, b)
	for _, n := range more {
		s = join(s, n)
	}
	return s
}

// Rebuild rebuilds the entire rope structure, resulting in a balanced tree.
func (n *Node) Rebuild() {
	switch n.kind {
	case tNode:
		n.value = concat(n.left.Value(), n.right.Value())
		n.left = nil
		n.right = nil
		n.adjust()
	}
}

// Rebalance finds unbalanced nodes and rebuilds them.
func (n *Node) Rebalance() {
	switch n.kind {
	case tNode:
		lratio := float64(n.left.length) / float64(n.right.length)
		rratio := float64(n.right.length) / float64(n.left.length)
		if lratio > RebalanceRatio || rratio > RebalanceRatio {
			n.Rebuild()
		} else {
			n.left.Rebalance()
			n.right.Rebalance()
		}
	}
}

// Each applies the given function to every node in the rope.
func (n *Node) Each(fn func(n *Node)) {
	fn(n)
	if n.kind == tNode {
		n.left.Each(fn)
		n.right.Each(fn)
	}
}

// EachLeaf applies the given function to every leaf node in order.
func (n *Node) EachLeaf(fn func(n *Node) bool) bool {
	switch n.kind {
	case tLeaf:
		return fn(n)
	default: // case tNode
		if n.left.EachLeaf(fn) {
			return true
		}
		return n.right.EachLeaf(fn)
	}
}

// ReadAt implements the io.ReaderAt interface.
func (n *Node) ReadAt(p []rune, off int64) (nread int, err error) {
	if off > int64(n.length) {
		return 0, io.EOF
	}

	end := off + int64(len(p))
	if end >= int64(n.length) {
		end = int64(n.length)
		err = io.EOF
	}
	b := n.Slice(int(off), int(end))
	nread = copy(p, b)
	return nread, err
}

/*
// WriteTo implements the io.WriterTo interface.
func (n *Node) WriteTo(w io.Writer) (int64, error) {
	var err error
	var ntotal int64
	n.EachLeaf(func(it *Node) bool {
		var nwritten int
		nwritten, err = w.Write(it.Value())
		ntotal += int64(nwritten)
		return err != nil
	})
	return ntotal, err
}
*/

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// from slice tricks
func insert(s []rune, k int, vs []rune) []rune {
	if n := len(s) + len(vs); n <= cap(s) {
		s2 := s[:n]
		copy(s2[k+len(vs):], s[k:])
		copy(s2[k:], vs)
		return s2
	}
	s2 := make([]rune, len(s)+len(vs))
	copy(s2, s[:k])
	copy(s2[k:], vs)
	copy(s2[k+len(vs):], s[k:])
	return s2
}

func concat(a, b []rune) []rune {
	c := make([]rune, 0, len(a)+len(b))
	c = append(c, a...)
	c = append(c, b...)
	return c
}
