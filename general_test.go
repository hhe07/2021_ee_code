package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func GetContent(filename string, count int) []byte {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ret := make([]byte, count)
	_, err = io.ReadFull(file, ret)
	if err != nil {
		panic(err)
	}
	return ret
}

type BufferTest struct {
	Start  []int // start also serves as the index in single-index operations
	End    []int // if this is -1, represents end
	Length int   // for operations that need some input, this will be used for additional loads
}

func TestLoad(t *testing.T) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := []BufferTest{
		{Length: 100},
		{Length: 200},
		{Length: 500},
		{Length: 1000},
	}
	for _, tc := range testCases {
		content := GetContent("testing.txt", tc.Length)
		for _, b := range buffers {
			t.Run(fmt.Sprintf("%T_load", b), func(t *testing.T) {
				t.Log(fmt.Sprintf("case: %d", tc.Length))
				b.Load(content)
				if b.ToString() != string(content) {
					t.Fatal("improperly loaded string")

				}
			})
		}
	}
}

func TestInsert(t *testing.T) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()
	testCases := []BufferTest{
		{Start: []int{0, 25, 50, -1}, Length: 100},
		{Start: []int{0, 50, 100, -1}, Length: 200},
		{Start: []int{0, 125, 250, -1}, Length: 500},
		{Start: []int{0, 250, 500, -1}, Length: 1000},
	}

	for _, tc := range testCases {
		content := GetContent("testing.txt", tc.Length)
		app := make([]rune, len(content))
		for i, bt := range content {
			app[i] = rune(bt)
		}

		//comparison := string(app)
		//comparison = comparison + comparison

		for _, b := range buffers {
			for sc_idx := range tc.Start {
				b.Load(content)
				point := tc.Start[sc_idx]
				if point == -1 {
					point = tc.Length - 1
				}
				t.Run(fmt.Sprintf("%T_insert", b), func(t *testing.T) {
					t.Log(fmt.Sprintf("case: %d", point))

					comparison := string(append(app[:point], append(app, app[point:]...)...))

					err := b.Insert(point, app)
					if err != nil {
						t.Fatal(err)
					}
					res, err := b.Report()
					if err != nil {
						t.Fatal(err)
					}
					if string(res) != comparison {
						t.Fatal(fmt.Sprintf("improperly appended \n expected: %s\n got: %s \n", comparison, string(res)))

					}
				})
			}

		}
	}
}

func TestReport(t *testing.T) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := []BufferTest{
		{Length: 100},
		{Length: 200},
		{Length: 500},
		{Length: 1000},
	}
	for _, tc := range testCases {
		content := GetContent("testing.txt", tc.Length)
		for _, b := range buffers {
			b.Load(content)
			t.Run(fmt.Sprintf("%T_report", b), func(t *testing.T) {
				t.Log(fmt.Sprintf("case: %d", tc.Length))
				res, err := b.Report()
				if err != nil {
					t.Fatal(err)
				}
				if string(res) != string(content) {
					t.Fatal("improperly reported contents")

				}
			})
		}
	}
}

func TestReportRange(t *testing.T) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := []BufferTest{
		{Start: []int{0, 0, 50, 10, 30}, End: []int{-1, 10, -1, 20, 50}, Length: 100},
		{Start: []int{0, 0, 50, 10, 30}, End: []int{-1, 10, -1, 20, 50}, Length: 200},
		{Start: []int{0, 0, 50, 10, 30}, End: []int{-1, 10, -1, 20, 50}, Length: 500},
		{Start: []int{0, 0, 50, 10, 30}, End: []int{-1, 10, -1, 20, 50}, Length: 1000},
	}
	for _, tc := range testCases {
		content := GetContent("testing.txt", tc.Length)
		for _, b := range buffers {
			b.Load(content)
			for sc_idx := range tc.Start {
				t.Run(fmt.Sprintf("%T_report_range", b), func(t *testing.T) {
					start, end := tc.Start[sc_idx], tc.End[sc_idx]
					if end == -1 {
						end = len(content) - 1
					}
					t.Log(fmt.Sprintf("len: %d, range: %d:%d", tc.Length, start, end))
					res, err := b.ReportRange(start, end)
					if err != nil {
						t.Fatal(err)
					}
					if string(res) != string(content[start:end+1]) {
						t.Fatal(fmt.Sprintf("improperly reported \n expected: %s\n got: %s \n", string(content[start:end+1]), string(res)))

					}
				})
			}

		}
	}
}

func TestAppend(t *testing.T) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := []BufferTest{
		{Length: 100},
		{Length: 200},
		{Length: 500},
		{Length: 1000},
	} // basically, starting with this len in buffer, and then doubling size with same content
	for _, tc := range testCases {
		content := GetContent("testing.txt", tc.Length)
		app := make([]rune, len(content))
		for i, bt := range content {
			app[i] = rune(bt)
		}
		comparison := string(app)
		comparison = comparison + comparison
		for _, b := range buffers {
			b.Load(content)
			t.Run(fmt.Sprintf("%T_append", b), func(t *testing.T) {
				t.Log(fmt.Sprintf("case: %d", tc.Length))
				t.Log(b.ToString())

				err := b.Append(app)
				if err != nil {
					t.Fatal(err)
				}
				res, err := b.Report()
				if err != nil {
					t.Fatal(err)
				}
				if string(res) != comparison {
					t.Fatal(fmt.Sprintf("improperly appended \n expected: %s\n got: %s \n", comparison, string(res)))

				}
			})
		}
	}

}

func TestSplit(t *testing.T) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := []BufferTest{
		{End: []int{25, 50}, Length: 100},
		{End: []int{50, 100}, Length: 200},
		{End: []int{125, 250}, Length: 500},
		{End: []int{250, 500}, Length: 1000},
	}

	for _, tc := range testCases {
		content := GetContent("testing.txt", tc.Length)
		for _, b := range buffers {
			b.Load(content)
			for sc_idx := range tc.End {
				t.Run(fmt.Sprintf("%T_split", b), func(t *testing.T) {
					point := tc.End[sc_idx]
					if point == -1 {
						point = len(content) - 1
					}
					t.Log(fmt.Sprintf("split: %d", point))
					res, err := b.Split(point)
					if err != nil {
						t.Fatal(err)
					}
					splitStr, err := res.Report()
					if err != nil {
						t.Fatal(err)
					}

					if string(splitStr) != string(content[point:]) {
						t.Log(len(splitStr))
						t.Log(len(content[point:]))
						t.Fatal(fmt.Sprintf("bad split off section\n expected: /%s/ \n got: /%s/ \n", content[point:], string(splitStr)))
					}
					remainStr, err := b.Report()
					if err != nil {
						t.Fatal(err)
					}
					if string(remainStr) != string(content[:point]) {
						t.Fatal(fmt.Sprintf("bad remaining section\n expected: /%s/ \n got: /%s/ \n", content[:point], string(remainStr)))
					}
					// cleanup
					//t.Log(string(splitStr))
					//t.Log(b.ToString())
					b.Append(splitStr)

				})
			}

		}
	}
}

func TestDelete(t *testing.T) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := []BufferTest{
		{Start: []int{0, 0, 50, 10, 30}, End: []int{-1, 10, -1, 20, 50}, Length: 100},
		{Start: []int{0, 0, 50, 10, 30}, End: []int{-1, 10, -1, 20, 50}, Length: 200},
		{Start: []int{0, 0, 50, 10, 30}, End: []int{-1, 10, -1, 20, 50}, Length: 500},
		{Start: []int{0, 0, 50, 10, 30}, End: []int{-1, 10, -1, 20, 50}, Length: 1000},
	}
	for _, tc := range testCases {
		content := GetContent("testing.txt", tc.Length)
		for _, b := range buffers {
			b.Load(content)
			for sc_idx := range tc.Start {
				t.Run(fmt.Sprintf("%T_Delete", b), func(t *testing.T) {
					start, end := tc.Start[sc_idx], tc.End[sc_idx]
					if end == -1 {
						end = len(content) - 1
					}
					t.Log(fmt.Sprintf("len: %d, range: %d:%d", tc.Length, start, end))
					res, err := b.DeleteRange(start, end)
					if err != nil {
						t.Fatal(err)
					}
					comp := content[start:end]
					if string(res) != string(comp) {
						t.Fatal(fmt.Sprintf("improperly deleted section \n expected: %s\n got: %s \n", string(comp), string(res)))

					}

					rem, err := b.Report()
					if err != nil {
						t.Fatal(err)
					}
					comp = make([]byte, len(content))
					copy(comp, content)

					comp = append(comp[:start], comp[end:]...)
					if string(rem) != string(comp) {
						t.Fatal(fmt.Sprintf("improper remaining section \n expected: %s\n got: %s \n", string(comp), string(rem)))

					}
					b.Insert(start, res)
				})
			}

		}
	}
}

func TestIndex(t *testing.T) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := []BufferTest{
		{Start: []int{0, 25, 50, 75, 99}, Length: 100},
		{Start: []int{0, 50, 100, 150, 199}, Length: 200},
		{Start: []int{0, 125, 250, 375, 499}, Length: 500},
		{Start: []int{0, 250, 500, 750, 999}, Length: 1000},
	}
	for _, tc := range testCases {
		content := GetContent("testing.txt", tc.Length)
		for _, b := range buffers {
			b.Load(content)
			for _, idx := range tc.Start {
				res, err := b.Index(idx)
				if err != nil {
					t.Fatal(err)
				}
				if rune(content[idx]) != res {
					t.Fatal(fmt.Sprintf("character incorrectly reported \n expected: %s \n got: %s", string(content[idx]), string(res)))
				}
			}
		}
	}
}

func TestConcat(t *testing.T) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	extras := make([]StorageType, 2)
	extras[0] = mkRope()
	extras[1] = mkGapBuf()

	testCases := []BufferTest{
		{Start: []int{25, 50, 75, 99}, Length: 100},
		{Start: []int{50, 100, 150, 199}, Length: 200},
		{Start: []int{125, 250, 375, 499}, Length: 500},
		{Start: []int{250, 500, 750, 999}, Length: 1000},
	}
	for _, tc := range testCases {
		content := GetContent("testing.txt", tc.Length)
		for i, b := range buffers {
			for _, idx := range tc.Start {
				b.Load(content)
				extras[i].Load(content[:idx])

				err := b.Concat(extras[i])
				if err != nil {
					t.Fatal(err)
				}
				expected := make([]byte, len(content))
				copy(expected, content)
				expected = append(expected, content[:idx]...)

				res, err := b.Report()
				if err != nil {
					t.Fatal(err)
				}

				if string(expected) != string(res) {
					t.Fatal(fmt.Sprintf("bad append.\n expected: %s\n got: %s\n", string(expected), string(res)))
				}

			}
		}
	}
}
