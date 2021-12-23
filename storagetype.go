package main

import (
	"io"
	"log"
	"os"
)

/*
type StorageType interface {
	PopContent(x1, y1, x2, y2 int) ([]rune, time.Time)
	CreateBuffer(initial [][]rune) time.Time
	GetCurrentSelection() ([4]int, time.Time) // x1 y1, x2 y2
	SetSelection(x1, y1, x2, y2 int) time.Time
	Insert(content []rune, startX, startY int) time.Time
}
*/

type StorageType interface {
	Report() ([]rune, error)              // report entire buffer
	ReportRange(i, j int) ([]rune, error) // report segment of buffer from i ... j, inclusive
	// ReportCharacter(i int) (rune, error)  // report single character
	Insert(i int, content []rune) error // insert content at position i.
	Append(content []rune) error        // used during testing, ignore
	//	Replace(i int, content []rune) error
	Split(i int) (StorageType, error)     // splits at point i, returning a new buffer [i ... n]
	DeleteRange(i, j int) ([]rune, error) // deletes from [i, j)
	//Concat(content []rune) error
	// Save(f *FileWrapper) error
	Load(f *FileWrapper, ct int) error // loads ``ct`` characters in from file f`
	ToString() string                  // self-explanatory
}

/*
	Report: only scaled with file length
	Load: only scaled with file length


	ReportRange: constant file length, scale reported length and depth
	Insert: different filen lengths, scale position and inserted length
	Split: scale file length and split point
	DeleteRange: scale file length, start, and length


*/

func FlattenXY(x, y, linelen int) (int, error) {
	return -1, nil
}

// ask gesell about how to do make, tree

type FileWrapper struct {
	Filename string
}

func (f *FileWrapper) Open(ct int) []byte {
	// Open file for reading
	file, err := os.Open(f.Filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	// The file.Read() function will happily read a tiny file in to a large
	// byte slice, but io.ReadFull() will return an
	// error if the file is smaller than the byte slice.
	byteSlice := make([]byte, ct)
	_, err = io.ReadFull(file, byteSlice)
	if err != nil {
		log.Fatal(err)
	}
	return byteSlice
}
