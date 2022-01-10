package main

type StorageType interface {
	Report() ([]rune, error)              // report entire buffer
	ReportRange(i, j int) ([]rune, error) // report segment of buffer from i ... j, inclusive
	Insert(i int, content []rune) error   // insert content at position i.
	Append(content []rune) error          // used during testing, ignore
	Concat(s StorageType) error           // self-explanatory
	Split(i int) (StorageType, error)     // splits at point i, returning a new buffer [i ... n]
	DeleteRange(i, j int) ([]rune, error) // deletes from [i, j)
	Load(contents []byte) error           // loads contents into buffer
	ToString() string                     // self-explanatory
	Index(i int) (rune, error)            //  returns character at position
}
