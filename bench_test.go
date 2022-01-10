package main

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkLoad(b *testing.B) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := make([]BufferTest, 0)
	for x := 1; x < 11; x++ {
		testCases = append(testCases, BufferTest{Length: x * 40000})
	}

	for _, tc := range testCases {
		content := GetContent("pg66576.txt", tc.Length)
		for _, bf := range buffers {
			b.Run(fmt.Sprintf("%T_load_%d", bf, tc.Length), func(t *testing.B) {
				bf.Load(content)
			})
		}
	}
}

func BenchmarkReport(b *testing.B) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := make([]BufferTest, 0)
	for x := 1; x < 11; x++ {
		testCases = append(testCases, BufferTest{Length: x * 40000})
	}
	var res []rune
	var err error
	for _, tc := range testCases {
		content := GetContent("pg66576.txt", tc.Length)
		for _, bf := range buffers {
			bf.Load(content)
			b.Run(fmt.Sprintf("%T_report_%d", bf, tc.Length), func(t *testing.B) {
				res, err = bf.Report()
				time.Sleep(time.Nanosecond)
				if err != nil {
					panic(err)
				}

			})
		}
	}
	fmt.Println(len(res))
}

func BenchmarkReportRange(b *testing.B) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := make([]BufferTest, 0)
	for x := 1; x < 11; x++ {
		amt := x * 40000
		testCases = append(testCases, BufferTest{Start: []int{amt / 8, amt / 4, amt / 2, (3 * amt) / 4}, End: []int{amt / 4, (3 * amt) / 8, (5 * amt) / 8, (7 * amt) / 8}, Length: amt})
	}
	var res []rune
	var err error
	for _, tc := range testCases {
		content := GetContent("pg66576.txt", tc.Length)
		for _, bf := range buffers {
			bf.Load(content)

			// first, changing size with constant report length
			b.Run(fmt.Sprintf("%T_reportRange_constLen_%d", bf, tc.Length), func(t *testing.B) {
				time.Sleep(time.Nanosecond)
				res, err = bf.ReportRange(10000, 30000)
				if err != nil {
					panic(err)
				}

			})

			for tIdx := range tc.End {
				start := tc.Start[tIdx]
				end := tc.End[tIdx]
				// second, changing length of report

				b.Run(fmt.Sprintf("%T_reportRange_diffLen_%d_%d:%d", bf, tc.Length, tc.Start[0], end), func(t *testing.B) {
					time.Sleep(time.Nanosecond)
					res, err = bf.ReportRange(tc.Start[0], end)
					if err != nil {
						panic(err)
					}

				})

				// third, changing position of the report with a constant length of n/8

				b.Run(fmt.Sprintf("%T_reportRange_changePos_%d_%d:%d", bf, tc.Length, start, end), func(t *testing.B) {
					time.Sleep(time.Nanosecond)
					res, err = bf.ReportRange(start, end)
					if err != nil {
						panic(err)
					}

				})

			}

		}
	}
	fmt.Println(len(res))
}

func BenchmarkInsert(b *testing.B) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := make([]BufferTest, 0)
	for x := 1; x < 6; x += 1 {
		testCases = append(testCases, BufferTest{
			Start:  []int{0, x * 10000, x * 20000, x * 40000, x * 60000, x * 80000},
			Length: x * 80000})
		//BufferTest{Start: []int{0, 250, 500, -1}, Length: 1000},
	}
	var res []rune
	var err error
	for _, tc := range testCases {
		content := GetContent("pg66576.txt", tc.Length)

		insert := GetContent("pg66576.txt", tc.Length)
		app := make([]rune, len(insert))
		for i := 0; i < len(insert); i++ {
			app[i] = rune(insert[i])
		}
		for _, bf := range buffers {
			bf.Load(content)
			// constant position 2000 character insert
			b.Run(fmt.Sprintf("%T_insert_%d_const", bf, tc.Length), func(t *testing.B) {
				err = bf.Insert(40000, app[:2000])
				//time.Sleep(time.Nanosecond)
				if err != nil {
					panic(err)
				}
				b.StopTimer()
				bf.Load(content)
			})
			// varying start

			for t_idx := range tc.Start {
				start := tc.Start[t_idx]
				b.Run(fmt.Sprintf("%T_insert_%d_Start_%d", bf, tc.Length, start), func(t *testing.B) {
					err = bf.Insert(start, app[:2000])
					if err != nil {
						panic(err)
					}
					b.StopTimer()
					bf.Load(content)
				})
				// varying input size:
				if t_idx != 0 {
					tmp := app[:t_idx*5000]
					b.Run(fmt.Sprintf("%T_insert_%d_Size_%d", bf, tc.Length, (t_idx*5000)), func(t *testing.B) {
						err = bf.Insert(100, tmp)
						if err != nil {
							panic(err)
						}
						b.StopTimer()
						bf.Load(content)
					})
				}

			}

		}
	}
	fmt.Println(len(res))
}

func BenchmarkSplit(b *testing.B) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := make([]BufferTest, 0)
	for x := 1; x < 11; x += 1 {
		amt := x * 40000
		testCases = append(testCases, BufferTest{
			Start:  []int{amt / 8, amt / 4, amt / 2, (3 * amt) / 4, (7 * amt) / 8},
			Length: amt})
	}
	var res StorageType
	var err error
	for _, tc := range testCases {
		content := GetContent("pg66576.txt", tc.Length)
		for _, bf := range buffers {
			bf.Load(content)
			// first, constant place

			b.Run(fmt.Sprintf("%T_split_%d_const", bf, len(content)), func(t *testing.B) {
				res, err = bf.Split(20000)
				//time.Sleep(time.Nanosecond)
				if err != nil {
					panic(err)
				}
				b.StopTimer()
				bf.Load(content)
			})

			// second, varying place
			for t_idx := range tc.Start {
				start := tc.Start[t_idx]
				b.Run(fmt.Sprintf("%T_split_%d_vary_%d", bf, len(content), start), func(t *testing.B) {
					res, err = bf.Split(start)
					//time.Sleep(time.Nanosecond)
					if err != nil {
						panic(err)
					}

					b.StopTimer()
					bf.Load(content)
				})
			}
		}
	}
	fmt.Println(len(res.ToString()))
}

func BenchmarkDelete(b *testing.B) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := make([]BufferTest, 0)
	for x := 1; x < 11; x += 1 {
		amt := x * 40000
		testCases = append(testCases, BufferTest{
			Start:  []int{amt / 8, amt / 4, amt / 2, (3 * amt) / 4, (7 * amt) / 8},
			End:    []int{amt / 4, (3 * amt) / 8, (5 * amt) / 8, (7 * amt) / 8, amt - 1},
			Length: amt})
	}
	var res []rune
	var err error
	for _, tc := range testCases {
		content := GetContent("pg66576.txt", tc.Length)
		for _, bf := range buffers {
			bf.Load(content)
			// first, constant position

			b.Run(fmt.Sprintf("%T_delete_%d_const", bf, tc.Length), func(t *testing.B) {
				res, err = bf.DeleteRange(10000, 30000)
				//time.Sleep(time.Nanosecond)
				if err != nil {
					panic(err)
				}
				b.StopTimer()
				bf.Load(content)
			})

			for tIdx := range tc.End {
				start := tc.Start[tIdx]
				end := tc.End[tIdx]
				// second, changing length of delete

				b.Run(fmt.Sprintf("%T_delete_diffLen_%d_%d:%d", bf, tc.Length, tc.Start[0], end), func(t *testing.B) {
					time.Sleep(time.Nanosecond)
					res, err = bf.DeleteRange(tc.Start[0], end)

					if err != nil {
						panic(err)
					}
					b.StopTimer()
					bf.Load(content)
				})

				// third, changing position of the delete with a constant length of n/8

				b.Run(fmt.Sprintf("%T_delete_changePos_%d_%d:%d", bf, tc.Length, start, end), func(t *testing.B) {
					time.Sleep(time.Nanosecond)
					res, err = bf.DeleteRange(start, end)

					if err != nil {
						panic(err)
					}
					b.StopTimer()
					bf.Load(content)
				})

			}
		}
	}
	fmt.Println(len(res))
}

func BenchmarkIndex(b *testing.B) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	testCases := make([]BufferTest, 0)
	for x := 1; x < 11; x += 1 {
		amt := x * 40000
		testCases = append(testCases, BufferTest{
			Start:  []int{0, amt / 8, amt / 4, amt / 2, (3 * amt) / 4, (7 * amt) / 8, amt - 1},
			Length: amt})
		//BufferTest{Start: []int{50, 10, 30}, End: []int{ -1, 20, 50}, Length: 100},
	}
	var res rune
	var err error
	for _, tc := range testCases {
		content := GetContent("pg66576.txt", tc.Length)
		for _, bf := range buffers {
			bf.Load(content)
			// first, constant position
			b.Run(fmt.Sprintf("%T_index_%d_const", bf, tc.Length), func(t *testing.B) {
				res, err = bf.Index(20000)
				time.Sleep(time.Nanosecond)
				if err != nil {
					panic(err)
				}

			})

			// second, changing position
			for _, idx := range tc.Start {
				b.Run(fmt.Sprintf("%T_index_%d_at_%d", bf, tc.Length, idx), func(t *testing.B) {
					res, err = bf.Index(idx)
					time.Sleep(time.Nanosecond)
					if err != nil {
						panic(err)
					}

				})
			}

		}
	}
	fmt.Println(res)
}

func BenchmarkAppend(b *testing.B) {
	buffers := make([]StorageType, 2)
	buffers[0] = mkRope()
	buffers[1] = mkGapBuf()

	extra := make([]StorageType, 2)
	extra[0] = mkRope()
	extra[1] = mkGapBuf()

	testCases := make([]BufferTest, 0)
	for x := 1; x < 11; x += 1 {
		amt := x * 40000
		testCases = append(testCases, BufferTest{
			Start:  []int{2000, amt / 8, amt / 4, amt / 2, (3 * amt) / 4, (7 * amt) / 8, amt},
			Length: amt})
	}

	var res []rune
	var err error
	for _, tc := range testCases {
		content := GetContent("pg66576.txt", tc.Length)
		for i, bf := range buffers {
			for _, length := range tc.Start {
				bf.Load(content)
				extra[i].Load(content[:length])
				b.Run(fmt.Sprintf("%T_Append_%d_%d", bf, len(content), length), func(t *testing.B) {
					err = bf.Concat(extra[i])
					time.Sleep(time.Nanosecond)
					if err != nil {
						panic(err)
					}
					b.StopTimer()
					bf.Load(content)

				})

				res, err = bf.Report()
				if err != nil {
					panic(err)
				}
			}

		}
	}
	fmt.Println(len(res))

}
