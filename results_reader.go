package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
)

func ReadRes(res []byte) map[string][]string {
	// convert res into string
	interStr := string(res)
	// first, split results along new line
	splitline := regexp.MustCompile("\n").Split(interStr, -1)

	ret := make(map[string][]string)
	testRE := regexp.MustCompile("Benchmark.*-8")
	opRE := regexp.MustCompile(" 0.*ns/op")
	for _, s := range splitline {
		testname := testRE.FindString(s)
		testname = regexp.MustCompile("Benchmark.*main").ReplaceAllString(testname, "")
		optime := opRE.FindString(s)
		optime = regexp.MustCompile("ns/op").ReplaceAllString(optime, "")
		if _, ok := ret[testname]; ok {
			// if the testname is already in, append to it
			ret[testname] = append(ret[testname], optime)
		} else {
			ret[testname] = make([]string, 1)
			ret[testname][0] = optime
		}
	}
	/*
		for k, r := range ret {
			fmt.Printf("/%s/: [", k)
			for _, s := range r {
				fmt.Printf("%s,", s)
			}
			fmt.Printf("]\n")
		}
	*/
	return ret
}

func PrepareForWrite(in map[string][]string) [][]string {
	ret := make([][]string, 0)
	for k, v := range in {
		tmp := make([]string, 0)
		tmp = append(tmp, k)
		tmp = append(tmp, v...)
		ret = append(ret, tmp)
	}
	return ret
}

func main() {
	test := "insert2"
	ret, err := os.ReadFile(fmt.Sprintf("results/%s.txt", test))
	if err != nil {
		panic(err)
	}
	res := PrepareForWrite(ReadRes(ret))

	f, err := os.Create(fmt.Sprintf("results/%s.csv", test))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	cwrite := csv.NewWriter(f)
	cwrite.WriteAll(res)

}
