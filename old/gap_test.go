package main

import (
	"testing"
)

func TestInit(t *testing.T) {
	gb := GapBuffer{
		Content: &DoublyLinkedList{
			Head:   nil,
			End:    nil,
			Length: 0,
		},
		GapStart:    nil,
		GapEnd:      nil,
		GapStartIdx: 0,
		GapLen:      0,
	}
	gb.makeGap()
	t.Log(gb.ToString())

	/*
		Insert at a non-zero position
		doesn't work with a newly inited gap buffer,
		but for the sake of this project that's fine.
	*/
	t.Logf("GapLen: %d, TextLen: %d", gb.GapLen, gb.Length())
	gb.Insert(0, []rune{'a'})
	t.Log(gb.ToString())
	// works
	t.Logf("GapLen: %d, TextLen: %d", gb.GapLen, gb.Length())

	// testing an insert at the start
	err := gb.Insert(0, []rune{'b'})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())
	// testing an insert at the end
	err = gb.Insert(2, []rune{'c'})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())
	// testing an insert at the centre
	err = gb.Insert(1, []rune{'c'})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())
	// getting some nodes
	node, err := gb.getNode(0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(node.Content)) // b

	node, err = gb.getNode(3)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(node.Content)) // c

	node, err = gb.getNode(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(node.Content)) // c

	node, err = gb.getNode(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(node.Content)) // a
	// testing invalid indices
	_, err = gb.getNode(-1)
	if err == nil {
		t.Fatal("bad index fell through")
	}
	_, err = gb.getNode(100)
	if err == nil {
		t.Fatal("bad index fell through")
	}
	t.Log(gb.GapStartIdx + gb.GapLen)
	// testing Report method
	r, err := gb.Report()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(r))

	// testing ReportRange method
	r, err = gb.ReportRange(0, 1) // bc
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(r))

	r, err = gb.ReportRange(1, 2) // ca
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(r))

	r, err = gb.ReportRange(1, 3) // cac
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(r))

	_, err = gb.ReportRange(2, 1) // fail
	if err == nil {
		t.Fatal("uncaught report error")
	}

	_, err = gb.ReportRange(1, 1) // fail
	if err == nil {
		t.Fatal("uncaught report error")
	}

	_, err = gb.ReportRange(-1, 0) // fail
	if err == nil {
		t.Fatal("uncaught report error")
	}

	_, err = gb.ReportRange(0, 5) // fail
	if err == nil {
		t.Fatal("uncaught report error")
	}

	// testing ReportCharacter method
	c, err := gb.ReportCharacter(0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(c))

	c, err = gb.ReportCharacter(3)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(c))

	c, err = gb.ReportCharacter(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(c))

	c, err = gb.ReportCharacter(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(c))

	err = gb.moveGap(0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())

	c, err = gb.ReportCharacter(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(c)) // a

	// testing insertAtGap method

	err = gb.insertAtGap([]rune{'h', 'e', 'l', 'l', 'o', 'w', 'o', 'r', 'l', 'd'})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())

	// testing Append method

	err = gb.Append([]rune{'a', 'a', 'a'})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())
	// moveGap has been implicitly tested already

	// testing Replace method

	err = gb.moveGap(10)
	if err != nil {
		t.Fatal(err)
	}

	err = gb.Replace(10, []rune{'a', 'a', 'a'})
	if err != nil {
		t.Fatal(err)
	}

	// testing DeleteRange method
	r, err = gb.DeleteRange(10, 15) // excludes last
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())
	t.Log(string(r))

	// testing Split method
	ngb, err := gb.Split(10) // new split off section includes the element
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())
	t.Log(ngb.ToString())

}
