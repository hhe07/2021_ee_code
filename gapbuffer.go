package main

import (
	"errors"
)

type GapBuffer struct {
	Content     *DoublyLinkedList
	GapStart    *DoubleLink
	GapEnd      *DoubleLink
	GapStartIdx int
	GapLen      int
}

// todo: unsafe: requires that Content be constructed without checking that it
// todo:  actually is

func (g *GapBuffer) Length() int {
	return g.Content.getLen() - g.GapLen
}

func (g *GapBuffer) Report() ([]rune, error) {
	// just special case of ReportRange where i = 0 and j = length - 1
	return g.ReportRange(0, g.Length()-1)
}

func (g *GapBuffer) getNode(i int) (*DoubleLink, error) {
	return g.Content.TraverseTo(i)
}

func (g *GapBuffer) ReportRange(i, j int) ([]rune, error) {
	if j > g.Length()-1 {
		return nil, errors.New("exceeds List Len")
	}
	if i < 0 {
		return nil, errors.New("index less than 0")
	}
	ret := make([]rune, 0)
	cn, err := g.getNode(i)
	if err != nil {
		return nil, err
	}
	inGap := false
	for idx := i; idx < j; idx++ {
		if cn == g.GapStart {
			inGap = true
		}
		if !inGap {
			ret = append(ret, cn.Content)
		}
		if cn == g.GapEnd {
			inGap = false
		}
		cn = cn.Next
	}
	return ret, nil

}

func (g *GapBuffer) ReportCharacter(i int) (rune, error) {
	if i < 0 || i > g.Length()-1 {
		return -1, errors.New("Index exceeds internal list constraints")
	}
	d, err := g.getNode(i)
	if err != nil {
		return -1, errors.New("bad internal traverse in reportCharacter")
	}
	return d.Content, nil
}

func (g *GapBuffer) insertAtGap(content []rune) error {
	if len(content) > g.GapLen {
		es := len(content) - g.GapLen
		g.expandGap(es + 2)
	}
	cn := g.GapStart
	for _, r := range content {
		g.GapStart = g.GapStart.Next
		cn.Content = r
		if cn.Next == nil {
			return errors.New("bad next item while inserting in gap")
		}
		cn = cn.Next
		g.GapLen--
		g.GapStartIdx++
	}
	return nil

}

func (g *GapBuffer) expandGap(size int) {
	for i := 0; i < size; i++ {

		currEnd := g.GapEnd
		nl := &DoubleLink{
			Content: -1,
			Next:    currEnd.Next,
			Prior:   currEnd,
		}
		currEnd.Next = nl
		g.GapEnd = nl
		g.GapLen++
		g.Content.Length++
	}
}

func (g *GapBuffer) Insert(i int, content []rune) error {
	err := g.CursorTo(i)
	if err != nil {
		return err
	}
	g.insertAtGap(content)
	return nil
}

func (g *GapBuffer) Append(content []rune) error {
	err := g.CursorTo(g.Length())
	if err != nil {
		return err
	}
	g.insertAtGap(content)
	// todo: make this a special case that directly puts stuff onto end rather than
	// todo: moving the gap?
	return nil
}

func (g *GapBuffer) Replace(i int, content []rune) error {
	// todo: should this move the gap?
	start, err := g.getNode(i)
	if err != nil {
		return err
	}
	for _, r := range content {
		start.Content = r
		start = start.Next
		if start == nil {
			return errors.New("overflow list")
		}
	}
	return nil
}

func (g *GapBuffer) Split(i int) (*GapBuffer, error) {
	// todo: should init a gap
	splitPt := i
	if i > g.GapStartIdx {
		splitPt += g.GapLen // todo: test
	}
	s, err := g.Content.Split(i)
	if err != nil {
		return nil, err
	}
	return &GapBuffer{
		Content:     s,
		GapStart:    s.Head,
		GapEnd:      s.Head,
		GapStartIdx: 0,
		GapLen:      0,
	}, nil
}

func (g *GapBuffer) DeleteRange(i, j int) ([]rune, error) {
	// TODO: single character removal doesn't account for gap
	if i == j {
		r, err := g.Content.Remove(i)
		return []rune{r}, err
	}
	delStart, err := g.getNode(i) // included: at position i-1
	if err != nil {
		return nil, err
	}
	delEnd, err := g.getNode(j) // excluded
	if err != nil {
		return nil, err
	}
	// case: start
	if i == 0 {
		g.Content.Head = delEnd
		delEnd.Prior = nil
	} else if j == g.Length()-1 {
		// todo: fix
		g.Content.End = delStart
		g.Content.End.Next = nil
	} else {
		delStart.Prior.Next = delEnd
		delEnd.Prior = delStart.Prior
	}
	ret := make([]rune, 0)
	for delStart != nil {
		ret = append(ret, delStart.Content)
		delStart = delStart.Next
		g.Content.Length--
	}
	return ret, nil
}

func (g *GapBuffer) CursorTo(i int) error {
	// todo: work on
	if g.GapStartIdx == i {
		return nil
	}
	// special case: if i = 0
	if i == 0 {
		tmp := g.Content.Head
		// couple prior and after gap together
		prior := g.GapStart.Prior
		end := g.GapEnd.Next
		prior.Next = end
		end.Prior = prior

		g.GapStart.Prior = nil

		g.Content.Head = g.GapStart
		g.GapEnd.Next = tmp
		tmp.Prior = g.Content.Head
		g.GapStartIdx = 0
		return nil
	} else if i >= g.Length() {
		prior := g.GapStart.Prior
		end := g.GapEnd.Next
		prior.Next = end
		end.Prior = prior

		tmp := g.Content.End
		g.GapEnd.Next = nil
		g.Content.End = g.GapEnd
		g.GapStart.Prior = tmp
		tmp.Next = g.GapStart
		g.GapStartIdx = g.Length()
		return nil
	}
	tAmt := i - g.GapStartIdx
	if tAmt > 0 {
		// move gap forwards
		for tAmt > 0 {
			mv := g.GapEnd.Next.Content // moved character
			g.GapStart.Content = mv
			g.GapStart = g.GapStart.Next
			g.GapEnd = g.GapEnd.Next
			g.GapEnd.Content = -1

			tAmt--
			g.GapStartIdx++
		}
		return nil
	}
	// move gap backwards
	for tAmt < 0 {
		mv := g.GapStart.Prior.Content // moved character
		g.GapEnd.Content = mv
		g.GapEnd = g.GapEnd.Prior
		g.GapStart = g.GapStart.Prior
		g.GapStart.Content = -1

		tAmt++
		g.GapStartIdx--
	}
	return nil

	/*
		TODO: this is good but doesnt work

		// moves the gap as a block to position i
		el, err := g.getNode(i - 1)
		if err != nil {
			return err
		}
		el.Next.Prior = g.GapEnd
		g.GapEnd.Next = el.Next

		el.Next = g.GapStart
		g.GapStart.Prior = el

		return nil

		// displaces
	*/
	// todo: shortcut for when i is already start index, start and end positions
	// todo: define behaviour well: does it displace indicated idx or otherwise

}

func (g *GapBuffer) Concat(content []rune) error {
	err := g.Append(content)
	return err
}

func (g *GapBuffer) Save(f *FileWrapper) error {
	return nil
}

func (g *GapBuffer) Load(f *FileWrapper) error {
	bs := f.Open()
	g.Content = &DoublyLinkedList{
		Length: 0,
	}
	for _, b := range bs {
		g.Content.Insert(g.Length(), rune(b))
	}
	// default: gap at very start
	err := g.Content.Insert(0, -1)
	if err != nil {
		return err
	}
	g.GapLen = 1
	g.GapStart = g.Content.Head
	g.GapEnd = g.Content.Head
	g.expandGap(9)
	g.GapStartIdx = 0
	return nil
}

func (g *GapBuffer) ToString() string {
	ret := ""
	head := g.Content.Head
	for head != nil {
		if head.Content > 0 {
			ret += string(head.Content)
		} else {
			ret += "_"
		}
		head = head.Next
	}
	return ret
}
