package main

import (
	"fmt"
	"testing"
)

func TestInsertA(t *testing.T) {
	t.Log("open file")
	f := FileWrapper{
		Filename: "testlist.md",
		CharLen:  0,
	}
	gb := GapBuffer{}
	gb.Load(&f)
	fmt.Println(gb.ToString())

	fmt.Println(gb.GapLen)

	fmt.Println()
	fmt.Println()
	err := gb.CursorTo(2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(gb.ToString())
	fmt.Println()
	fmt.Println()

	fmt.Println(gb.Length())
	fmt.Println(gb.GapLen)
	gb.insertAtGap([]rune{'g', 'a', 'p', ' ', 'b', 'u', 'f', 'f', 'e', 'r', ' ', 'a', 'n', 'd'})
	fmt.Println(gb.ToString())
	fmt.Println()
	fmt.Println()

	fmt.Println(gb.Length())
	gb.Append([]rune{'w', 'a', 'w'})
	fmt.Println(gb.ToString())

	gb.DeleteRange(0, 2)
	fmt.Println(gb.ToString())
	// todo: remove, remove block, replace, split

}
