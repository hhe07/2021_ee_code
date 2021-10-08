package main

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	t.Log("Traverse To test")

	ll := DoublyLinkedList{Head: nil, End: nil, Length: 0}
	ll.Insert(0, 'h')
	ll.Insert(1, 'a')

	ll.Insert(1, 'b')
	res, _ := ll.TraverseTo(1)
	fmt.Println(string(res.Content))
	fmt.Println(ll.ToString())
	ll.Insert(0, 'c')
	fmt.Println(ll.ToString())
	//ll.Swap(1, 2)
	//fmt.Println(ll.ToString())



}
