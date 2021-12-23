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

func (g *GapBuffer) isEmpty() error {
	if g.Length() == 0 {
		return errors.New("list is empty")
	}
	return nil
}

func (g *GapBuffer) Length() int {
	return g.Content.getLen() - g.GapLen
}

func (g *GapBuffer) checkIdxStrict(i ...int) error {
	err := g.isEmpty()
	if err != nil {
		return err
	}
	prev := i[0]
	for _, el := range i {
		if el < prev {
			return errors.New("incorrect index order")
		}
		if (el < 0) || (g.Length() <= el) {
			return errors.New("Index Out of bounds")
		}
		prev = el
	}
	return nil
}

func (g *GapBuffer) checkIdxLoose(i ...int) error {
	// for anything that allows something to be put at the very end
	err := g.isEmpty()
	if err != nil {
		return err
	}
	prev := i[0]
	for _, el := range i {
		if el < prev {
			return errors.New("incorrect index order")
		}
		if (el < 0) || (g.Length() < el) {
			return errors.New("Index Out of bounds")
		}
		prev = el
	}
	return nil
}

func (g *GapBuffer) getNode(i int) (*DoubleLink, error) {
	err := g.checkIdxStrict(i)
	if err != nil {
		return nil, err
	}
	if i < g.GapStartIdx {
		res, err := g.Content.TraverseTo(i)
		if err != nil {
			return nil, err
		}
		return res, nil
	} else {
		res, err := g.Content.TraverseTo(i + g.GapLen)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (g *GapBuffer) Report() ([]rune, error) {
	// just special case of ReportRange where i = 0 and j = length - 1
	return g.ReportRange(0, g.Length()-1)
}

func (g *GapBuffer) ReportRange(i, j int) ([]rune, error) {
	err := g.checkIdxStrict(i, j)

	if err != nil {
		return nil, err
	}
	if i == j {
		return nil, errors.New("use ReportCharacter instead")
	}
	ret := make([]rune, 0)
	if g.GapStartIdx == 0 {
		i += g.GapLen
	}
	cn, err := g.Content.TraverseTo(i)
	if err != nil {
		return nil, err
	}
	idx := i
	if j >= g.GapStartIdx {
		j += g.GapLen
	}
	for idx <= j {
		if idx >= g.GapStartIdx+g.GapLen || idx < g.GapStartIdx {
			ret = append(ret, cn.Content)
		}
		idx++
		cn = cn.Next
	}
	return ret, nil

}

func (g *GapBuffer) ReportCharacter(i int) (rune, error) {
	err := g.checkIdxStrict(i)
	if err != nil {
		return -1, err
	}

	d, err := g.getNode(i)
	if err != nil {
		return -1, err
	}
	return d.Content, nil
}

func (g *GapBuffer) insertAtGap(content []rune) error {

	if len(content) >= g.GapLen {
		es := len(content) - g.GapLen
		g.expandGap(es + 2)
	}
	cn := g.GapStart
	for _, r := range content {
		g.GapStart = g.GapStart.Next
		cn.Content = r
		if cn.Next == nil {
			return errors.New("bad next link while inserting in gap")
		}
		cn = cn.Next
		g.GapLen--
		g.GapStartIdx++
	}
	return nil

}

func (g *GapBuffer) expandGap(size int) {
	for i := 0; i < size; i++ {
		g.Content.Insert(g.GapStartIdx+g.GapLen, -1) // todo: this is probably the part you want to tune for indices
		g.GapEnd = g.GapEnd.Next
		g.GapLen++
	}
}

func (g *GapBuffer) Insert(i int, content []rune) error {
	err := g.moveGap(i)
	if err != nil {
		return err
	}
	g.insertAtGap(content)
	return nil
}

func (g *GapBuffer) Append(content []rune) error {
	err := g.moveGap(g.Length())
	if err != nil {
		return err
	}
	g.insertAtGap(content)
	return nil
}

func (g *GapBuffer) moveGap(i int) error {
	if i == g.GapStartIdx {
		return nil
	}
	err := g.checkIdxLoose(i)
	if err != nil {
		return err
	}
	var prevIdx, nxt *DoubleLink

	if i != 0 && i != g.Length() {
		prevIdx, err = g.getNode(i - 1)
		if err != nil {
			return err
		}
		nxt = prevIdx.Next
	}

	if g.GapStartIdx == 0 {
		g.Content.Head = g.GapEnd.Next
		g.GapEnd.DelinkR()
	} else if g.GapStartIdx+g.GapLen == g.Content.Length {
		g.Content.End = g.GapStart.Prior
		g.GapStart.DelinkL()
	} else {
		bf := g.GapStart.Prior
		af := g.GapEnd.Next
		g.GapStart.DelinkL()
		g.GapEnd.DelinkR()
		bf.LinkR(af)
	}

	if i == 0 {
		g.GapEnd.LinkR(g.Content.Head)
		g.Content.Head = g.GapStart
	} else if i == g.Length() {
		g.GapStart.LinkL(g.Content.End)
		g.Content.End = g.GapEnd
	} else {
		prevIdx.DelinkR()
		g.GapStart.LinkL(prevIdx)
		g.GapEnd.LinkR(nxt)
	}
	g.GapStartIdx = i
	return nil
}

func (g *GapBuffer) Replace(i int, content []rune) error {
	err := g.checkIdxStrict(i)
	if err != nil {
		return err
	}

	start, err := g.getNode(i)
	if err != nil {
		return err
	}
	for _, r := range content {
		if start == g.GapStart {
			for start != g.GapEnd {
				start = start.Next
			}
			start = start.Next
		}
		start.Content = r
		start = start.Next
		if start == nil {
			return errors.New("overflow list")
		}

	}
	return nil
}

func (g *GapBuffer) compressGap() error {
	if g.GapLen == 0 {
		return nil
	}
	if g.GapStartIdx == 0 {
		g.Content.Head = g.GapEnd.Next
		g.Content.Head.DelinkL()
	} else if g.GapStartIdx+g.GapLen == g.Content.Length {
		g.Content.End = g.GapStart.Prior
		g.Content.End.DelinkR()
	} else {
		bf := g.GapStart.Prior
		af := g.GapEnd.Next
		g.GapStart.DelinkL()
		g.GapEnd.DelinkR()
		bf.LinkR(af)
	}
	g.Content.Length -= g.GapLen
	g.GapStart, g.GapEnd = nil, nil
	g.GapStartIdx, g.GapLen = 0, 0
	return nil
}

func (g *GapBuffer) Split(i int) (StorageType, error) {
	err := g.checkIdxStrict(i)
	if err != nil || i == 0 {
		return nil, errors.New("bad split point")
	}

	err = g.compressGap()
	if err != nil {
		return nil, err
	}
	s, err := g.Content.Split(i)
	if err != nil {
		return nil, err
	}
	nb := &GapBuffer{
		Content:     s,
		GapStart:    s.Head,
		GapEnd:      s.Head,
		GapStartIdx: 0,
		GapLen:      0,
	}
	err = nb.makeGap()
	if err != nil {
		return nil, err
	}
	err = g.makeGap()
	if err != nil {
		return nil, err
	}
	return nb, nil
}

func (g *GapBuffer) DeleteRange(i, j int) ([]rune, error) {
	err := g.checkIdxStrict(i, j)
	if err != nil {
		return nil, err
	}

	ret := make([]rune, 0)
	inCt := j - i
	g.moveGap(i)
	for i := 0; i < inCt; i++ { // todo: tune
		g.GapEnd = g.GapEnd.Next
		ret = append(ret, g.GapEnd.Content)
		g.GapEnd.Content = -1
		g.GapLen++
	}
	return ret, nil
}

func (g *GapBuffer) Concat(content []rune) error {
	err := g.Append(content)
	return err
}

func (g *GapBuffer) Save(f *FileWrapper) error {
	return nil
}

func (g *GapBuffer) Load(f *FileWrapper, ct int) error {
	bs := f.Open(ct)
	g.Content = &DoublyLinkedList{
		Length: 0,
	}
	g.GapLen = 0
	for _, b := range bs {
		err := g.Content.Insert(g.Length(), rune(b))
		if err != nil {
			return err
		}
	}
	err := g.makeGap()
	return err
}

func (g *GapBuffer) makeGap() error {
	if g.GapLen > 0 {
		return nil
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
	ct := 0
	for head != nil {
		if head.Content == -1 { //(ct >= g.GapStartIdx) && (ct <= g.GapStartIdx+g.GapLen)
			ret += "_"
		} else {
			ret += string(head.Content)
		}
		ct++
		head = head.Next
	}
	return ret
}
