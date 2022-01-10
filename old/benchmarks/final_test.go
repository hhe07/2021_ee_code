package main

import (
	"fmt"
	"testing"
)

func BenchmarkGapLoad(b *testing.B) {
	gb := &GapBuffer{
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
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}
	for i := 1; i <= 5; i++ {
		size := 10000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", gb, size)
		// test loads
		b.Run(name+"_load", func(b *testing.B) {
			err := gb.Load(file, size)
			if err != nil {
				panic(err)
			}
		})
		gb.Content = &DoublyLinkedList{
			Head:   nil,
			End:    nil,
			Length: 0,
		}
		b.ResetTimer()
	}

}

func BenchmarkRopeLoad(b *testing.B) {
	r := &Rope{
		Root: nil,
	}
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}
	for i := 1; i <= 5; i++ {
		size := 10000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", r, size)
		// test loads
		b.Run(name+"_load", func(b *testing.B) {
			err := r.Load(file, size)
			if err != nil {
				panic(err)
			}
		})
		r.Root = nil
		b.ResetTimer()
	}

}

func BenchmarkGapReport(b *testing.B) {
	gb := &GapBuffer{
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
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}

	for i := 1; i <= 5; i++ {

		size := 10000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", gb, size)
		// test loads
		err := gb.Load(file, size)
		if err != nil {
			panic(err)
		}
		b.ResetTimer()

		b.Run(name+"_report", func(b *testing.B) {
			disc, err := gb.Report()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%d", len(disc))
		})

		gb.Content = &DoublyLinkedList{
			Head:   nil,
			End:    nil,
			Length: 0,
		}
		b.ResetTimer()
	}
}

func BenchmarkRopeReport(b *testing.B) {
	r := &Rope{
		Root: nil,
	}
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}

	for i := 1; i <= 5; i++ {

		size := 10000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", r, size)
		// test loads
		err := r.Load(file, size)
		if err != nil {
			panic(err)
		}
		b.ResetTimer()

		b.Run(name+"_report", func(b *testing.B) {
			disc, err := r.Report()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%d", len(disc))
		})

		r.Root = nil
		b.ResetTimer()
	}
}

func BenchmarkGapReportR(b *testing.B) {
	gb := &GapBuffer{
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
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}
	var disc []rune
	for i := 1; i <= 1; i++ {

		size := 100000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", gb, size)
		// test loads
		err := gb.Load(file, size)
		if err != nil {
			panic(err)
		}
		b.ResetTimer()

		for j := 1; j <= 5; j++ {
			s, e := 15000, 15000*(j+1)
			nameadd := fmt.Sprintf("report_len_%d", e-s)
			b.Run(name+nameadd, func(b *testing.B) {
				disc, err = gb.ReportRange(s, e)
				if err != nil {
					panic(err)
				}
			})

		}
		b.Log(len(disc))
		b.ResetTimer()

		// second, scaling report start
		for j := 1; j <= 5; j++ {

			s, e := 15000*j, 15000*(j+1)
			nameadd := fmt.Sprintf("_report_%d:%d", s, e)
			b.Run(name+nameadd, func(b *testing.B) {
				disc, err = gb.ReportRange(s, e)
				if err != nil {
					panic(err)
				}
			})
		}
		b.Log(len(disc))
		b.ResetTimer()

		gb.Content = &DoublyLinkedList{
			Head:   nil,
			End:    nil,
			Length: 0,
		}
		b.ResetTimer()
	}
}

func BenchmarkRopeReportR(b *testing.B) {
	r := &Rope{
		Root: nil,
	}
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}
	var disc []rune
	for i := 1; i <= 1; i++ {

		size := 100000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", r, size)
		// test loads
		err := r.Load(file, size)
		if err != nil {
			panic(err)
		}
		b.ResetTimer()

		for j := 1; j <= 5; j++ {
			s, e := 15000, 15000*(j+1)
			nameadd := fmt.Sprintf("report_len_%d", e-s)
			b.Run(name+nameadd, func(b *testing.B) {
				disc, err = r.ReportRange(s, e)
				if err != nil {
					panic(err)
				}
			})

		}
		b.Log(len(disc))
		b.ResetTimer()

		// second, scaling report start
		for j := 1; j <= 5; j++ {

			s, e := 15000*j, 15000*(j+1)
			nameadd := fmt.Sprintf("_report_%d:%d", s, e)
			b.Run(name+nameadd, func(b *testing.B) {
				disc, err = r.ReportRange(s, e)
				if err != nil {
					panic(err)
				}
			})
		}

		b.ResetTimer()

		r.Root = nil
		b.ResetTimer()
	}
}

func BenchmarkGapInsert(b *testing.B) {
	gb := &GapBuffer{
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
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}

	for i := 1; i <= 5; i++ {

		size := 10000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", gb, size)
		// test loads
		err := gb.Load(file, size)
		if err != nil {
			panic(err)
		}
		b.ResetTimer()

		// test inserts
		// first, scaling input size
		base := "aaaaaaaaaa"
		for j := 100; j <= 150; j += 10 {
			thisins := StrToSlice(MultString(base, j))
			nameadd := fmt.Sprintf("_insert_len:%d", len(thisins))
			b.Run(name+nameadd, func(b *testing.B) {
				err = gb.Insert(100, thisins)
				if err != nil {
					panic(err)
				}
			})
		}
		b.ResetTimer()
		// second, scaling position
		for j := 1; j <= 5; j++ {
			thisins := StrToSlice(MultString(base, 110))
			pos := 1000 * j
			nameadd := fmt.Sprintf("_insert_at:%d", pos)
			b.Run(name+nameadd, func(b *testing.B) {
				err = gb.Insert(pos, thisins)
				if err != nil {
					panic(err)
				}
			})
		}
		b.ResetTimer()
	}
}

func BenchmarkRopeInsert(b *testing.B) {
	r := &Rope{Root: nil}
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}

	for i := 1; i <= 5; i++ {

		size := 10000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", r, size)
		// test loads
		err := r.Load(file, size)
		if err != nil {
			panic(err)
		}
		b.ResetTimer()

		// test inserts
		// first, scaling input size
		base := "aaaaaaaaaa"
		for j := 100; j <= 150; j += 10 {
			thisins := StrToSlice(MultString(base, j))
			nameadd := fmt.Sprintf("_insert_len:%d", len(thisins))
			b.Run(name+nameadd, func(b *testing.B) {
				err = r.Insert(100, thisins)
				if err != nil {
					panic(err)
				}
			})
		}
		b.ResetTimer()
		// second, scaling position
		for j := 1; j <= 5; j++ {
			thisins := StrToSlice(MultString(base, 110))
			pos := 1000 * j
			nameadd := fmt.Sprintf("_insert_at:%d", pos)
			b.Run(name+nameadd, func(b *testing.B) {
				err = r.Insert(pos, thisins)
				if err != nil {
					panic(err)
				}
			})
		}
		b.ResetTimer()
	}
}

func BenchmarkGapSplit(b *testing.B) {
	gb := &GapBuffer{
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
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}

	for i := 1; i <= 5; i++ {

		size := 10000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", gb, size)
		// test loads
		err := gb.Load(file, size)
		if err != nil {
			panic(err)
		}
		b.ResetTimer()
		// test split: scaling position
		for j := 1; j <= 5; j++ {
			pos := 10 * j
			nameadd := fmt.Sprintf("_split_at:%d", pos)

			ret := 0
			b.Run(name+nameadd, func(b *testing.B) {
				disct, err := gb.Split(gb.Length() - pos)
				if err != nil {
					panic(err)
				}
				ret += len(disct.ToString())

			})
		}
		b.ResetTimer()

		gb.Content = &DoublyLinkedList{
			Head:   nil,
			End:    nil,
			Length: 0,
		}
		b.ResetTimer()
	}
}

func BenchmarkRopeSplit(b *testing.B) {
	r := Rope{Root: nil}
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}

	for i := 1; i <= 5; i++ {

		size := 10000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", r, size)
		// test loads
		err := r.Load(file, size)
		if err != nil {
			panic(err)
		}
		b.ResetTimer()
		// test split: scaling position
		for j := 1; j <= 5; j++ {
			pos := 10 * j
			nameadd := fmt.Sprintf("_split_at:%d", pos)

			ret := 0
			b.Run(name+nameadd, func(b *testing.B) {
				disct, err := r.Split(r.Root.length - pos)
				if err != nil {
					panic(err)
				}
				ret += len(disct.ToString())

			})
		}
		b.ResetTimer()

		r.Root = nil
		b.ResetTimer()
	}
}

func BenchmarkGapDelete(b *testing.B) {
	gb := &GapBuffer{
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
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}

	var disc []rune
	for i := 1; i <= 1; i++ {

		size := 100000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", gb, size)
		// test loads
		err := gb.Load(file, size)
		if err != nil {
			panic(err)
		}
		b.ResetTimer()

		for j := 3; j <= 8; j++ {
			s, e := 10000, 10000*(j+1)
			nameadd := fmt.Sprintf("delete_len_%d", e-s)
			b.Run(name+nameadd, func(b *testing.B) {
				disc, err = gb.DeleteRange(s, e)
				if err != nil {
					panic(err)
				}
				gb.Insert(s, disc)
			})

		}
		b.Log(len(disc))
		b.ResetTimer()

		// second, scaling report start
		for j := 2; j <= 7; j++ {

			s, e := 10000*j, 10000*(j+1)
			nameadd := fmt.Sprintf("_delete_%d:%d", s, e)
			b.Run(name+nameadd, func(b *testing.B) {
				disc, err = gb.DeleteRange(s, e)
				if err != nil {
					panic(err)
				}
				gb.Insert(s, disc)
			})
		}

		b.ResetTimer()

		gb.Content = &DoublyLinkedList{
			Head:   nil,
			End:    nil,
			Length: 0,
		}
		b.ResetTimer()
	}
}

func BenchmarkRopeDelete(b *testing.B) {
	r := &Rope{Root: nil}
	file := &FileWrapper{
		Filename: "pg66576.txt",
	}

	var disc []rune
	for i := 1; i <= 1; i++ {

		size := 100000 * IntPow(2, i)
		name := fmt.Sprintf("%T_size_%d", r, size)
		// test loads
		err := r.Load(file, size)
		if err != nil {
			panic(err)
		}
		b.ResetTimer()

		for j := 3; j <= 8; j++ {
			s, e := 10000, 10000*(j+1)
			nameadd := fmt.Sprintf("delete_len_%d", e-s)
			b.Run(name+nameadd, func(b *testing.B) {
				disc, err = r.DeleteRange(s, e)
				if err != nil {
					panic(err)
				}
				r.Insert(s, disc)
			})

		}
		b.Log(len(disc))
		b.ResetTimer()

		// second, scaling report start
		for j := 2; j <= 7; j++ {

			s, e := 10000*j, 10000*(j+1)
			nameadd := fmt.Sprintf("_delete_%d:%d", s, e)
			b.Run(name+nameadd, func(b *testing.B) {
				disc, err = r.DeleteRange(s, e)
				if err != nil {
					panic(err)
				}
				r.Insert(s, disc)
			})
		}

		b.ResetTimer()

		r.Root = nil
		b.ResetTimer()
	}
}

func BenchmarkAll(b *testing.B) {
	testEls := []StorageType{
		&GapBuffer{
			Content: &DoublyLinkedList{
				Head:   nil,
				End:    nil,
				Length: 0,
			},
			GapStart:    nil,
			GapEnd:      nil,
			GapStartIdx: 0,
			GapLen:      0,
		},
		&Rope{
			Root: nil,
		},
	}

	file := &FileWrapper{
		Filename: "pg66576.txt",
	}
	var disc []rune
	var err error
	//var disct StorageType
	b.ResetTimer()
	for _, t := range testEls {
		for i := 1; i <= 2; i++ {

			size := 100000 * IntPow(2, i)
			name := fmt.Sprintf("%T_size_%d", t, size)
			// test loads
			b.Run(name+"_load", func(b *testing.B) {
				err = t.Load(file, size)
				if err != nil {
					panic(err)
				}
			})
			b.ResetTimer()

			// test reports
			b.Run(name+"_report", func(b *testing.B) {
				disc, err = t.Report()
				if err != nil {
					panic(err)
				}
			})
			b.Log(len(disc))
			b.ResetTimer()

			// test ranged reports
			// first, scaling length
			for j := 1; j <= 5; j++ {
				s, e := 5000, 5000*(j+1)
				nameadd := fmt.Sprintf("report_len_%d", e-s)
				b.Run(name+nameadd, func(b *testing.B) {
					disc, err = t.ReportRange(s, e)
					if err != nil {
						panic(err)
					}
				})

			}
			b.Log(len(disc))
			b.ResetTimer()

			// second, scaling report start
			for j := 1; j <= 5; j++ {

				s, e := 10000*j, 10000*(j+1)
				nameadd := fmt.Sprintf("_report_%d:%d", s, e)
				b.Run(name+nameadd, func(b *testing.B) {
					disc, err = t.ReportRange(s, e)
					if err != nil {
						panic(err)
					}
				})
			}
			b.Log(len(disc))
			b.ResetTimer()
			/*
				// test inserts
				// first, scaling input size
				base := "aaaaaaaaaa"
				for j := 1; j <= 5; j++ {
					thisins := StrToSlice(MultString(base, j))
					nameadd := fmt.Sprintf("_insert_len:%d", len(thisins))
					b.Run(name+nameadd, func(b *testing.B) {
						err = t.Insert(100, thisins)
						if err != nil {
							panic(err)
						}
					})
				}
				b.ResetTimer()
				// second, scaling position
				for j := 1; j <= 5; j++ {
					thisins := StrToSlice(base)
					pos := 100 * j
					nameadd := fmt.Sprintf("_insert_at:%d", pos)
					b.Run(name+nameadd, func(b *testing.B) {
						err = t.Insert(pos, thisins)
						if err != nil {
							panic(err)
						}
					})
				}
				b.ResetTimer()

				// test split: scaling position
				for j := 1; j <= 5; j++ {
					pos := 100 * j
					nameadd := fmt.Sprintf("_split_at:%d", pos)
					b.Run(name+nameadd, func(b *testing.B) {
						disct, err = t.Split(pos)
						if err != nil {
							panic(err)
						}
					})
					pb := disct.ToString()
					t.Append(StrToSlice(pb))
					b.ResetTimer()
				}
				b.ResetTimer()

				// test delete
				// first, scaling size
				for j := 1; j <= 5; j++ {
					s, e := 1000*i, 1000*(j+1)
					nameadd := fmt.Sprintf("_delete_len:%d", e-s)
					b.Run(name+nameadd, func(b *testing.B) {
						disc, err = t.DeleteRange(s, e)
						if err != nil {
							panic(err)
						}
					})
					t.Insert(s, disc)
					b.ResetTimer()
				}
				b.ResetTimer()
				// second, scaling position
				for j := 1; j <= 5; j++ {
					s, e := 1000*j, 1000*(j+1)
					nameadd := fmt.Sprintf("_delete_%d:%d", s, e)
					b.Run(name+nameadd, func(b *testing.B) {
						disc, err = t.DeleteRange(s, e)
						if err != nil {
							panic(err)
						}
					})
					t.Insert(s, disc)
					b.ResetTimer()
				}
				b.ResetTimer()
			*/
		}
		b.ResetTimer()
	}
	b.ResetTimer()

}

func BenchmarkReport(b *testing.B) {
	testEls := []StorageType{
		&GapBuffer{
			Content: &DoublyLinkedList{
				Head:   nil,
				End:    nil,
				Length: 0,
			},
			GapStart:    nil,
			GapEnd:      nil,
			GapStartIdx: 0,
			GapLen:      0,
		},
		&Rope{
			Root: nil,
		},
	}

	file := &FileWrapper{
		Filename: "pg66576.txt",
	}

	//fmt.Printf("size: %d\n", 10000)
	err := testEls[0].Load(file, 10000)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		testEls[0].ReportRange(10, 1000)
	}
}
