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
	ReportRange(i, j int) ([]rune, error) // report segment of buffer
	ReportCharacter(i int) (rune, error)  // report single character
	Insert(i int, content []rune) error
	Append(i int, content []rune) error
	Replace(i int, content []rune) error
	Split(i int) (StorageType, error) // TODO: will this accept both rope and gapbuffer?
	DeleteRange(i, j int) ([]rune, error)
	Concat(content []rune) error
	Save(f *FileWrapper) error
	IsReady() (bool, error)
	Load(f *FileWrapper) error
}

func FlattenXY(x, y, linelen int) (int, error) {
	return -1, nil
}

// ask gesell about how to do make, tree shit

type FileWrapper struct {
	Filename string
	CharLen  int
}

func (f *FileWrapper) Open() []byte {
	// Open file for reading
	file, err := os.Open(f.Filename)
	if err != nil {
		log.Fatal(err)
	}

	// The file.Read() function will happily read a tiny file in to a large
	// byte slice, but io.ReadFull() will return an
	// error if the file is smaller than the byte slice.
	byteSlice := make([]byte, 200)
	_, err = io.ReadFull(file, byteSlice)
	if err != nil {
		log.Fatal(err)
	}
	return byteSlice
}
