package main

import (
	"testing"
)

func TestRope(t *testing.T) {
	gb := Rope{Root: New([]rune{'a'})}
	res, err := gb.Report()
	t.Log(string(res))

	err = gb.Insert(0, []rune{'b'})
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

	// testing ReportRange method
	r, err := gb.ReportRange(0, 1) // bc
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
	/*

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
	*/

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

	c, err = gb.ReportCharacter(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(c)) // a

	// testing Append method

	err = gb.Append([]rune{'a', 'a', 'a'})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())

	// testing DeleteRange method
	r, err = gb.DeleteRange(1, 3)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())
	t.Log(string(r))

	// testing Split method
	ngb, err := gb.Split(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gb.ToString())
	t.Log(ngb.ToString())

}
