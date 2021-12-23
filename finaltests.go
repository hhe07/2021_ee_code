package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
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

	for i := 1; i <= 1; i++ {
		for t, s := range testEls {
			if t == 0 {
				fmt.Println("testing gap buffer")
			} else {
				fmt.Println("testing rope")
			}

			// test load
			start := time.Now()
			size := 100000 * IntPow(3, i)
			fmt.Printf("size: %d\n", size)
			err := s.Load(file, size)
			var disc []rune
			var disct StorageType
			if err != nil {
				panic(err)
			}
			t := time.Now()
			fmt.Printf("load time for %d chars: %v \n", size, t.Sub(start))

			// test report
			start = time.Now()
			disc, err = s.Report()
			//log.Println(string(disc))
			if err != nil {
				panic(err)
			}
			t = time.Now()
			fmt.Printf("report time for %d chars: %v \n", len(disc), t.Sub(start))

			// test report range: first, scaling report length
			for j := 1; j <= 5; j++ {
				start := time.Now()
				disc, err = s.ReportRange(10000*i, 10000*(j+1))
				if err != nil {
					panic(err)

				}
				//t := time.Now()
				//log.Printf(string(disc))
				// !: t.Sub(start) doesn't return meaningful values, neither does anything else
				fmt.Printf("report time for %d chars from %d to %d: %d \n", len(disc), 1000*i, 1000*(j+1), time.Since(start).Nanoseconds())
			}
			// test report range: second, scaling report start
			for j := 1; j <= 5; j++ {
				start = time.Now()
				disc, err = s.ReportRange(j*100, (j+1)*100)
				if err != nil {
					panic(err)
				}
				t = time.Now()
				fmt.Printf("report time for %d chars from %d to %d: %v \n", len(disc), 100*j, 100*(j+1), t.Sub(start))
			}

			// test insert: first, scaling input size
			base := "aaaaaaaaaa"
			for j := 1; j <= 5; j++ {
				start = time.Now()
				thisins := StrToSlice(MultString(base, j))
				err = s.Insert(100, thisins)
				if err != nil {
					panic(err)
				}
				t = time.Now()
				fmt.Printf("insert time for %d chars at 100: %v \n", len(thisins), t.Sub(start))
			}
			// test insert: second, scaling position
			for j := 1; j <= 5; j++ {
				start = time.Now()
				thisins := StrToSlice(base)
				err = s.Insert(100*j, thisins)
				if err != nil {
					panic(err)
				}
				t = time.Now()
				fmt.Printf("insert time for 10 chars at %d: %v \n", 100*j, t.Sub(start))
			}

			// test split: scaling position
			for j := 1; j <= 5; j++ {
				p := time.Now().UnixNano()
				//start = time.Now()
				disct, err = s.Split(100 * j)
				if err != nil {
					panic(err)
				}
				//t = time.Now()
				v := time.Now().UnixNano()
				fmt.Printf("split time at %d: %v \n", 100*j, v-p) //t.Sub(start)
				pb := disct.ToString()
				s.Append(StrToSlice(pb))

			}
			// test delete: first, scaling size
			for j := 1; j <= 5; j++ {
				start = time.Now()
				disc, err = s.DeleteRange(100, 100+100*j)
				if err != nil {
					panic(err)
				}
				t = time.Now()
				fmt.Printf("delete time for %d chars from %d to %d: %v \n", len(disc), 100, 100+100*j, t.Sub(start))
				s.Insert(100, disc)
			}
			// test delete: second, scaling position
			for j := 1; j <= 5; j++ {
				start = time.Now()
				disc, err = s.DeleteRange(j*100, (j+1)*100)
				if err != nil {
					panic(err)
				}
				t = time.Now()
				fmt.Printf("delete time for %d chars from %d to %d: %v \n", len(disc), 100*j, 100*(j+1), t.Sub(start))
			}
		}

	}

}

func IntPow(b, p int) int {
	return int(math.Pow(float64(b), float64(p)))

}

func StrToSlice(s string) []rune {
	ret := make([]rune, 0)
	for _, c := range s {
		ret = append(ret, c)
	}
	return ret
}

func MultString(s string, ct int) string {
	ret := ""
	for p := 0; p < ct; p++ {
		ret += s
	}
	return ret
}
